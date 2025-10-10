package main

import (
   "41.neocities.org/media/hulu"
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

type flag_set struct {
   address  string
   cdm      net.Cdm
   email    string
   filters  net.Filters
   media    string
   password string
}

func (f *flag_set) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.cdm.ClientId = f.media + "/client_id.bin"
   f.cdm.PrivateKey = f.media + "/private_key.pem"
   flag.StringVar(&f.cdm.ClientId, "C", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.cdm.PrivateKey, "P", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.email, "e", "", "email")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   return nil
}

func (f *flag_set) do_email() error {
   data, err := hulu.NewAuthenticate(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.media+"/hulu/Authenticate", data)
}

func (f *flag_set) email_password() bool {
   if f.email != "" {
      if f.password != "" {
         return true
      }
   }
   return false
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   switch {
   case set.address != "":
      err = set.do_address()
   case set.email_password():
      err = set.do_email()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.media + "/hulu/Authenticate")
   if err != nil {
      return err
   }
   var auth hulu.Authenticate
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   err = auth.Refresh()
   if err != nil {
      return err
   }
   id, err := hulu.Id(f.address)
   if err != nil {
      return err
   }
   deep, err := auth.DeepLink(id)
   if err != nil {
      return err
   }
   data, err = auth.Playlist(deep)
   if err != nil {
      return err
   }
   var play hulu.Playlist
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.StreamUrl)
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return play.License(data)
   }
   return f.filters.Filter(resp, &f.cdm)
}

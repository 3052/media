package main

import (
   "41.neocities.org/media/cineMember"
   "41.neocities.org/net"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

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

type flag_set struct {
   address  string
   cdm      net.Cdm
   email    string
   filters  net.Filters
   media    string
   password string
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) do_email() error {
   data, err := cineMember.NewUser(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.media+"/cineMember/User", data)
}

func (f *flag_set) email_password() bool {
   if f.email != "" {
      if f.password != "" {
         return true
      }
   }
   return false
}

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.media + "/cineMember/User")
   if err != nil {
      return err
   }
   var user cineMember.User
   err = user.Unmarshal(data)
   if err != nil {
      return err
   }
   var address cineMember.Address
   err = address.Parse(f.address)
   if err != nil {
      return err
   }
   article, err := address.Article()
   if err != nil {
      return err
   }
   asset, ok := article.Film()
   if !ok {
      return errors.New(".Film()")
   }
   data, err = user.Play(article, asset)
   if err != nil {
      return err
   }
   var play cineMember.Play
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   title, ok := play.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(title.Manifest)
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return title.License(data)
   }
   return f.filters.Filter(resp, &f.cdm)
}

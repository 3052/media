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

func (f *flag_set) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.cdm.ClientId = f.media + "/client_id.bin"
   f.cdm.PrivateKey = f.media + "/private_key.pem"
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.dash, "d", "", "DASH ID")
   flag.StringVar(&f.email, "email", "", "email")
   flag.StringVar(&f.cdm.PrivateKey, "p", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.password, "password", "", "password")
   flag.Parse()
   return nil
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.email != "" {
      if set.password != "" {
         err = set.do_email()
      }
   } else if set.address != "" {
      err = set.do_address()
   } else if set.dash != "" {
      err = set.do_dash()
   } else {
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

type flag_set struct {
   address  string
   dash     string
   cdm        net.License
   email    string
   media    string
   password string
}

func (f *flag_set) do_email() error {
   data, err := hulu.NewAuthenticate(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.media+"/hulu/Authenticate", data)
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
   deep, err := auth.DeepLink(hulu.Id(f.address))
   if err != nil {
      return err
   }
   data, err = auth.Playlist(deep)
   if err != nil {
      return err
   }
   err = write_file(f.media+"/hulu/Playlist", data)
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
   return net.Mpd(f.media+"/Mpd", resp)
}

func (f *flag_set) do_dash() error {
   data, err := os.ReadFile(f.media + "/hulu/Playlist")
   if err != nil {
      return err
   }
   var play hulu.Playlist
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   f.cdm.Widevine = func(data []byte) ([]byte, error) {
      return play.Widevine(data)
   }
   return f.cdm.Download(f.media+"/Mpd", f.dash)
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}


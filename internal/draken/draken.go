package main

import (
   "41.neocities.org/media/draken"
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) authenticate() error {
   data, err := draken.NewLogin(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/draken/Login", data)
}

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.cdm.ClientId = f.cache + "/L3/client_id.bin"
   f.cdm.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.StringVar(&f.cdm.ClientId, "C", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.cdm.PrivateKey, "P", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.email, "e", "", "email")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   return nil
}

func (f *flag_set) email_password() bool {
   if f.email != "" {
      if f.password != "" {
         return true
      }
   }
   return false
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   switch {
   case set.address != "":
      err = set.download()
   case set.email_password():
      err = set.authenticate()
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
   cache    string
   password string
}

func (f *flag_set) download() error {
   data, err := os.ReadFile(f.cache + "/draken/Login")
   if err != nil {
      return err
   }
   var login draken.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   var movie draken.Movie
   err = movie.New(path.Base(f.address))
   if err != nil {
      return err
   }
   title, err := login.Entitlement(movie)
   if err != nil {
      return err
   }
   data, err = login.Playback(&movie, title)
   if err != nil {
      return err
   }
   var play draken.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.Playlist)
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return login.License(&play, data)
   }
   return f.filters.Filter(resp, &f.cdm)
}

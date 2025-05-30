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

func (f *flag_set) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.e.ClientId = f.media + "/client_id.bin"
   f.e.PrivateKey = f.media + "/private_key.pem"
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
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
   case set.password != "":
      err = set.authenticate()
   case set.address != "":
      err = set.download()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) authenticate() error {
   data, err := draken.NewLogin(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.media+"/draken/Login", data)
}

func (f *flag_set) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/draken/Login")
      if err != nil {
         return err
      }
      var login draken.Login
      err = login.Unmarshal(data)
      if err != nil {
         return err
      }
      data, err = os.ReadFile(f.media + "/draken/Playback")
      if err != nil {
         return err
      }
      var play draken.Playback
      err = play.Unmarshal(data)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return login.Widevine(&play, data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/draken/Login")
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
   err = write_file(f.media+"/draken/Playback", data)
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
   return net.Mpd(f.media+"/Mpd", resp)
}

type flag_set struct {
   address  string
   dash     string
   e        net.License
   email    string
   media    string
   password string
}


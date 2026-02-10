package main

import (
   "41.neocities.org/media/draken"
   "41.neocities.org/media/internal"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

type flags struct {
   address  string
   dash     string
   e        internal.License
   email    string
   media    string
   password string
}

func (f *flags) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.e.ClientId = f.media + "/client_id.bin"
   f.e.PrivateKey = f.media + "/private_key.pem"
   return nil
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   switch {
   case f.password != "":
      err := f.authenticate()
      if err != nil {
         panic(err)
      }
   case f.address != "":
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}

func (f *flags) write_file(name string, data []byte) error {
   log.Println("WriteFile", f.media + name)
   return os.WriteFile(f.media + name, data, os.ModePerm)
}

func (f *flags) authenticate() error {
   data, err := draken.NewLogin(f.email, f.password)
   if err != nil {
      return err
   }
   return f.write_file("/draken/Login", data)
}

func (f *flags) download() error {
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
      return f.e.Download(f.media + "/Mpd", f.dash)
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
   err = f.write_file("/draken/Playback", data)
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
   return internal.Mpd(f.media + "/Mpd", resp)
}

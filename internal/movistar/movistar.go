package main

import (
   "41.neocities.org/media/movistar"
   "41.neocities.org/media/internal"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

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

type flags struct {
   e        internal.License
   email    string
   media    string
   password string
   dash     string
   movistar int64
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.Int64Var(&f.movistar, "b", 0, "movistar ID")
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
   case f.movistar >= 1:
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}

///

func (f *flags) write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flags) authenticate() error {
   data, err := movistar.NewToken(f.email, f.password)
   if err != nil {
      return err
   }
   err = f.write_file("/movistar/Token", data)
   if err != nil {
      t.Fatal(err)
   }
}

func (f *flags) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/movistar/Playlist")
      if err != nil {
         return err
      }
      var play movistar.Playlist
      err = play.Unmarshal(data)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return play.Widevine(data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/movistar/Authenticate")
   if err != nil {
      return err
   }
   var auth movistar.Authenticate
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   err = auth.Refresh()
   if err != nil {
      return err
   }
   deep, err := auth.DeepLink(f.movistar)
   if err != nil {
      return err
   }
   data, err = auth.Playlist(deep)
   if err != nil {
      return err
   }
   err = f.write_file("/movistar/Playlist", data)
   if err != nil {
      return err
   }
   var play movistar.Playlist
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.StreamUrl)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}


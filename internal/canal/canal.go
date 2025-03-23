package main

import (
   "41.neocities.org/media/canal"
   "41.neocities.org/media/internal"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) write_file(name string, data []byte) error {
   log.Println("WriteFile", f.media+name)
   return os.WriteFile(f.media+name, data, os.ModePerm)
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

type flags struct {
   e        internal.License
   email    string
   media    string
   password string
   dash     string
   address string
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.StringVar(&f.address, "a", "", "address")
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

///

func (f *flags) authenticate() error {
   var ticket1 ticket
   err = ticket1.New()
   if err != nil {
      t.Fatal(err)
   }
   data, err = ticket1.token(username, password)
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home+"/media/canal/token", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
   // OLD
   data, err := canal.NewAuthenticate(f.email, f.password)
   if err != nil {
      return err
   }
   return f.write_file("/canal/Authenticate", data)
}

func (f *flags) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/canal/Playlist")
      if err != nil {
         return err
      }
      var play canal.Playlist
      err = play.Unmarshal(data)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return play.Widevine(data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/canal/Authenticate")
   if err != nil {
      return err
   }
   var auth canal.Authenticate
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   err = auth.Refresh()
   if err != nil {
      return err
   }
   deep, err := auth.DeepLink(f.entity)
   if err != nil {
      return err
   }
   data, err = auth.Playlist(deep)
   if err != nil {
      return err
   }
   err = f.write_file("/canal/Playlist", data)
   if err != nil {
      return err
   }
   var play canal.Playlist
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

package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/max"
   "41.neocities.org/platform/mullvad"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) do_initiate() error {
   var st max.St
   err := st.New()
   if err != nil {
      return err
   }
   log.Println("Create", f.media + "/max/St")
   file, err := os.Create(f.media + "/max/St")
   if err != nil {
      return err
   }
   defer file.Close()
   _, err = fmt.Fprint(file, st)
   if err != nil {
      return err
   }
   initiate, err := st.Initiate()
   if err != nil {
      return err
   }
   fmt.Println(initiate)
   return nil
}

type flags struct {
   dash     string
   e        internal.License
   initiate bool
   login    bool
   media    string
   mullvad  bool
   url      max.WatchUrl
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
   flag.Var(&f.url, "a", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.BoolVar(
      &f.initiate, "initiate", false, "/authentication/linkDevice/initiate",
   )
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.BoolVar(
      &f.login, "login", false, "/authentication/linkDevice/login",
   )
   flag.BoolVar(&f.mullvad, "m", false, "Mullvad")
   flag.Parse()
   if f.mullvad {
      http.DefaultClient.Transport = &mullvad.Transport{}
   }
   switch {
   case f.initiate:
      err := f.do_initiate()
      if err != nil {
         panic(err)
      }
   case f.login:
      err := f.do_login()
      if err != nil {
         panic(err)
      }
   case f.url.VideoId != "":
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

func (f *flags) do_login() error {
   data, err := os.ReadFile(f.media + "/max/St")
   if err != nil {
      return err
   }
   var st max.St
   err = st.Set(string(data))
   if err != nil {
      return err
   }
   data, err = st.Login()
   if err != nil {
      return err
   }
   return f.write_file("/max/Login", data)
}

func (f *flags) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/max/Playback")
      if err != nil {
         return err
      }
      var play max.Playback
      err = play.Unmarshal(data)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return play.Widevine(data)
      }
      return f.e.Download(f.media + "/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/max/Login")
   if err != nil {
      return err
   }
   var login max.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = login.Playback(&f.url)
   if err != nil {
      return err
   }
   var play max.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   err = f.write_file("/max/Playback", data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.Fallback.Manifest.Url[0])
   if err != nil {
      return err
   }
   return internal.Mpd(f.media + "/Mpd", resp)
}

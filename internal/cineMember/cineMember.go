package main

import (
   "41.neocities.org/media/cineMember"
   "41.neocities.org/media/internal"
   "41.neocities.org/platform/mullvad"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

type flags struct {
   address  cineMember.Address
   dash     string
   e        internal.License
   email    string
   media    string
   mullvad  bool
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
   flag.Var(&f.address, "a", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.BoolVar(&f.mullvad, "m", false, "Mullvad")
   flag.Parse()
   if f.mullvad {
      http.DefaultClient.Transport = &mullvad.Transport{}
   }
   switch {
   case f.password != "":
      err := f.write_user()
      if err != nil {
         panic(err)
      }
   case f.address[0] != "":
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}

func (f *flags) write_file(name string, data []byte) error {
   log.Println("WriteFile", f.media+name)
   return os.WriteFile(f.media+name, data, os.ModePerm)
}

func (f *flags) write_user() error {
   data, err := cineMember.NewUser(f.email, f.password)
   if err != nil {
      return err
   }
   return f.write_file("/cineMember/User", data)
}

func (f *flags) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/cineMember/Play")
      if err != nil {
         return err
      }
      var play cineMember.Play
      err = play.Unmarshal(data)
      if err != nil {
         return err
      }
      title, _ := play.Dash()
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return title.Widevine(data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/cineMember/User")
   if err != nil {
      return err
   }
   var user cineMember.User
   err = user.Unmarshal(data)
   if err != nil {
      return err
   }
   article, err := f.address.Article()
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
   err = f.write_file("/cineMember/Play", data)
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
   return internal.Mpd(f.media+"/Mpd", resp)
}

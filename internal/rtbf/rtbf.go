package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/rtbf"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

type flags struct {
   address  rtbf.Address
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
   flag.Var(&f.address, "a", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   switch {
   case f.password != "":
      err := f.authenticate()
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

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flags) authenticate() error {
   data, err := rtbf.NewLogin(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.media + "/rtbf/Login", data)
}

func (f *flags) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/rtbf/Entitlement")
      if err != nil {
         return err
      }
      var title rtbf.Entitlement
      err = title.Unmarshal(data)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return title.Widevine(data)
      }
      return f.e.Download(f.media + "/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/rtbf/Login")
   if err != nil {
      return err
   }
   var login rtbf.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   jwt, err := login.Jwt()
   if err != nil {
      return err
   }
   gigya, err := jwt.Login()
   if err != nil {
      return err
   }
   content, err := f.address.Content()
   if err != nil {
      return err
   }
   data, err = gigya.Entitlement(content)
   if err != nil {
      return err
   }
   var title rtbf.Entitlement
   err = title.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.media + "/rtbf/Entitlement", data)
   if err != nil {
      return err
   }
   format, ok := title.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(format.MediaLocator)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media + "/Mpd", resp)
}

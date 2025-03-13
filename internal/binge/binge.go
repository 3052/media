package main

import (
   "41.neocities.org/media/binge"
   "41.neocities.org/media/internal"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

type flags struct {
   binge    int
   dash     string
   e        internal.License
   email    string
   media    string
   password string
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.IntVar(&f.binge, "b", 0, "binge ID")
   flag.StringVar(&f.email, "email", "", "email")
   flag.StringVar(&f.password, "password", "", "password")
   flag.Parse()
   switch {
   case f.password != "":
      err := f.authenticate()
      if err != nil {
         panic(err)
      }
   case f.binge >= 1:
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}

func (f *flags) authenticate() error {
   data, err := binge.NewAuth(f.email, f.password)
   if err != nil {
      return err
   }
   return f.write_file("/binge/Auth", data)
}

func (f *flags) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/binge/TokenService")
      if err != nil {
         return err
      }
      var token binge.TokenService
      err = token.Unmarshal(data)
      if err != nil {
         return err
      }
      data, err = os.ReadFile(f.media + "/binge/Play")
      if err != nil {
         return err
      }
      var play binge.Play
      err = play.Unmarshal(data)
      if err != nil {
         return err
      }
      stream, _ := play.Dash()
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return token.Widevine(stream, data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/binge/Auth")
   if err != nil {
      return err
   }
   var auth binge.Auth
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = auth.Token()
   if err != nil {
      return err
   }
   var token binge.TokenService
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   err = f.write_file("/binge/TokenService", data)
   if err != nil {
      return err
   }
   data, err = token.Play(f.binge)
   if err != nil {
      return err
   }
   var play binge.Play
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   err = f.write_file("/binge/Play", data)
   if err != nil {
      return err
   }
   stream, ok := play.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(stream.Manifest)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
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

func (f *flags) write_file(name string, data []byte) error {
   log.Println("WriteFile", f.media+name)
   return os.WriteFile(f.media+name, data, os.ModePerm)
}

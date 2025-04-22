package main

import (
   "41.neocities.org/media/criterion"
   "41.neocities.org/media/internal"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.IntVar(&internal.ThreadCount, "t", 1, "thread count")
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

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flags) authenticate() error {
   data, err := criterion.NewToken(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.media + "/criterion/Token", data)
}

func (f *flags) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/criterion/Files")
      if err != nil {
         return err
      }
      var files criterion.Files
      err = files.Unmarshal(data)
      if err != nil {
         return err
      }
      file, _ := files.Dash()
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return file.Widevine(data)
      }
      return f.e.Download(f.media + "/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/criterion/Token")
   if err != nil {
      return err
   }
   var token criterion.Token
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = token.Refresh()
   if err != nil {
      return err
   }
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.media + "/criterion/Token", data)
   if err != nil {
      return err
   }
   video, err := token.Video(path.Base(f.address))
   if err != nil {
      return err
   }
   data, err = token.Files(video)
   if err != nil {
      return err
   }
   var files criterion.Files
   err = files.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.media + "/criterion/Files", data)
   if err != nil {
      return err
   }
   file, ok := files.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(file.Links.Source.Href)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media + "/Mpd", resp)
}
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

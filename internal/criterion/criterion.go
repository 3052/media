package main

import (
   "41.neocities.org/media/criterion"
   "41.neocities.org/net"
   "errors"
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
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.IntVar(&net.Threads, "t", 2, "threads")
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
   data, err := criterion.NewToken(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.media+"/criterion/Token", data)
}

func (f *flag_set) download() error {
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
      return f.e.Download(f.media+"/Mpd", f.dash)
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
   err = write_file(f.media+"/criterion/Token", data)
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
   err = write_file(f.media+"/criterion/Files", data)
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

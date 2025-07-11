package main

import (
   "41.neocities.org/media/criterion"
   "41.neocities.org/net"
   "errors"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

func main() {
   http.DefaultTransport = &http.Transport{
      Proxy: func(req *http.Request) (*url.URL, error) {
         log.Println(req.Method, req.URL)
         return nil, nil
      },
   }
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

type flag_set struct {
   address  string
   dash     string
   cdm        net.Cdm
   email    string
   media    string
   password string
   filters net.Filters
}
func (f *flag_set) download() error {
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
   file, ok := files.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(file.Links.Source.Href)
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return file.License(data)
   }
   return f.filters.Filter(resp, &f.cdm)
}

func (f *flag_set) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.cdm.ClientId = f.media + "/client_id.bin"
   f.cdm.PrivateKey = f.media + "/private_key.pem"
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.StringVar(&f.cdm.PrivateKey, "k", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   return nil
}

package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/criterion"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   job   maya.WidevineJob
   name    string
   // 1
   email    string
   password string
   // 2
   address  string
   // 3
   dash string
}

///

func (f *command) New() error {
   switch {
   case set.address != "":
      err = set.do_address()
   case set.email_password():
      err = set.do_token()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
   var err error
   f.name, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.name = filepath.ToSlash(f.name)
   f.job.ClientId = f.name + "/L3/client_id.bin"
   f.job.PrivateKey = f.name + "/L3/private_key.pem"
   flag.StringVar(&f.job.ClientId, "C", f.job.ClientId, "client ID")
   flag.StringVar(&f.job.PrivateKey, "P", f.job.PrivateKey, "private key")
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.email, "e", "", "email")
   flag.Var(&f.filters, "f", maya.FilterUsage)
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   return nil
   if f.email != "" {
      if f.password != "" {
         return true
      }
   }
   return false
}
func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *command) do_token() error {
   data, err := criterion.FetchToken(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.name+"/criterion/Token", data)
}

func (f *command) do_address() error {
   data, err := os.ReadFile(f.name + "/criterion/Token")
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
   err = write_file(f.name+"/criterion/Token", data)
   if err != nil {
      return err
   }
   video, err := token.Video(path.Base(f.address))
   if err != nil {
      return err
   }
   files, err := token.Files(video)
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
   f.job.Send = func(data []byte) ([]byte, error) {
      return file.Widevine(data)
   }
   return f.filters.Filter(resp, &f.job)
}


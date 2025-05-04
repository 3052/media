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
   e        internal.License
   media    string
   
   email    string
   password string
   
   address string
   
   dash     string
}

func (f *flags) do_email() error {
   data, err := cineMember.NewUser(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.media + "/cineMember/User", data)
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "d", "", "DASH ID")
   flag.StringVar(&f.email, "email", "", "email")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "password", "", "password")
   flag.IntVar(&internal.ThreadCount, "t", 1, "thread count")
   flag.Parse()
   if f.email != "" {
      if f.password != "" {
         err = f.do_email()
      }
   } else if f.address != "" {
      err = f.do_address()
   } else if f.dash != "" {
      err = f.do_dash()
   } else {
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

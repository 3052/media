package main

import (
   "41.neocities.org/media/rtbf"
   "41.neocities.org/net"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

type flags struct {
   dash     string
   cdm        net.Cdm
   email    string
   media    string
   password string
   address  string
}

func (f *flags) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.cdm.ClientId = f.media + "/client_id.bin"
   f.cdm.PrivateKey = f.media + "/private_key.pem"
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.cdm.PrivateKey, "k", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.StringVar(&f.address, "a", "", "address")
   flag.Parse()
   return nil
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flags) do_password() error {
   data, err := rtbf.NewLogin(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.media+"/rtbf/Login", data)
}

func main() {
   var set flags
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.email != "" {
      if set.password != "" {
         err = set.do_password()
      }
   } else if set.address != "" {
      err = set.do_address()
   } else {
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flags) do_address() error {
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
   var address rtbf.Address
   address.New(f.address)
   content, err := address.Content()
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
   format, ok := title.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(format.MediaLocator)
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return title.License(data)
   }
   return f.cdm.Download(f.media+"/Mpd", f.dash)
}

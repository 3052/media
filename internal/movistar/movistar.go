package main

import (
   "41.neocities.org/media/movistar"
   "41.neocities.org/net"
   "flag"
   "log"
   "os"
   "path/filepath"
)

func (f *flag_set) do_movistar() error {
   data, err := os.ReadFile(f.media + "/movistar/Token")
   if err != nil {
      return err
   }
   var token movistar.Token
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = os.ReadFile(f.media + "/movistar/Device")
   if err != nil {
      return err
   }
   var device movistar.Device
   err = device.Unmarshal(data)
   if err != nil {
      return err
   }
   oferta, err := token.Oferta()
   if err != nil {
      return err
   }
   initData, err := oferta.InitData(device)
   if err != nil {
      return err
   }
   data, err = movistar.NewDetails(f.movistar)
   if err != nil {
      return err
   }
   var details movistar.Details
   err = details.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := details.Mpd()
   if err != nil {
      return err
   }
   session, err := device.Session(initData, &details)
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return session.License(data)
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
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.email, "email", "", "email")
   flag.Int64Var(&f.movistar, "m", 0, "movistar ID")
   flag.StringVar(&f.cdm.PrivateKey, "p", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.password, "password", "", "password")
   flag.Parse()
   return nil
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.email != "" {
      if set.password != "" {
         err = set.do_email()
      }
   } else if set.movistar >= 1 {
      err = set.do_movistar()
   } else {
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flag_set) do_email() error {
   data, err := movistar.NewToken(f.email, f.password)
   if err != nil {
      return err
   }
   var token movistar.Token
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.media+"/movistar/Token", data)
   if err != nil {
      return err
   }
   oferta, err := token.Oferta()
   if err != nil {
      return err
   }
   data1, err := token.Device(oferta)
   if err != nil {
      return err
   }
   return write_file(f.media+"/movistar/Device", data1)
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

type flag_set struct {
   cdm        net.Cdm
   email    string
   media    string
   movistar int64
   password string
   filters net.Filters
}

package main

import (
   "41.neocities.org/media/canal"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flags) do_email() error {
   var ticket canal.Ticket
   err := ticket.New()
   if err != nil {
      return err
   }
   token, err := ticket.Token(f.email, f.password)
   if err != nil {
      return err
   }
   data, err := canal.NewSession(token.SsoToken)
   if err != nil {
      return err
   }
   return write_file(f.media+"/canal/Session", data)
}

func (f *flags) do_refresh() error {
   data, err := os.ReadFile(f.media + "/canal/Session")
   if err != nil {
      return err
   }
   var session canal.Session
   err = session.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = canal.NewSession(session.SsoToken)
   if err != nil {
      return err
   }
   return write_file(f.media+"/canal/Session", data)
}

func (f *flags) do_address() error {
   var fields canal.Fields
   err := fields.New(f.address)
   if err != nil {
      return err
   }
   fmt.Println(
      canal.AlgoliaConvertTracking, "=",
      fields.Get(canal.AlgoliaConvertTracking),
   )
   return nil
}

func (f *flags) do_season() error {
   data, err := os.ReadFile(f.media + "/canal/Session")
   if err != nil {
      return err
   }
   var session canal.Session
   err = session.Unmarshal(data)
   if err != nil {
      return err
   }
   assets, err := session.Assets(f.asset, f.season)
   if err != nil {
      return err
   }
   for i, asset := range assets {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&asset)
   }
   return nil
}

func (f *flags) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.license.ClientId = f.media + "/client_id.bin"
   f.license.PrivateKey = f.media + "/private_key.pem"
   return nil
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(
      &f.license.PrivateKey, "p", f.license.PrivateKey, "private key",
   )
   flag.StringVar(&f.license.ClientId, "c", f.license.ClientId, "client ID")
   /////////////////////////////////////////////////////////////////////////
   flag.StringVar(&f.email, "email", "", "email")
   flag.StringVar(&f.password, "password", "", "password")
   ///////////////////////////////////////////////////////
   flag.BoolVar(&f.refresh, "r", false, "refresh")
   ///////////////////////////////////////////////
   flag.StringVar(&f.address, "a", "", "address")
   //////////////////////////////////////////////
   flag.Int64Var(&f.season, "s", 0, "season")
   //////////////////////////////////////////
   flag.StringVar(&f.asset, "asset", "", "asset ID")
   flag.Var(&f.bitrate, "b", "bitrate")
   flag.Parse()
   if f.email != "" {
      if f.password != "" {
         err = f.do_email()
      }
   } else if f.refresh {
      err = f.do_refresh()
   } else if f.address != "" {
      err = f.do_address()
   } else if f.asset != "" {
      if f.season >= 1 {
         err = f.do_season()
      } else {
         err = f.do_asset()
      }
   } else {
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

type flags struct {
   media   string
   license net.License
   ///////////////////
   email    string
   password string
   ///////////////
   refresh bool
   /////////////
   address string
   ///////////////
   season int64
   ///////////////
   asset   string
   bitrate net.Bitrate
}

func (f *flags) do_asset() error {
   data, err := os.ReadFile(f.media + "/canal/Session")
   if err != nil {
      return err
   }
   var session canal.Session
   err = session.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = session.Play(f.asset)
   if err != nil {
      return err
   }
   var play canal.Play
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.Url)
   if err != nil {
      return err
   }
   f.license.Widevine = func(data []byte) ([]byte, error) {
      return play.Widevine(data)
   }
   return f.license.Bitrate(resp, &f.bitrate)
}

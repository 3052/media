package main

import (
   "41.neocities.org/media/canal"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   http.DefaultTransport = &http.Transport{
      Proxy: func(req *http.Request) (*url.URL, error) {
         urlVar, err := http.ProxyFromEnvironment(req)
         if err != nil {
            return nil, err
         }
         if urlVar != nil {
            log.Println("Proxy", urlVar)
         }
         log.Println(req.Method, req.URL)
         return urlVar, nil
      },
   }
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.email != "" {
      if set.password != "" {
         err = set.do_email()
      }
   } else if set.refresh {
      err = set.do_refresh()
   } else if set.address != "" {
      err = set.do_address()
   } else if set.asset != "" {
      if set.season >= 1 {
         err = set.do_season()
      } else {
         err = set.do_asset()
      }
   } else {
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

type flag_set struct {
   media   string
   cdm net.Cdm
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
   filters net.Filters
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) do_email() error {
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

func (f *flag_set) do_asset() error {
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
   f.cdm.License = func(data []byte) ([]byte, error) {
      return play.License(data)
   }
   return f.filters.Filter(resp, &f.cdm)
}

func (f *flag_set) do_refresh() error {
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

func (f *flag_set) do_address() error {
   var fields canal.Fields
   err := fields.New(f.address)
   if err != nil {
      return err
   }
   fmt.Println("asset id =", fields.AssetId())
   return nil
}

func (f *flag_set) do_season() error {
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
   flag.StringVar(&f.asset, "A", "", "asset ID")
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.email, "email", "", "email")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(
      &f.cdm.PrivateKey, "p", f.cdm.PrivateKey, "private key",
   )
   flag.StringVar(&f.password, "password", "", "password")
   flag.BoolVar(&f.refresh, "r", false, "refresh")
   flag.Int64Var(&f.season, "s", 0, "season")
   flag.Parse()
   return nil
}

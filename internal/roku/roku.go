package main

import (
   "41.neocities.org/media/roku"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func main() {
   http.DefaultTransport = &roku.Transport
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      log.Fatal(err)
   }
   switch {
   case set.code_write:
      err = set.do_code()
   case set.roku != "":
      err = set.do_roku()
   case set.token_write:
      err = set.do_token()
   default:
      flag.Usage()
   }
   if err != nil {
      log.Fatal(err)
   }
}

func (f *flag_set) do_code() error {
   data, err := (*roku.Code).AccountToken(nil)
   if err != nil {
      return err
   }
   var token roku.AccountToken
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.cache+"/roku/AccountToken", data)
   if err != nil {
      return err
   }
   data1, err := token.Activation()
   if err != nil {
      return err
   }
   var activation roku.Activation
   err = activation.Unmarshal(data1)
   if err != nil {
      return err
   }
   fmt.Println(&activation)
   return write_file(f.cache+"/roku/Activation", data1)
}

func (f *flag_set) do_token() error {
   data, err := os.ReadFile(f.cache + "/roku/AccountToken")
   if err != nil {
      return err
   }
   var token roku.AccountToken
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = os.ReadFile(f.cache + "/roku/Activation")
   if err != nil {
      return err
   }
   var activation roku.Activation
   err = activation.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = token.Code(&activation)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/roku/Code", data)
}

type flag_set struct {
   cache       string
   code_write  bool
   config      net.Config
   filters     net.Filters
   roku        string
   token_read  bool
   token_write bool
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) do_roku() error {
   var code *roku.Code
   if f.token_read {
      data, err := os.ReadFile(f.cache + "/roku/Code")
      if err != nil {
         return err
      }
      code = &roku.Code{}
      err = code.Unmarshal(data)
      if err != nil {
         return err
      }
   }
   data, err := code.AccountToken()
   if err != nil {
      return err
   }
   var token roku.AccountToken
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   data1, err := token.Playback(f.roku)
   if err != nil {
      return err
   }
   var play roku.Playback
   err = play.Unmarshal(data1)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.Url)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return play.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
}
func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.config.ClientId = f.cache + "/L3/client_id.bin"
   f.config.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.StringVar(&f.config.ClientId, "C", f.config.ClientId, "client ID")
   flag.StringVar(&f.config.PrivateKey, "P", f.config.PrivateKey, "private key")
   flag.BoolVar(&f.token_write, "T", false, "write token")
   flag.BoolVar(&f.code_write, "c", false, "write code")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.roku, "r", "", "Roku ID")
   flag.BoolVar(&f.token_read, "t", false, "read token")
   flag.Parse()
   return nil
}

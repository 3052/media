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

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.e.ClientId = f.cache + "/L3/client_id.bin"
   f.e.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "d", "", "dash ID")
   flag.BoolVar(&f.code_write, "code", false, "write code")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.StringVar(&f.roku, "r", "", "Roku ID")
   flag.BoolVar(&f.token_read, "t", false, "read token")
   flag.BoolVar(&f.token_write, "token", false, "write token")
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
   case set.code_write:
      err = set.do_code()
   case set.token_write:
      err = set.do_token()
   case set.roku != "":
      err = set.do_roku()
   case set.dash != "":
      err = set.do_dash()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
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
   e          net.License
   cache      string
   token_read bool
   code_write  bool
   token_write bool
   roku        string
   dash        string
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
   err = write_file(f.cache+"/roku/Playback", data1)
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
   return net.Mpd(f.cache+"/Mpd", resp)
}

func (f *flag_set) do_dash() error {
   data, err := os.ReadFile(f.cache + "/roku/Playback")
   if err != nil {
      return err
   }
   var play roku.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   f.e.Widevine = func(data []byte) ([]byte, error) {
      return play.Widevine(data)
   }
   return f.e.Download(f.cache+"/Mpd", f.dash)
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

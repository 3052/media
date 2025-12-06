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
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
      return "L"
   })
   var program runner
   err := program.run()
   if err != nil {
      log.Fatal(err)
   }
}

func (r *runner) run() error {
   var err error
   r.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   r.cache = filepath.ToSlash(r.cache)
   r.config.ClientId = r.cache + "/L3/client_id.bin"
   r.config.PrivateKey = r.cache + "/L3/private_key.pem"
   
   flag.StringVar(&r.config.ClientId, "C", r.config.ClientId, "client ID")
   flag.StringVar(&r.config.PrivateKey, "P", r.config.PrivateKey, "private key")
   flag.BoolVar(&r.connection, "c", false, "connection")
   flag.StringVar(&r.dash, "d", "", "DASH ID")
   flag.BoolVar(&r.get_credentials, "g", false, "get credentials")
   flag.StringVar(&r.roku, "r", "", "Roku ID")
   flag.BoolVar(&r.set_credentials, "s", false, "set credentials")
   flag.Parse()
   if r.connection {
      return r.do_connection()
   }
   if r.set_credentials {
      return r.do_credentials()
   }
   if r.roku != "" {
      return r.do_roku()
   }
   if r.dash != "" {
      return r.do_dash()
   }
   flag.Usage()
   return nil
}

func (r *runner) write(cache *roku.Cache) error {
   data, err := json.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", r.cache + "/roku/Cache")
   return os.WriteFile(r.cache + "/roku/Cache", data, os.ModePerm)
}

func (r *runner) read() (*roku.Cache, error) {
   data, err := os.ReadFile(r.cache + "/roku/Cache")
   if err != nil {
      return nil, err
   }
   var cache roku.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return nil, err
   }
   return &cache, nil
}

func (r *runner) do_connection() error {
   var (
      cache roku.Cache
      credentials *roku.Credentials
      err error
   )
   cache.Connection, err = credentials.NewConnection()
   if err != nil {
      return err
   }
   cache.LinkCode, err = cache.Connection.RequestLinkCode()
   if err != nil {
      return err
   }
   fmt.Println(cache.LinkCode)
   return r.write(&cache)
}

type runner struct {
   cache       string
   config      net.Config
   // 1
   connection  bool
   // 2
   set_credentials bool
   // 3
   roku        string
   get_credentials  bool
   // 4
   dash string
}

///

func (r *runner) do_credentials() error {
   data, err := os.ReadFile(r.cache + "/roku/AccountToken")
   if err != nil {
      return err
   }
   var token roku.AccountToken
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = os.ReadFile(r.cache + "/roku/Activation")
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
   return write_file(r.cache+"/roku/Code", data)
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (r *runner) do_roku() error {
   var code *roku.Code
   if r.get_credentials {
      data, err := os.ReadFile(r.cache + "/roku/Code")
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
   data1, err := token.Playback(r.roku)
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
   r.config.Send = func(data []byte) ([]byte, error) {
      return play.Widevine(data)
   }
   return r.filters.Filter(resp, &r.config)
}

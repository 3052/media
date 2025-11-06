package main

import (
   "41.neocities.org/media/mubi"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.cache + "/mubi/Authenticate")
   if err != nil {
      return err
   }
   var auth mubi.Authenticate
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   slug, err := mubi.FilmSlug(f.address)
   if err != nil {
      return err
   }
   film_id, err := mubi.FilmId(slug)
   if err != nil {
      return err
   }
   err = auth.Viewing(film_id)
   if err != nil {
      return err
   }
   secure, err := auth.SecureUrl(film_id)
   if err != nil {
      return err
   }
   resp, err := http.Get(secure.Url)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return auth.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) do_code() error {
   data, err := mubi.NewLinkCode()
   if err != nil {
      return err
   }
   var code mubi.LinkCode
   err = code.Unmarshal(data)
   if err != nil {
      return err
   }
   fmt.Println(&code)
   return write_file(f.cache+"/mubi/LinkCode", data)
}

func (f *flag_set) do_auth() error {
   data, err := os.ReadFile(f.cache + "/mubi/LinkCode")
   if err != nil {
      return err
   }
   var code mubi.LinkCode
   err = code.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = code.Authenticate()
   if err != nil {
      return err
   }
   return write_file(f.cache+"/mubi/Authenticate", data)
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
   flag.StringVar(&f.address, "a", "", "address")
   flag.BoolVar(&f.auth, "auth", false, "authenticate")
   flag.BoolVar(&f.code, "code", false, "link code")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.IntVar(&f.config.Threads, "t", 12, "threads")
   flag.Parse()
   return nil
}
func main() {
   http.DefaultTransport = &mubi.Transport
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   switch {
   case set.address != "":
      err = set.do_address()
   case set.auth:
      err = set.do_auth()
   case set.code:
      err = set.do_code()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

type flag_set struct {
   address string
   auth    bool
   cache   string
   code    bool
   config  net.Config
   filters net.Filters
}

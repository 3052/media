package main

import (
   "41.neocities.org/media/mubi"
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
   http.DefaultTransport = &http.Transport{
      Protocols: &http.Protocols{}, // github.com/golang/go/issues/25793
      Proxy: func(req *http.Request) (*url.URL, error) {
         log.Println(req.Method, req.URL)
         return nil, nil
      },
   }
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   switch {
   case set.code:
      err = set.do_code()
   case set.auth:
      err = set.do_auth()
   case set.slug != "":
      err = set.do_slug()
   default:
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

type flag_set struct {
   auth    bool
   cdm     net.Cdm
   code    bool
   filters net.Filters
   cache   string
   slug    mubi.Slug
}

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.cdm.ClientId = f.cache + "/L3/client_id.bin"
   f.cdm.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.Func("a", "address", func(data string) error {
      return f.slug.Parse(data)
   })
   flag.BoolVar(&f.auth, "auth", false, "authenticate")
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.BoolVar(&f.code, "code", false, "link code")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.cdm.PrivateKey, "p", f.cdm.PrivateKey, "private key")
   flag.Parse()
   return nil
}

func (f *flag_set) do_slug() error {
   data, err := os.ReadFile(f.cache + "/mubi/Authenticate")
   if err != nil {
      return err
   }
   var auth mubi.Authenticate
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return auth.License(data)
   }
   film, err := f.slug.Film()
   if err != nil {
      return err
   }
   err = auth.Viewing(film)
   if err != nil {
      return err
   }
   data, err = auth.SecureUrl(film)
   if err != nil {
      return err
   }
   var secure mubi.SecureUrl
   err = secure.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(secure.Url)
   if err != nil {
      return err
   }
   return f.filters.Filter(resp, &f.cdm)
}

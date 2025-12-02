package main

import (
   "41.neocities.org/media/mubi"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "io"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

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
   flag.BoolVar(&f.code, "c", false, "link code")
   flag.StringVar(&f.dash, "d", "", "DASH ID")
   flag.BoolVar(&f.session, "s", false, "session")
   flag.BoolVar(&f.text, "t", false, "text")
   flag.Parse()
   return nil
}

func get(address string) error {
   resp, err := http.Get(address)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   return write_file(path.Base(address), data)
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func main() {
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".dash" {
         return ""
      }
      return "L"
   })
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      log.Fatal(err)
   }
   if set.code {
      err = set.do_code()
   } else if set.session {
      err = set.do_session()
   } else if set.address != "" {
      if set.text {
         err = set.do_text()
      } else {
         err = set.do_address()
      }
   } else if set.dash != "" {
      err = set.do_dash()
   } else {
      flag.Usage()
   }
   if err != nil {
      log.Fatal(err)
   }
}

func (f *flag_set) do_code() error {
   var cache mubi.Cache
   err = cache.LinkCode.Fetch()
   if err != nil {
      return err
   }
   fmt.Println(cache.LinkCode)
   data, err := json.Marshal(cache)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/mubi/Cache", data)
}

func (f *flag_set) do_session() error {
   data, err := os.ReadFile(f.cache + "/mubi/Cache")
   if err != nil {
      return err
   }
   var cache mubi.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   cache.Session, err = cache.LinkCode.Session()
   if err != nil {
      return err
   }
   data, err = json.Marshal(cache)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/mubi/Cache", data)
}

func (f *flag_set) do_text() error {
   slug, err := mubi.FilmSlug(f.address)
   if err != nil {
      return err
   }
   film_id, err := mubi.FilmId(slug)
   if err != nil {
      return err
   }
   data, err := os.ReadFile(f.cache + "/mubi/Cache")
   if err != nil {
      return err
   }
   var cache mubi.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   secure, err := cache.Session.SecureUrl(film_id)
   if err != nil {
      return err
   }
   for _, text := range secure.TextTrackUrls {
      err = get(text.Url)
      if err != nil {
         return err
      }
   }
   return nil
}

type flag_set struct {
   config  net.Config
   cache   string
   // 1
   code    bool
   // 2
   session    bool
   // 3
   text    bool
   // 4
   address string
   // 5
   dash string
}

///

func (f *flag_set) do_address() error {
   slug, err := mubi.FilmSlug(f.address)
   if err != nil {
      return err
   }
   film_id, err := mubi.FilmId(slug)
   if err != nil {
      return err
   }
   data, err := os.ReadFile(f.cache + "/mubi/Cache")
   if err != nil {
      return err
   }
   var cache mubi.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   err = cache.Session.Viewing(film_id)
   if err != nil {
      return err
   }
   secure, err := cache.Session.SecureUrl(film_id)
   if err != nil {
      return err
   }
   
   resp, err := http.Get(secure.Url)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return session.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
}

func (f *flag_set) do_dash() error {
   data, err := os.ReadFile(f.cache + "/mubi/Cache")
   if err != nil {
      return err
   }
   var session mubi.Session
   err = session.Unmarshal(data)
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
   err = session.Viewing(film_id)
   if err != nil {
      return err
   }
   secure, err := session.SecureUrl(film_id)
   if err != nil {
      return err
   }
   resp, err := http.Get(secure.Url)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return session.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
}

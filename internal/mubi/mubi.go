package main

import (
   "41.neocities.org/media/mubi"
   "41.neocities.org/net"
   "encoding/json"
   "flag"
   "fmt"
   "io"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

type command struct {
   config  net.Config
   name   string
   // 1
   code    bool
   // 2
   session    bool
   // 3
   address string
   // 4
   dash string
}

///

type user_cache struct {
   LinkCode *LinkCode
   Mpd      *url.URL
   MpdBody  []byte
   Session *Session
}

func (f *command) New() error {
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
   flag.IntVar(&f.config.Threads, "T", 2, "threads")
   flag.StringVar(&f.address, "a", "", "address")
   flag.BoolVar(&f.code, "c", false, "link code")
   flag.StringVar(&f.dash, "d", "", "DASH ID")
   flag.BoolVar(&f.session, "s", false, "session")
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
      return "LP"
   })
   log.SetFlags(log.Ltime)
   var set command
   err := set.New()
   if err != nil {
      log.Fatal(err)
   }
   if set.code {
      err = set.do_code()
   } else if set.session {
      err = set.do_session()
   } else if set.address != "" {
      err = set.do_address()
   } else if set.dash != "" {
      err = set.do_dash()
   } else {
      flag.Usage()
   }
   if err != nil {
      log.Fatal(err)
   }
}

func (f *command) do_code() error {
   var code mubi.LinkCode
   err := code.Fetch()
   if err != nil {
      return err
   }
   fmt.Println(&code)
   data, err := json.Marshal(mubi.user_cache{LinkCode: &code})
   if err != nil {
      return err
   }
   return write_file(f.name+"/mubi/user_cache", data)
}

func (f *command) do_session() error {
   data, err := os.ReadFile(f.name + "/mubi/user_cache")
   if err != nil {
      return err
   }
   var cache mubi.user_cache
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
   return write_file(f.name+"/mubi/user_cache", data)
}

func (f *command) do_address() error {
   slug, err := mubi.FilmSlug(f.address)
   if err != nil {
      return err
   }
   film_id, err := mubi.FilmId(slug)
   if err != nil {
      return err
   }
   data, err := os.ReadFile(f.name + "/mubi/user_cache")
   if err != nil {
      return err
   }
   var cache mubi.user_cache
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
   err = secure.Mpd(&cache)
   if err != nil {
      return err
   }
   data, err = json.Marshal(cache)
   if err != nil {
      return err
   }
   err = write_file(f.name + "/mubi/user_cache", data)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}

func (f *command) do_dash() error {
   data, err := os.ReadFile(f.name + "/mubi/user_cache")
   if err != nil {
      return err
   }
   var cache mubi.user_cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return cache.Session.Widevine(data)
   }
   return f.config.Download(cache.MpdBody, cache.Mpd, f.dash)
}

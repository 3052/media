package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/mubi"
   "encoding/xml"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Session.Widevine(data)
   }
   return c.job.DownloadDash(cache.Mpd.Body, cache.Mpd.Url, c.dash)
}

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func read(name string) (*user_cache, error) {
   data, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   cache := &user_cache{}
   err = xml.Unmarshal(data, cache)
   if err != nil {
      return nil, err
   }
   return cache, nil
}

type user_cache struct {
   LinkCode *mubi.LinkCode
   Mpd      *mubi.Mpd
   Session  *mubi.Session
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".dash" {
         return ""
      }
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/mubi/userCache.xml"

   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.StringVar(&c.address, "a", "", "address")
   flag.BoolVar(&c.code, "c", false, "link code")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.BoolVar(&c.session, "s", false, "session")
   flag.Parse()

   if c.code {
      return c.do_code()
   }
   if c.session {
      return c.do_session()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
}

func (c *command) do_code() error {
   var code mubi.LinkCode
   err := code.Fetch()
   if err != nil {
      return err
   }
   fmt.Println(&code)
   return write(c.name, &user_cache{LinkCode: &code})
}

func (c *command) do_session() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   cache.Session, err = cache.LinkCode.Session()
   if err != nil {
      return err
   }
   return write(c.name, cache)
}

func (c *command) do_address() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   slug, err := mubi.FilmSlug(c.address)
   if err != nil {
      return err
   }
   film_id, err := mubi.FetchId(slug)
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
   cache.Mpd, err = secure.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Mpd.Body, cache.Mpd.Url)
}

type command struct {
   job  maya.WidevineJob
   name string
   // 1
   code bool
   // 2
   session bool
   // 3
   address string
   // 4
   dash string
}

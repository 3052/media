package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/mubi"
   "encoding/xml"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

type command struct {
   address string
   code    bool
   config  maya.Config
   dash    string
   name    string
   session bool
}

type user_cache struct {
   LinkCode *mubi.LinkCode
   Mpd      *url.URL
   MpdBody  []byte
   Session  *mubi.Session
}

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      return cache.Session.Widevine(data)
   }
   return c.config.Download(cache.Mpd, cache.MpdBody, c.dash)
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
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/mubi/user_cache.xml"

   flag.StringVar(&c.config.ClientId, "C", c.config.ClientId, "client ID")
   flag.StringVar(&c.config.PrivateKey, "P", c.config.PrivateKey, "private key")
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

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
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
   film_id, err := mubi.FilmId(slug)
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
   cache.Mpd, cache.MpdBody, err = secure.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.Representations(cache.Mpd, cache.MpdBody)
}

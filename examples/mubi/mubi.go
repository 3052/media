package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/rosso/mubi.xml"
   // 1
   flag.BoolVar(&c.code, "c", false, "link code")
   // 2
   flag.BoolVar(&c.session, "s", false, "session")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   // 3, 4
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.IntVar(&c.job.Threads, "t", 2, "threads")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   maya.SetProxy(func(req *http.Request) (string, bool) {
      if path.Ext(req.URL.Path) == ".dash" {
         return "", false
      }
      return c.proxy, true
   })
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
   return maya.Usage([][]string{
      {"c"},
      {"s"},
      {"a", "x"},
      {"d", "x", "t", "C", "P"},
   })
}

type command struct {
   name string
   // 1
   code bool
   // 2
   session bool
   // 3
   address string
   // 3, 4
   proxy string
   // 4
   dash string
   job  maya.WidevineJob
}

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = cache.Session.Widevine
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

type user_cache struct {
   Dash     *mubi.Dash
   LinkCode *mubi.LinkCode
   Session  *mubi.Session
}

func (c *command) do_code() error {
   var code mubi.LinkCode
   err := code.Fetch()
   if err != nil {
      return err
   }
   fmt.Println(&code)
   return maya.Write(c.name, &user_cache{LinkCode: &code})
}

func (c *command) do_session() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   cache.Session, err = cache.LinkCode.Session()
   if err != nil {
      return err
   }
   return maya.Write(c.name, cache)
}

func (c *command) do_address() error {
   cache, err := maya.Read[user_cache](c.name)
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
   cache.Dash, err = secure.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

func main() {
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

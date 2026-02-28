package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("mubi")
   // 1
   flag.BoolVar(&c.code, "c", false, "link code")
   // 2
   flag.BoolVar(&c.session, "s", false, "session")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
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
   return maya.Usage([][]string{
      {"c"},
      {"s"},
      {"a", "x"},
      {"d", "x", "C", "P"},
   })
}

type command struct {
   cache maya.Cache
   // 1
   code bool
   // 2
   session bool
   // 3
   address string
   // 4
   dash string
   job  maya.WidevineJob
}

func (c *command) do_dash() error {
   var dash mubi.Dash
   err := c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   var session mubi.Session
   err = c.cache.Get("Session", &session)
   if err != nil {
      return err
   }
   c.job.Send = session.Widevine
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func (c *command) do_code() error {
   var link_code mubi.LinkCode
   err := link_code.Fetch()
   if err != nil {
      return err
   }
   fmt.Println(&link_code)
   return c.cache.Set("LinkCode", link_code)
}

func (c *command) do_session() error {
   var link_code mubi.LinkCode
   err := c.cache.Get("LinkCode", &link_code)
   if err != nil {
      return err
   }
   session, err := link_code.Session()
   if err != nil {
      return err
   }
   return c.cache.Set("Session", session)
}

func (c *command) do_address() error {
   slug, err := mubi.FilmSlug(c.address)
   if err != nil {
      return err
   }
   film_id, err := mubi.FetchId(slug)
   if err != nil {
      return err
   }
   var session mubi.Session
   err = c.cache.Get("Session", &session)
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
   dash, err := secure.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".dash"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

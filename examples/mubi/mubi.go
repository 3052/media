package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "flag"
   "fmt"
   "log"
)

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := cache.Setup("rosso/mubi.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 2
   flag.BoolVar(&c.code, "c", false, "link code")
   // 3
   flag.BoolVar(&c.session, "s", false, "session")
   // 4
   flag.StringVar(&c.address, "a", "", "address")
   // 5
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "C", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "P", job.PrivateKey, "private key")
   flag.Parse()
   err = maya.SetProxy(c.proxy, "*.dash")
   if err != nil {
      return err
   }
   if c.code {
      return c.do_code()
   }
   if c.session {
      return c.do_session()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"x"},
      {"c"},
      {"s"},
      {"a"},
      {"d", "C", "P"},
   })
}

func (c *client) do_code() error {
   c.LinkCode = &mubi.LinkCode{}
   err := c.LinkCode.Fetch()
   if err != nil {
      return err
   }
   fmt.Println(c.LinkCode)
   return cache.Write(c)
}

var cache maya.Cache

var job maya.WidevineJob

func (c *client) do_session() error {
   return cache.Update(c, func() error {
      var err error
      c.Session, err = c.LinkCode.Session()
      return err
   })
}

func (c *client) do_address() error {
   slug, err := mubi.FilmSlug(c.address)
   if err != nil {
      return err
   }
   film_id, err := mubi.FetchId(slug)
   if err != nil {
      return err
   }
   err = cache.Read(c)
   if err != nil {
      return err
   }
   err = c.Session.Viewing(film_id)
   if err != nil {
      return err
   }
   secure, err := c.Session.SecureUrl(film_id)
   if err != nil {
      return err
   }
   c.Dash, err = secure.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = c.Session.Widevine
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
}

type client struct {
   Dash     *mubi.Dash
   LinkCode *mubi.LinkCode
   Session  *mubi.Session
   // 1
   proxy string
   // 2
   code bool
   // 3
   session bool
   // 4
   address string
   // 5
   dash_id string
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

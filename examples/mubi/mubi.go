package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "flag"
   "fmt"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/mubi.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   // 1
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   // 2
   code := maya.BoolVar(new(bool), "c", "link code")
   // 3
   session := maya.BoolVar(new(bool), "s", "session")
   // 4
   address := maya.StringVar(&c.address, "a", "address")
   // 5
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case set[widevine]:
      return cache.Write(c)
   case set[code]:
      return c.do_code()
   case set[session]:
      return with_cache(c.do_session)
   case set[address]:
      return with_cache(c.do_address)
   case set[dash_id]:
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {code},
      {session},
      {address},
      {dash_id},
   })
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.dash")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_code() error {
   var err error
   c.LinkCode, err = mubi.FetchLinkCode()
   if err != nil {
      return err
   }
   fmt.Println(c.LinkCode)
   return cache.Write(c)
}

func (c *client) do_session() error {
   var err error
   c.Session, err = c.LinkCode.Session()
   if err != nil {
      return err
   }
   return cache.Write(c)
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
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Session.Widevine,
   )
}

type client struct {
   Dash     *mubi.Dash
   LinkCode *mubi.LinkCode
   Session  *mubi.Session
   // 1
   Job maya.Job
   // 4
   address string
   // 5
   dash_id string
}

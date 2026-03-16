package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/peacock"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *client) do() error {
   err := cache.Setup("rosso/peacock.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c, true)
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.Job.Widevine, "w", c.Job.Widevine, "Widevine")
   // 4
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.IntVar(&c.Job.Threads, "t", 2, "threads")
   set := maya.Parse()
   if set["e"] {
      if set["p"] {
         return c.do_email_password()
      }
   }
   if set["a"] {
      return c.do_address()
   }
   if set["w"] {
      return cache.Write(c)
   }
   if set["d"] {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"a"},
      {"w"},
      {"d"},
   })
}

func (c *client) do_email_password() error {
   var err error
   c.Cookie, err = peacock.FetchIdSession(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
   if cache.Error != nil {
      return cache.Error
   }
   var token peacock.Token
   err := token.Fetch(c.Cookie)
   if err != nil {
      return err
   }
   c.Playout, err = token.Playout(path.Base(c.address))
   if err != nil {
      return err
   }
   endpoint, err := c.Playout.Fastly()
   if err != nil {
      return err
   }
   c.Dash, err = endpoint.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

type client struct {
   Cookie  *http.Cookie
   Dash    *peacock.Dash
   Playout *peacock.Playout
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   Job maya.Job
   // 4
   dash_id string
}
func (c *client) do_dash_id() error {
   if cache.Error != nil {
      return cache.Error
   }
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playout.Widevine,
   )
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

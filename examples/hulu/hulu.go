package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "flag"
   "log"
)

func (c *client) do_address() error {
   err := c.Session.TokenRefresh()
   if err != nil {
      return err
   }
   deep_link, err := c.Session.DeepLink(hulu.Id(c.address))
   if err != nil {
      return err
   }
   c.Playlist, err = c.Session.Playlist(deep_link)
   if err != nil {
      return err
   }
   c.Dash, err = c.Playlist.Dash()
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
   Dash     *hulu.Dash
   Playlist *hulu.Playlist
   Session  *hulu.Session
   // 1
   Job maya.Job
   // 2
   email    string
   password string
   // 3
   address string
   // 4
   dash_id string
}
var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4,*.mp4a")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   err := cache.Setup("rosso/hulu.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c)
   // 1
   flag.StringVar(&c.Job.PlayReady, "P", c.Job.PlayReady, "PlayReady")
   // 2
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   // 4
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   set := maya.Parse()
   switch {
   case set["P"]:
      return cache.Write(c)
   case set["e"] && set["p"]:
      return c.do_email_password()
   }
   if err != nil {
      return err
   }
   switch {
   case set["a"]:
      return c.do_address()
   case set["d"]:
      return c.Job.DownloadDash(
         c.Dash.Body, c.Dash.Url, c.dash_id, c.Playlist.PlayReady,
      )
   }
   return maya.Usage([][]string{
      {"P"},
      {"e", "p"},
      {"a"},
      {"d"},
   })
}

func (c *client) do_email_password() error {
   var err error
   c.Session, err = hulu.FetchSession(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "flag"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/hulu.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c, true)
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 2
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   // 4
   flag.StringVar(&c.Job.PlayReady, "P", c.Job.PlayReady, "PlayReady")
   // 5
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   set := maya.Parse()
   err = maya.SetProxy(c.proxy, "*.mp4,*.mp4a")
   if err != nil {
      return err
   }
   if set["e"] {
      if set["p"] {
         return c.do_email_password()
      }
   }
   if set["a"] {
      return c.do_address()
   }
   if set["P"] {
      return cache.Write(c)
   }
   if set["d"] {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"x"},
      {"e", "p"},
      {"a"},
      {"P"},
      {"d"},
   })
}

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playlist.PlayReady,
   )
}

func (c *client) do_email_password() error {
   c.Session = &hulu.Session{}
   err := c.Session.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

func (c *client) do_address() error {
   err := cache.Update(c, func() error {
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
      return err
   })
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
   proxy string
   // 2
   email    string
   password string
   // 3
   address string
   // 4
   Job maya.Job
   // 5
   dash_id string
}

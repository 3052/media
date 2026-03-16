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
   cache.Read(c)
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.Job.PlayReady, "P", c.Job.PlayReady, "PlayReady")
   // 4
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   set := maya.Parse()
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
      {"e", "p"},
      {"a"},
      {"P"},
      {"d"},
   })
}

func (c *client) do_email_password() error {
   c.Session = &hulu.Session{}
   err := c.Session.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

type client struct {
   Dash     *hulu.Dash
   Playlist *hulu.Playlist
   Session  *hulu.Session
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

func (c *client) do_address() error {
   if cache.Error != nil {
      return cache.Error
   }
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

func (c *client) do_dash_id() error {
   if cache.Error != nil {
      return cache.Error
   }
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playlist.PlayReady,
   )
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

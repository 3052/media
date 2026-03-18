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
   read_err := cache.Read(c)
   // 1
   playReady := maya.StringVar(&c.Job.PlayReady, "P", "PlayReady")
   // 2
   email := maya.StringVar(&c.email, "e", "email")
   password := maya.StringVar(&c.password, "p", "password")
   // 3
   address := maya.StringVar(&c.address, "a", "address")
   // 4
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case len(set) == 0:
      return maya.Usage([][]*flag.Flag{
         {playReady},
         {email, password},
         {address},
         {dash_id},
      })
   case set[playReady.Name]:
      return cache.Write(c)
   case set[email.Name] && set[password.Name]:
      return c.do_email_password()
   case read_err != nil:
      return read_err
   case set[address.Name]:
      return c.do_address()
   case set[dash_id.Name]:
      return c.Job.DownloadDash(
         c.Dash.Body, c.Dash.Url, c.dash_id, c.Playlist.PlayReady,
      )
   }
   return nil
}

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

func (c *client) do_email_password() error {
   var err error
   c.Session, err = hulu.FetchSession(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

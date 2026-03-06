package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "flag"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4,*.mp4a")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

var job maya.PlayReadyJob

func (c *client) do() error {
   job.CertificateChain, _ = maya.ResolveCache("SL2000/CertificateChain")
   job.EncryptSignKey, _ = maya.ResolveCache("SL2000/EncryptSignKey")
   err := cache.Setup("rosso/hulu.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.email, "E", "", "email")
   flag.StringVar(&c.password, "P", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.CertificateChain, "c", job.CertificateChain, "certificate chain")
   flag.StringVar(&job.EncryptSignKey, "e", job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"E", "P"},
      {"a"},
      {"d", "c", "e"},
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
   dash_id string
}

///

func (c *client) do_address() error {
   id, err := hulu.Id(c.address)
   if err != nil {
      return err
   }
   err = cache.Update(c, func() error {
      err := c.Session.TokenRefresh()
      if err != nil {
         return err
      }
      deep_link, err := c.Session.DeepLink(id)
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

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = c.Playlist.PlayReady
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
}

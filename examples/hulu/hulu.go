package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *command) run() error {
   c.cache.Init("SL2000")
   c.job.CertificateChain = c.cache.Join("CertificateChain")
   c.job.EncryptSignKey = c.cache.Join("EncryptSignKey")
   c.cache.Init("hulu")
   // 1
   flag.StringVar(&c.email, "E", "", "email")
   flag.StringVar(&c.password, "P", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.CertificateChain, "c", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "e", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"E", "P"},
      {"a", "x"},
      {"d", "c", "e"},
   })
}

type command struct {
   cache maya.Cache
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash string
   job  maya.PlayReadyJob
}

func (c *command) do_dash() error {
   var playlist hulu.Playlist
   err := c.cache.Get("Playlist", &playlist)
   if err != nil {
      return err
   }
   c.job.Send = playlist.PlayReady
   var dash hulu.Dash
   err = c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func (c *command) do_address() error {
   var session hulu.Session
   err := c.cache.Get("Session", &session)
   if err != nil {
      return err
   }
   err = session.TokenRefresh()
   if err != nil {
      return err
   }
   err = c.cache.Set("Session", session)
   if err != nil {
      return err
   }
   id, err := hulu.Id(c.address)
   if err != nil {
      return err
   }
   deep_link, err := session.DeepLink(id)
   if err != nil {
      return err
   }
   playlist, err := session.Playlist(deep_link)
   if err != nil {
      return err
   }
   err = c.cache.Set("Playlist", playlist)
   if err != nil {
      return err
   }
   dash, err := playlist.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func (c *command) do_email_password() error {
   var session hulu.Session
   err := session.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set("Session", session)
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      switch path.Ext(req.URL.Path) {
      case ".mp4", ".mp4a":
         return "", false
      }
      return "", true
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

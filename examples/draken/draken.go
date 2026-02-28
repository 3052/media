package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/draken"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *command) do_dash() error {
   var login draken.Login
   err := c.cache.Get("Login", &login)
   if err != nil {
      return err
   }
   var playback draken.Playback
   err = c.cache.Get("Playback", &playback)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return login.Widevine(&playback, data)
   }
   var dash draken.Dash
   err = c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".m4s"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}
func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("draken")
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
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
      {"e", "p"},
      {"a"},
      {"d", "C", "P"},
   })
}

func (c *command) do_email_password() error {
   var login draken.Login
   err := login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set("Login", login)
}

func (c *command) do_address() error {
   movie, err := draken.FetchMovie(path.Base(c.address))
   if err != nil {
      return err
   }
   var login draken.Login
   err = c.cache.Get("Login", &login)
   if err != nil {
      return err
   }
   entitlement, err := login.Entitlement(movie)
   if err != nil {
      return err
   }
   playback, err := login.Playback(movie, entitlement)
   if err != nil {
      return err
   }
   err = c.cache.Set("Playback", playback)
   if err != nil {
      return err
   }
   dash, err := playback.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
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
   job  maya.WidevineJob
}

package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/kanopy"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("kanopy")
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.IntVar(&c.kanopy, "k", 0, "Kanopy ID")
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
   if c.kanopy >= 1 {
      return c.do_kanopy()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"k"},
      {"d", "C", "P"},
   })
}

type command struct {
   cache maya.Cache
   // 1
   email    string
   password string
   // 2
   kanopy int
   // 3
   dash string
   job  maya.WidevineJob
}

func (c *command) do_email_password() error {
   var login kanopy.Login
   err := login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set("Login", login)
}
func (c *command) do_dash() error {
   var login kanopy.Login
   err := c.cache.Get("Login", &login)
   if err != nil {
      return err
   }
   var manifest kanopy.PlayManifest
   err = c.cache.Get("PlayManifest", &manifest)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return login.Widevine(&manifest, data)
   }
   var dash kanopy.Dash
   err = c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func (c *command) do_kanopy() error {
   var login kanopy.Login
   err := c.cache.Get("Login", &login)
   if err != nil {
      return err
   }
   member, err := login.Membership()
   if err != nil {
      return err
   }
   plays, err := login.Plays(member, c.kanopy)
   if err != nil {
      return err
   }
   manifest, err := plays.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("PlayManifest", manifest)
   if err != nil {
      return err
   }
   dash, err := manifest.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
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

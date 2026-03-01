package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *client) do() error {
   c.job.CertificateChain, _ = maya.ResolveCache("SL2000/CertificateChain")
   c.job.EncryptSignKey, _ = maya.ResolveCache("SL2000/EncryptSignKey")
   err := c.cache.Init("rosso/hulu.xml")
   if err != nil {
      return err
   }
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
      {"a"},
      {"d", "c", "e"},
   })
}

type client struct {
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

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = state.Playlist.PlayReady
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

type saved_state struct {
   Dash     *hulu.Dash
   Playlist *hulu.Playlist
   Session  *hulu.Session
}

func (c *client) do_address() error {
   id, err := hulu.Id(c.address)
   if err != nil {
      return err
   }
   var state saved_state
   err = c.cache.Update(&state, func() error {
      err := state.Session.TokenRefresh()
      if err != nil {
         return err
      }
      deep_link, err := state.Session.DeepLink(id)
      if err != nil {
         return err
      }
      state.Playlist, err = state.Session.Playlist(deep_link)
      if err != nil {
         return err
      }
      state.Dash, err = state.Playlist.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

func (c *client) do_email_password() error {
   var session hulu.Session
   err := session.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set(saved_state{Session: &session})
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      switch path.Ext(req.URL.Path) {
      case ".mp4", ".mp4a":
         return "", false
      }
      return "", true
   })
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "flag"
   "log"
   "net/http"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s,*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   c.job.CertificateChain, _ = maya.ResolveCache("SL2000/CertificateChain")
   c.job.EncryptSignKey, _ = maya.ResolveCache("SL2000/EncryptSignKey")
   err := c.cache.Setup("rosso/paramount.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.username, "U", "", "username")
   flag.StringVar(&c.password, "P", "", "password")
   // 2
   flag.StringVar(&c.paramount, "p", "", "paramount ID")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.BoolVar(&c.cookie, "c", false, "cookie")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   if c.username != "" {
      if c.password != "" {
         return c.do_username_password()
      }
   }
   if c.paramount != "" {
      return c.do_paramount()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"U", "P"},
      {"p"},
      {"d", "c", "C", "E"},
   })
}

func (c *client) do_dash() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   var state saved_state
   err = c.cache.Read(&state)
   if err != nil {
      return err
   }
   if !c.cookie {
      state.Cookie = nil
   }
   token, err := paramount.PlayReady(at, state.ContentId, state.Cookie)
   if err != nil {
      return err
   }
   c.job.Send = token.Send
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

type saved_state struct {
   ContentId string
   Cookie    *http.Cookie
   Dash      *paramount.Dash
}

func (c *client) do_username_password() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   cookie, err := paramount.Login(at, c.username, c.password)
   if err != nil {
      return err
   }
   return c.cache.Write(saved_state{Cookie: cookie})
}

func (c *client) do_paramount() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   item, err := paramount.FetchItem(at, c.paramount)
   if err != nil {
      return err
   }
   var state saved_state
   state.Dash, err = item.Dash()
   if err != nil {
      return err
   }
   state.ContentId = c.paramount
   err = c.cache.Write(state)
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

type client struct {
   cache maya.Cache
   // 1
   username string
   password string
   // 2
   paramount string
   // 3
   dash   string
   cookie bool
   job    maya.PlayReadyJob
}

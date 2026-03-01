package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/criterion"
   "flag"
   "log"
   "net/http"
   "path"
)

type saved_state struct {
   Dash      *criterion.Dash
   MediaFile *criterion.MediaFile
   Token     *criterion.Token
}

func (c *client) do_address() error {
   var state saved_state
   err := c.cache.Update(&state, func() error {
      err := state.Token.Refresh()
      if err != nil {
         return err
      }
      item, err := state.Token.Item(path.Base(c.address))
      if err != nil {
         return err
      }
      files, err := state.Token.Files(item)
      if err != nil {
         return err
      }
      state.MediaFile, err = files.Dash()
      if err != nil {
         return err
      }
      state.Dash, err = state.MediaFile.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Init("rosso/criterion.xml")
   if err != nil {
      return err
   }
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

type client struct {
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

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".mp4"
   })
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email_password() error {
   var token criterion.Token
   err := token.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set(saved_state{Token: &token})
}

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = state.MediaFile.Widevine
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

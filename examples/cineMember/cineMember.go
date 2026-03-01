package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/cineMember"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".m4s"
   })
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   err := c.cache.Init("rosso/cineMember.xml")
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
      {"d"},
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
   job  maya.Job
}

func (c *client) do_email_password() error {
   cookie, err := cineMember.FetchSession()
   if err != nil {
      return err
   }
   err = cineMember.FetchLogin(cookie, c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set(saved_state{Cookie: cookie})
}

type saved_state struct {
   Cookie *http.Cookie
   Dash   *cineMember.Dash
}

func (c *client) do_address() error {
   id, err := cineMember.FetchId(c.address)
   if err != nil {
      return err
   }
   var state saved_state
   err = c.cache.Update(&state, func() error {
      stream, err := cineMember.FetchStream(state.Cookie, id)
      if err != nil {
         return err
      }
      link, err := stream.Dash()
      if err != nil {
         return err
      }
      state.Dash, err = link.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

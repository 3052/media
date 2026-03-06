package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/cineMember"
   "flag"
   "log"
   "net/http"
)

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

var job maya.Job

func (c *client) do() error {
   err := cache.Setup("rosso/cineMember.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
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
      {"e", "p"},
      {"a"},
      {"d"},
   })
}

func (c *client) do_email_password() error {
   var err error
   c.Cookie, err = cineMember.FetchSession()
   if err != nil {
      return err
   }
   err = cineMember.FetchLogin(c.Cookie, c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

type client struct {
   Cookie *http.Cookie
   Dash   *cineMember.Dash
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash_id string
}

func (c *client) do_address() error {
   id, err := cineMember.FetchId(c.address)
   if err != nil {
      return err
   }
   err = cache.Update(c, func() error {
      stream, err := cineMember.FetchStream(c.Cookie, id)
      if err != nil {
         return err
      }
      link, err := stream.Dash()
      if err != nil {
         return err
      }
      c.Dash, err = link.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

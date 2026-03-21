package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/cineMember"
   "flag"
   "log"
   "net/http"
)

func (c *client) do() error {
   err := cache.Setup("rosso/cineMember.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   email := maya.StringVar(&c.email, "e", "email")
   password := maya.StringVar(&c.password, "p", "password")
   //------------------------------------------------------
   address := maya.StringVar(&c.address, "a", "address")
   //---------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   if set[email] {
      if set[password] {
         return c.do_email_password()
      }
   }
   if set[address] {
      return with_cache(c.do_address)
   }
   if set[dash_id] {
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {email, password},
      {address},
      {dash_id},
   })
}

func (c *client) do_address() error {
   id, err := cineMember.FetchId(c.address)
   if err != nil {
      return err
   }
   stream, err := cineMember.FetchStream(c.Session, id)
   if err != nil {
      return err
   }
   link, err := stream.Dash()
   if err != nil {
      return err
   }
   c.Dash, err = link.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, nil)
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

func (c *client) do_email_password() error {
   var err error
   c.Session, err = cineMember.FetchSession()
   if err != nil {
      return err
   }
   err = cineMember.FetchLogin(c.Session, c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

type client struct {
   Dash   *cineMember.Dash
   Session *http.Cookie
   //---------------------
   Job maya.Job
   //-------------
   email    string
   password string
   //-------------
   address string
   //------------
   dash_id string
}

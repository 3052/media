package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/cineMember"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".m4s"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) do_email_password() error {
   var session cineMember.Session
   err := session.Fetch()
   if err != nil {
      return err
   }
   err = session.Login(c.email, c.password)
   if err != nil {
      return err
   }
   return maya.Write(c.name, &user_cache{Session: &session})
}

type user_cache struct {
   Dash    *cineMember.Dash
   Session *cineMember.Session
}

type command struct {
   name string
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash string
   job  maya.Job
}

func (c *command) do_address() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   id, err := cineMember.FetchId(c.address)
   if err != nil {
      return err
   }
   stream, err := cache.Session.Stream(id)
   if err != nil {
      return err
   }
   link, err := stream.Dash()
   if err != nil {
      return err
   }
   cache.Dash, err = link.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   c.name = filepath.ToSlash(cache) + "/cineMember/userCache.xml"
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.Parse()
   // 1
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   // 2
   if c.address != "" {
      return c.do_address()
   }
   // 3
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"a"},
      {"d"},
   })
}

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

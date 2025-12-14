package main

import (
   "41.neocities.org/media/cineMember"
   "41.neocities.org/net"
   "encoding/xml"
   "errors"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
         return ""
      }
      return "L"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)

   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
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
   flag.Usage()
   return nil
}

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
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
   return write(c.name, &user_cache{Session: &session})
}

func read(name string) (*user_cache, error) {
   data, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   cache := &user_cache{}
   err = xml.Unmarshal(data, cache)
   if err != nil {
      return nil, err
   }
   return cache, nil
}

func (c *command) do_address() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   id, err := cineMember.Id(c.address)
   if err != nil {
      return err
   }
   stream, err := cache.Session.Stream(id)
   if err != nil {
      return err
   }
   link, ok := stream.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   cache.Mpd, cache.MpdBody, err = link.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return net.Representations(cache.Mpd, cache.MpdBody)
}

type command struct {
   name   string
   config net.Config
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash string
}

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   return c.config.Download(cache.Mpd, cache.MpdBody, c.dash)
}

type user_cache struct {
   Mpd     *url.URL
   MpdBody []byte
   Session *cineMember.Session
}

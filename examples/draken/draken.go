package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/draken"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) do_dash() error {
   c.job.Send = func(data []byte) ([]byte, error) {
      return login.Widevine(play, data)
   }
   return c.filters.Filter(resp, &c.job)
}

type user_cache struct {
   Dash *draken.Dash
   Login *draken.Login
   Playback *draken.Playback
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(*http.Request) string {
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
   c.name = cache + "/draken/userCache.xml"
   c.job.ClientId = filepath.Join(cache, "/L3/client_id.bin")
   c.job.PrivateKey = filepath.Join(cache, "/L3/private_key.pem")
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
   maya.Usage([][]string{
      {"e", "p"},
      {"a"},
      {"d", "C", "P"},
   })
   return nil
}

func (c *command) do_email_password() error {
   var login draken.Login
   err := login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return write(c.name, &user_cache{Login: &login})
}

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
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
   var movie draken.Movie
   err = movie.Fetch(path.Base(c.address))
   if err != nil {
      return err
   }
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   entitlement, err := cache.Login.Entitlement(movie)
   if err != nil {
      return err
   }
   cache.Playback, err = cache.Login.Playback(&movie, entitlement)
   if err != nil {
      return err
   }
   cache.Dash, err = cache.Playback.Dash()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

type command struct {
   name string
   // 1
   email    string
   password string
   // 2
   address  string
   // 3
   dash string
   job   maya.WidevineJob
}

package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/draken"
   "encoding/json"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(*http.Request) string {
      return "LP"
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
   c.name = cache + "/draken/userCache.json"
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
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
   data, err := json.Marshal(cache)
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
   err = json.Unmarshal(data, cache)
   if err != nil {
      return nil, err
   }
   return cache, nil
}
func (c *command) do_address() error {
   movie, err := draken.FetchMovie(path.Base(c.address))
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
   cache.Playback, err = cache.Login.Playback(movie, entitlement)
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
   address string
   // 3
   dash string
   job  maya.WidevineJob
}

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Login.Widevine(cache.Playback, data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

type user_cache struct {
   Dash     *draken.Dash
   Login    *draken.Login
   Playback *draken.Playback
}

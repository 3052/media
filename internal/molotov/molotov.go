package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/molotov"
   "encoding/json"
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
   maya.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
         return ""
      }
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
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/molotov/user_cache.json"

   flag.StringVar(&c.config.ClientId, "C", c.config.ClientId, "client ID")
   flag.StringVar(&c.config.PrivateKey, "P", c.config.PrivateKey, "private key")
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
   data, err := json.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (c *command) do_email_password() error {
   var login molotov.Login
   err := login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return write(c.name, &user_cache{Login: &login})
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

type command struct {
   address  string
   config   maya.Config
   dash     string
   email    string
   name     string
   password string
}

func (c *command) do_address() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   err = cache.Login.Refresh()
   if err != nil {
      return err
   }
   var media molotov.MediaId
   err = media.Parse(c.address)
   if err != nil {
      return err
   }
   view, err := cache.Login.ProgramView(&media)
   if err != nil {
      return err
   }
   cache.Playback, err = cache.Login.Playback(view)
   if err != nil {
      return err
   }
   cache.Mpd, cache.MpdBody, err = cache.Playback.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.Representations(cache.Mpd, cache.MpdBody)
}
func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      return cache.Playback.Widevine(data)
   }
   return c.config.Download(cache.Mpd, cache.MpdBody, c.dash)
}

type user_cache struct {
   Login    *molotov.Login
   Mpd      *url.URL
   MpdBody  []byte
   Playback *molotov.Playback
}

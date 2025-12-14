package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/molotov"
   "flag"
   "log"
   "net/http"
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
      return "L"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) New() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/molotov/user_cache.xml"

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
   var login molotiv.Login
   err := login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   err = login.Refresh()
   if err != nil {
      return err
   }
   return write(c.name, &user_cache{Login: &login})
}

type command struct {
   name   string
   config maya.Config
   // 1
   email    string
   password string
   
   // 2
   address string
   // 3
   dash string
}

type user_cache struct {
   Login *Login
   Mpd *url.URL
   MpdBody []byte
}

///

func (c *command) do_address() error {
   data, err := os.ReadFile(c.name + "/molotov/Login")
   if err != nil {
      return err
   }
   var login molotov.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = login.Refresh()
   if err != nil {
      return err
   }
   err = write_file(c.name+"/molotov/Login", data)
   if err != nil {
      return err
   }
   var media molotov.MediaId
   err = media.Parse(c.address)
   if err != nil {
      return err
   }
   play_url, err := login.PlayUrl(&media)
   if err != nil {
      return err
   }
   playback, err := login.Playback(play_url)
   if err != nil {
      return err
   }
   resp, err := http.Get(playback.FhdReady())
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      return playback.Widevine(data)
   }
   return c.filters.Filter(resp, &c.config)
}

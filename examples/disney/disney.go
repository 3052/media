package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/disney"
   "encoding/xml"
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
   c.name = cache + "/disney/userCache.xml"
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"

   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
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

func (c *command) do_email_password() error {
   var device disney.Device
   err := device.Register()
   if err != nil {
      return err
   }
   account_without, err := device.Login(c.email, c.password)
   if err != nil {
      return err
   }
   var cache user_cache
   cache.Account, err = account_without.SwitchProfile()
   if err != nil {
      return err
   }
   return write(c.name, &cache)
}

type user_cache struct {
   Account *disney.Account
}

type command struct {
   name     string
   job   maya.WidevineJob
   // 1
   email    string
   password string
   
   // 2
   address  string
   // 3
   dash     string
}

///

func (c *command) do_address() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   err = cache.Session.TokenRefresh()
   if err != nil {
      return err
   }
   id, err := disney.Id(c.address)
   if err != nil {
      return err
   }
   deep_link, err := cache.Session.DeepLink(id)
   if err != nil {
      return err
   }
   cache.Playlist, err = cache.Session.Playlist(deep_link)
   if err != nil {
      return err
   }
   cache.Mpd, err = cache.Playlist.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.Representations(cache.Mpd.Url, cache.Mpd.Body)
}

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Playlist.Widevine(data)
   }
   return c.job.Download(cache.Mpd.Url, cache.Mpd.Body, c.dash)
}

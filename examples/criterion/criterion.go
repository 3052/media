package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/criterion"
   "encoding/xml"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.name = cache + "/criterion/userCache.xml"
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
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
   var token criterion.Token
   err := token.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return write(c.name, &user_cache{Token: &token})
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

type command struct {
   job  maya.WidevineJob
   name string
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
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.MediaFile.Widevine(data)
   }
   return c.job.DownloadDash(cache.Mpd.Body, cache.Mpd.Url, c.dash)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".mp4" {
         return ""
      }
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}
type user_cache struct {
   MediaFile *criterion.MediaFile
   Mpd       *criterion.Mpd
   Token     *criterion.Token
}

func (c *command) do_address() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   err = cache.Token.Refresh()
   if err != nil {
      return err
   }
   item, err := cache.Token.Item(path.Base(c.address))
   if err != nil {
      return err
   }
   files, err := cache.Token.Files(item)
   if err != nil {
      return err
   }
   var ok bool
   cache.MediaFile, ok = files.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   cache.Mpd, err = cache.MediaFile.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Mpd.Body, cache.Mpd.Url)
}

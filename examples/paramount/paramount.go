package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/paramount"
   "encoding/xml"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(c.app_secret())
   if err != nil {
      return err
   }
   token, err := paramount.PlayReady(at, cache.Item.ContentId, cache.Cookie)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return token.Send(data)
   }
   return c.job.DownloadDash(cache.Mpd.Body, cache.Mpd.Url, c.dash)
}

func (c *command) do_paramount() error {
   at, err := paramount.GetAt(c.app_secret())
   if err != nil {
      return err
   }
   cache, err := read(c.name)
   if err != nil {
      cache = &user_cache{}
   }
   cache.Item, err = paramount.FetchItem(at, c.paramount)
   if err != nil {
      return err
   }
   cache.Mpd, err = cache.Item.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Mpd.Body, cache.Mpd.Url)
}

type command struct {
   job  maya.PlayReadyJob
   name string
   // 1
   username string
   password string
   // 2
   paramount string
   intl      bool
   // 3
   dash string
   cookie bool
}

func (c *command) app_secret() string {
   if c.intl {
      return paramount.ComCbsCa.AppSecret
   }
   return paramount.ComCbsApp.AppSecret
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      switch path.Ext(req.URL.Path) {
      case ".m4s", ".mp4":
         return ""
      }
      switch path.Base(req.URL.Path) {
      case "anonymous-session-token.json", "getlicense":
         return "L"
      }
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type user_cache struct {
   Cookie *http.Cookie
   Item *paramount.Item
   Mpd  *paramount.Mpd
}

func (c *command) do_username_password() error {
   at, err := paramount.GetAt(paramount.ComCbsApp.AppSecret)
   if err != nil {
      return err
   }
   var cache user_cache
   cache.Cookie, err = paramount.Login(at, c.username, c.password)
   if err != nil {
      return err
   }
   return write(c.name, &cache)
}

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.name = cache + "/paramount/userCache.xml"
   c.job.CertificateChain = cache + "/SL2000/CertificateChain"
   c.job.EncryptSignKey = cache + "/SL2000/EncryptSignKey"
   // 1
   flag.StringVar(&c.username, "U", "", "username")
   flag.StringVar(&c.password, "P", "", "password")
   // 2
   flag.StringVar(&c.paramount, "p", "", "paramount ID")
   flag.BoolVar(&c.intl, "i", false, "intl")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.BoolVar(&c.cookie, "g", false, "cookie")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   // 1
   if c.username != "" {
      if c.password != "" {
         return c.do_username_password()
      }
   }
   // 2
   if c.paramount != "" {
      return c.do_paramount()
   }
   // 3
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
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

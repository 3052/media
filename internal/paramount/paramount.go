package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/paramount"
   "encoding/xml"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

type command struct {
   config    maya.Config
   dash      string
   intl      bool
   name      string
   paramount string
}

func (c *command) do_dash() error {
   data, err := os.ReadFile(c.name)
   if err != nil {
      return err
   }
   var cache mpd
   err = xml.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   // INTL does NOT allow anonymous key request, so if you are INTL you
   // will need to use US VPN until someone codes the INTL login
   at, err := paramount.GetAt(paramount.ComCbsApp.AppSecret)
   if err != nil {
      return err
   }
   var token paramount.SessionToken
   err = token.Fetch(at, c.paramount)
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      return token.Widevine(data)
   }
   return c.config.Download(cache.Url, cache.Body, c.dash)
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

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/paramount/user_cache.json"

   flag.StringVar(&c.config.ClientId, "C", c.config.ClientId, "client ID")
   flag.StringVar(&c.config.PrivateKey, "P", c.config.PrivateKey, "private key")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.BoolVar(&c.intl, "i", false, "intl")
   flag.StringVar(&c.paramount, "p", "", "paramount ID")
   flag.Parse()

   if c.paramount != "" {
      return c.do_paramount()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
}

func (c *command) app_secret() string {
   if c.intl {
      return paramount.ComCbsCa.AppSecret
   }
   return paramount.ComCbsApp.AppSecret
}

type mpd struct {
   Body []byte
   Url  *url.URL
}

func (c *command) do_paramount() error {
   at, err := paramount.GetAt(c.app_secret())
   if err != nil {
      return err
   }
   item, err := paramount.FetchItem(at, c.paramount)
   if err != nil {
      return err
   }
   var cache mpd
   cache.Url, cache.Body, err = item.Mpd()
   if err != nil {
      return err
   }
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", c.name)
   err = os.WriteFile(c.name, data, os.ModePerm)
   if err != nil {
      return err
   }
   return maya.Representations(cache.Url, cache.Body)
}

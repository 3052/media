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
   data, err := os.ReadFile(c.name)
   if err != nil {
      return err
   }
   var cache paramount.Mpd
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
   c.job.Send = func(data []byte) ([]byte, error) {
      return token.Widevine(data)
   }
   return c.job.DownloadDash(cache.Body, cache.Url, c.dash)
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

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/paramount/userCache.xml"
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   // 1
   flag.StringVar(&c.paramount, "p", "", "paramount ID")
   flag.BoolVar(&c.intl, "i", false, "intl")
   // 2
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.Parse()
   // 1
   if c.paramount != "" {
      return c.do_paramount()
   }
   // 2
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
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
   cache, err := item.Mpd()
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
   return maya.ListDash(cache.Body, cache.Url)
}

type command struct {
   job    maya.WidevineJob
   name      string
   // 1
   paramount string
   intl      bool
   // 2
   dash      string
}

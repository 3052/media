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

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.name = cache + "/paramount/userCache.xml"
   c.job.CertificateChain = cache + "/SL3000/CertificateChain"
   c.job.EncryptSignKey = cache + "/SL3000/EncryptSignKey"
   // 1
   flag.BoolVar(&c.intl, "i", false, "intl")
   flag.StringVar(&c.paramount, "p", "", "paramount ID")
   // 2
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.IntVar(&c.job.Threads, "t", 2, "threads")
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
   Item *paramount.Item
   Mpd  *paramount.Mpd
}

func (c *command) app_secret() string {
   if c.intl {
      return paramount.ComCbsCa.AppSecret
   }
   return paramount.ComCbsApp.AppSecret
}

type command struct {
   job  maya.PlayReadyJob
   name string
   // 1
   paramount string
   intl      bool
   // 2
   dash string
}

func (c *command) do_paramount() error {
   at, err := paramount.GetAt(c.app_secret())
   if err != nil {
      return err
   }
   var cache user_cache
   cache.Item, err = paramount.FetchItem(at, c.paramount)
   if err != nil {
      return err
   }
   cache.Mpd, err = cache.Item.Mpd()
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
   return maya.ListDash(cache.Mpd.Body, cache.Mpd.Url)
}

func (c *command) do_dash() error {
   data, err := os.ReadFile(c.name)
   if err != nil {
      return err
   }
   var cache user_cache
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
   token, err := paramount.PlayReady(at, cache.Item.ContentId)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return token.Send(data)
   }
   return c.job.DownloadDash(cache.Mpd.Body, cache.Mpd.Url, c.dash)
}

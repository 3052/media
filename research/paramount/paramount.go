package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/paramount"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(paramount.AppSecrets[0].ComCbsApp)
   if err != nil {
      return err
   }
   token, err := paramount.PlayReady(at, cache.ContentId, nil)
   if err != nil {
      return err
   }
   c.job.Send = token.Send
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

func main() {
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   name string
   // 1
   paramount string
   // 2
   dash string
   job  maya.PlayReadyJob
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
   flag.StringVar(&c.paramount, "p", "", "paramount ID")
   // 2
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.CertificateChain, "c", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "e", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   maya.SetProxy(func(req *http.Request) (string, bool) {
      switch path.Ext(req.URL.Path) {
      case ".m4s", ".mp4":
         return "", false
      }
      return "", true
   })
   if c.paramount != "" {
      return c.do_paramount()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"p"},
      {"d", "c", "e"},
   })
}

func (c *command) do_paramount() error {
   at, err := paramount.GetAt(paramount.AppSecrets[0].ComCbsApp)
   if err != nil {
      return err
   }
   item, err := paramount.FetchItem(at, c.paramount)
   if err != nil {
      return err
   }
   var cache user_cache
   cache.Dash, err = item.Dash()
   if err != nil {
      return err
   }
   cache.ContentId = c.paramount
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

type user_cache struct {
   ContentId string
   Dash      *paramount.Dash
}

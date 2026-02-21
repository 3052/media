package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/disney"
   "flag"
   "fmt"
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
   c.name = cache + "/disney/userCache.xml"
   c.job.CertificateChain = cache + "/SL3000/CertificateChain"
   c.job.EncryptSignKey = cache + "/SL3000/EncryptSignKey"
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.season, "s", "", "season")
   // 4
   flag.StringVar(&c.media_id, "m", "", "media ID")
   // 5
   flag.StringVar(&c.hls, "h", "", "HLS ID")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   maya.SetProxy(func(req *http.Request) (string, bool) {
      switch path.Ext(req.URL.Path) {
      case ".mp4", ".mp4a":
         return "", false
      }
      return c.proxy, true
   })
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.season != "" {
      return c.do_season()
   }
   if c.media_id != "" {
      return c.do_media_id()
   }
   if c.hls != "" {
      return c.do_hls()
   }
   return maya.Usage([][]string{
      {"e", "p", "x"},
      {"a"},
      {"s"},
      {"m"},
      {"h", "C", "E"},
   })
}

func (c *command) do_email_password() error {
   device, err := disney.RegisterDevice()
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
   return maya.Write(c.name, &cache)
}

func (c *command) do_address() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   entity, err := disney.GetEntity(c.address)
   if err != nil {
      return err
   }
   page, err := cache.Account.Page(entity)
   if err != nil {
      return err
   }
   fmt.Println(page)
   return nil
}

func (c *command) do_season() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   season, err := cache.Account.Season(c.season)
   if err != nil {
      return err
   }
   fmt.Println(season)
   return nil
}

func (c *command) do_media_id() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   stream, err := cache.Account.Stream(c.media_id)
   if err != nil {
      return err
   }
   cache.Hls, err = stream.Hls()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListHls(cache.Hls.Body, cache.Hls.Url)
}

func main() {
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) do_hls() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = cache.Account.PlayReady
   return c.job.DownloadHls(cache.Hls.Body, cache.Hls.Url, c.hls)
}

type user_cache struct {
   Account *disney.Account
   Hls     *disney.Hls
}

type command struct {
   name string
   // 1
   email    string
   password string
   proxy    string
   // 2
   address string
   // 3
   season string
   // 4
   media_id string
   // 5
   hls string
   job maya.PlayReadyJob
}

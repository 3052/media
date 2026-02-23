package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/hboMax"
   "flag"
   "fmt"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

type command struct {
   name string
   // 1
   initiate bool
   market   string
   // 2
   login bool
   // 3
   address string
   season  int
   proxy   string
   // 4
   edit string
   // 5
   dash string
   job  maya.PlayReadyJob
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.name = cache + "/hboMax/userCache.xml"
   c.job.CertificateChain = cache + "/SL3000/CertificateChain"
   c.job.EncryptSignKey = cache + "/SL3000/EncryptSignKey"
   // 1
   flag.BoolVar(&c.initiate, "i", false, "device initiate")
   flag.StringVar(
      &c.market, "m", hboMax.Markets[0], fmt.Sprint(hboMax.Markets),
   )
   // 2
   flag.BoolVar(&c.login, "l", false, "device login")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   flag.IntVar(&c.season, "s", 0, "season")
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 4
   flag.StringVar(&c.edit, "e", "", "edit ID")
   // 5
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   err = c.do_proxy()
   if err != nil {
      return err
   }
   if c.initiate {
      return c.do_initiate()
   }
   if c.login {
      return c.do_login()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.edit != "" {
      return c.do_edit()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"i", "m"},
      {"l"},
      {"a", "s", "x"},
      {"e"},
      {"d", "C", "E"},
   })
}

func (c *command) do_proxy() error {
   if c.edit != "" {
      cache, err := maya.Read[user_cache](c.name)
      if err != nil {
         return err
      }
      c.proxy = cache.Proxy
   }
   maya.SetProxy(func(req *http.Request) (string, bool) {
      if path.Ext(req.URL.Path) == ".mp4" {
         return "", false
      }
      return c.proxy, true
   })
   return nil
}

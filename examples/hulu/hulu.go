package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/hulu"
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
      switch path.Ext(req.URL.Path) {
      case ".mp4", ".mp4a":
         return ""
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
   c.job.CertificateChain = cache + "/SL2000/CertificateChain"
   c.job.EncryptSignKey = cache + "/SL2000/EncryptSignKey"
   c.name = cache + "/hulu/userCache.xml"
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
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
   return maya.Usage([][]string{
      {"e", "p"},
      {"a"},
      {"d", "C", "E"},
   })
}

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Playlist.PlayReady(data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

func (c *command) do_address() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   err = cache.Session.TokenRefresh()
   if err != nil {
      return err
   }
   id, err := hulu.Id(c.address)
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
   cache.Dash, err = cache.Playlist.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

type user_cache struct {
   Dash     *hulu.Dash
   Playlist *hulu.Playlist
   Session  *hulu.Session
}

func (c *command) do_email_password() error {
   var session hulu.Session
   err := session.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return maya.Write(c.name, &user_cache{Session: &session})
}

type command struct {
   name string
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash string
   job  maya.PlayReadyJob
}

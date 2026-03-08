package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "flag"
   "fmt"
   "log"
   "net/http"
)

type client struct {
   Dash     *hboMax.Dash
   Login    *hboMax.Login
   Playback *hboMax.Playback
   St       *http.Cookie
   // 1
   Proxy       string
   proxy_write bool
   // 2
   initiate bool
   market   string
   // 3
   login bool
   // 4
   address string
   season  int
   // 5
   edit string
   // 6
   dash_id string
}

func (c *client) do() error {
   job.CertificateChain, _ = maya.ResolveCache("SL3000/CertificateChain")
   job.EncryptSignKey, _ = maya.ResolveCache("SL3000/EncryptSignKey")
   err := cache.Setup("rosso/hboMax.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c, true)
   if err != nil {
      return err
   }
   // 1
   flag.Func("x", "proxy", func(proxy string) error {
      c.Proxy = proxy
      c.proxy_write = true
      return nil
   })
   // 2
   flag.BoolVar(&c.initiate, "i", false, "device initiate")
   flag.StringVar(
      &c.market, "m", hboMax.Markets[0], fmt.Sprint(hboMax.Markets),
   )
   // 3
   flag.BoolVar(&c.login, "l", false, "device login")
   // 4
   flag.StringVar(&c.address, "a", "", "address")
   flag.IntVar(&c.season, "s", 0, "season")
   // 5
   flag.StringVar(&c.edit, "e", "", "edit ID")
   // 6
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.CertificateChain, "C", job.CertificateChain, "certificate chain")
   flag.StringVar(&job.EncryptSignKey, "E", job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   err = maya.SetProxy(c.Proxy, "*.mp4")
   if err != nil {
      return err
   }
   if c.proxy_write {
      return cache.Write(c)
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
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"x"},
      {"i", "m"},
      {"l"},
      {"a", "s"},
      {"e"},
      {"d", "C", "E"},
   })
}

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = c.Playback.PlayReady
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
}

var job maya.PlayReadyJob

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_initiate() error {
   return cache.Update(c, func() error {
      var err error
      c.St, err = hboMax.FetchSt()
      if err != nil {
         return err
      }
      initiate, err := hboMax.FetchInitiate(c.St, c.market)
      if err != nil {
         return err
      }
      fmt.Println(initiate)
      return nil
   }, true)
}

func (c *client) do_login() error {
   return cache.Update(c, func() error {
      c.Login = &hboMax.Login{}
      return c.Login.Fetch(c.St)
   })
}

func (c *client) do_address() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   show, err := hboMax.ParseShow(c.address)
   if err != nil {
      return err
   }
   var videos *hboMax.Videos
   if c.season >= 1 {
      videos, err = c.Login.Season(show, c.season)
   } else {
      videos, err = c.Login.Movie(show)
   }
   if err != nil {
      return err
   }
   videos.FilterAndSort()
   for i, video := range videos.Included {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(video)
   }
   return nil
}

func (c *client) do_edit() error {
   err := cache.Update(c, func() error {
      var err error
      c.Playback, err = c.Login.PlayReady(c.edit)
      if err != nil {
         return err
      }
      c.Dash, err = c.Playback.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

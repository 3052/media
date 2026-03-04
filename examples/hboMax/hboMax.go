package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "flag"
   "fmt"
   "log"
   "net/http"
)

func (c *client) do_proxy() error {
   ok, err := cache.Read(&state)
   if !ok {
      return err
   }
   if c.proxy != nil {
      state.Proxy = *c.proxy
      err = cache.Write(state)
      if err != nil {
         return err
      }
   }
   return maya.SetProxy(state.Proxy, "*.mp4")
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var job maya.PlayReadyJob

var cache maya.Cache

func (c *client) do() error {
   job.CertificateChain, _ = maya.ResolveCache("SL3000/CertificateChain")
   job.EncryptSignKey, _ = maya.ResolveCache("SL3000/EncryptSignKey")
   err := cache.Setup("rosso/hboMax.xml")
   if err != nil {
      return err
   }
   // 1
   flag.Func("x", "proxy", func(s string) error {
      c.proxy = &s
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
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&job.CertificateChain, "C", job.CertificateChain, "certificate chain")
   flag.StringVar(&job.EncryptSignKey, "E", job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   err = c.do_proxy()
   if err != nil {
      return err
   }
   if c.proxy != nil {
      return nil
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
      {"x"},
      {"i", "m"},
      {"l"},
      {"a", "s"},
      {"e"},
      {"d", "C", "E"},
   })
}

func (c *client) do_initiate() error {
   var err error
   state.St, err = hboMax.FetchSt()
   if err != nil {
      return err
   }
   initiate, err := hboMax.FetchInitiate(state.St, c.market)
   if err != nil {
      return err
   }
   fmt.Println(initiate)
   return cache.Write(state)
}

type client struct {
   // 1
   proxy *string
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
   dash string
}

var state struct {
   Dash     *hboMax.Dash
   Login    *hboMax.Login
   Playback *hboMax.Playback
   Proxy    string
   St       *http.Cookie
}

func (c *client) do_dash() error {
   if state.Playback == nil {
      _, err := cache.Read(&state)
      if err != nil {
         return err
      }
   }
   job.Send = state.Playback.PlayReady
   return job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

func (c *client) do_address() error {
   show, err := hboMax.ParseUrl(c.address)
   if err != nil {
      return err
   }
   if state.Login == nil {
      _, err = cache.Read(&state)
      if err != nil {
         return err
      }
   }
   var videos *hboMax.Videos
   if c.season >= 1 {
      videos, err = state.Login.Season(show, c.season)
   } else {
      videos, err = state.Login.Movie(show)
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

func (c *client) do_login() error {
   if state.St == nil {
      _, err := cache.Read(&state)
      if err != nil {
         return err
      }
   }
   var err error
   state.Login, err = hboMax.FetchLogin(state.St)
   if err != nil {
      return err
   }
   return cache.Write(state)
}

func (c *client) do_edit() error {
   if state.Login == nil {
      _, err := cache.Read(&state)
      if err != nil {
         return err
      }
   }
   var err error
   state.Playback, err = state.Login.PlayReady(c.edit)
   if err != nil {
      return err
   }
   state.Dash, err = state.Playback.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(state)
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

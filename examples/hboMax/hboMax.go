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
   c.cache.Optional = true
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   if c.proxy != nil {
      state.Proxy = *c.proxy
      err = c.cache.Set(state)
      if err != nil {
         return err
      }
   }
   return maya.SetProxy(state.Proxy, "*.mp4")
}

func (c *client) do() error {
   c.job.CertificateChain, _ = maya.ResolveCache("SL3000/CertificateChain")
   c.job.EncryptSignKey, _ = maya.ResolveCache("SL3000/EncryptSignKey")
   err := c.cache.Init("rosso/hboMax.xml")
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
   flag.IntVar(&c.job.Threads, "t", 2, "threads")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
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
      {"d", "t", "C", "E"},
   })
}

type client struct {
   cache maya.Cache
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
   job  maya.PlayReadyJob
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type saved_state struct {
   Dash     *hboMax.Dash
   Login    *hboMax.Login
   Playback *hboMax.Playback
   Proxy    string
   St       *http.Cookie
}

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Get(&state)
   if err != nil {
      return err
   }
   c.job.Send = state.Playback.PlayReady
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

func (c *client) do_edit() error {
   var state saved_state
   err := c.cache.Update(&state, func() error {
      var err error
      state.Playback, err = state.Login.PlayReady(c.edit)
      if err != nil {
         return err
      }
      state.Dash, err = state.Playback.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

func (c *client) do_initiate() error {
   st, err := hboMax.FetchSt()
   if err != nil {
      return err
   }
   initiate, err := hboMax.FetchInitiate(st, c.market)
   if err != nil {
      return err
   }
   fmt.Println(initiate)
   return c.cache.Set(saved_state{St: st})
}

func (c *client) do_login() error {
   var state saved_state
   return c.cache.Update(&state, func() error {
      var err error
      state.Login, err = hboMax.FetchLogin(state.St)
      return err
   })
}

func (c *client) do_address() error {
   show, err := hboMax.ParseUrl(c.address)
   if err != nil {
      return err
   }
   var state saved_state
   err = c.cache.Get(&state)
   if err != nil {
      return err
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

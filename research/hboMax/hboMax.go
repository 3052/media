package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "flag"
   "fmt"
   "net/http"
)

func (c *client) do() error {
   c.job.CertificateChain, _ = maya.ResolveCache("SL2000/CertificateChain")
   c.job.EncryptSignKey, _ = maya.ResolveCache("SL2000/EncryptSignKey")
   err := c.cache.Init("rosso/hboMax.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.proxy, "x", "", "proxy")
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
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   err = maya.SetProxy(c.proxy, "*.mp4")
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
      {"x"},
      {"i", "m"},
      {"l"},
      {"a", "s"},
      {"e"},
      {"d", "C", "E"},
   })
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

type saved_state struct {
   Dash     *hboMax.Dash
   Login    *hboMax.Login
   Playback *hboMax.Playback
   St       *http.Cookie
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


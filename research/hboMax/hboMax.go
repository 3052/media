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
   if c.proxy_save {
      err := cache.Write(state)
      if err != nil {
         return err
      }
   }
   return maya.SetProxy(state.Proxy, "*.mp4")
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

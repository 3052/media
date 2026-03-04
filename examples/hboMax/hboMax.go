package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "fmt"
)

func (c *client) do_address() error {
   show, err := hboMax.ParseUrl(c.address)
   if err != nil {
      return err
   }
   err = cache.Read(&state)
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

func (c *client) do_edit() error {
   err := cache.Update(&state, func() error {
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

func (c *client) do_dash() error {
   err := cache.Read(&state)
   if err != nil {
      return err
   }
   job.Send = state.Playback.PlayReady
   return job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}


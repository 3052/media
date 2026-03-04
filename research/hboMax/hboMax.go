package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "flag"
   "fmt"
   "log"
   "net/http"
)

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

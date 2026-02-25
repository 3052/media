package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "fmt"
   "net/http"
   "path"
)

func (p *program) run_proxy() error {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      if path.Ext(req.URL.Path) == ".mp4" {
         return "", false
      }
      return p.proxy, true
   })
   return nil
}

func (p *program) run_address() error {
   var show hboMax.ShowKey
   err := show.Parse(p.address)
   if err != nil {
      return err
   }
   cache, err := maya.Read[cache_data](p.cache_file)
   if err != nil {
      return err
   }
   var videos *hboMax.Videos
   if p.season >= 1 {
      videos, err = cache.Login.Season(&show, p.season)
   } else {
      videos, err = cache.Login.Movie(&show)
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

func (p *program) run_login() error {
   cache, err := maya.Read[cache_data](p.cache_file)
   if err != nil {
      return err
   }
   cache.Login, err = cache.St.Login()
   if err != nil {
      return err
   }
   return maya.Write(p.cache_file, cache)
}

type cache_data struct {
   Dash     *hboMax.Dash
   Login    *hboMax.Login
   Playback *hboMax.Playback
   St       *hboMax.St
}

func (p *program) run_dash() error {
   cache, err := maya.Read[cache_data](p.cache_file)
   if err != nil {
      return err
   }
   p.job.Send = cache.Playback.PlayReady
   return p.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, p.dash)
}

func (p *program) run_edit() error {
   cache, err := maya.Read[cache_data](p.cache_file)
   if err != nil {
      return err
   }
   cache.Playback, err = cache.Login.PlayReady(p.edit)
   if err != nil {
      return err
   }
   cache.Dash, err = cache.Playback.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(p.cache_file, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

func (p *program) run_initiate() error {
   var st hboMax.St
   err := st.Fetch()
   if err != nil {
      return err
   }
   err = maya.Write(p.cache_file, &cache_data{St: &st})
   if err != nil {
      return err
   }
   initiate, err := st.Initiate(p.market)
   if err != nil {
      return err
   }
   fmt.Println(initiate)
   return nil
}

package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "flag"
   "fmt"
   "log"
   "net/http"
)

func (c *client) do() error {
   err := cache.Setup("rosso/hboMax.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   // 1
   playReady := maya.StringVar(&c.Job.PlayReady, "p", "PlayReady")
   // 2
   initiate := maya.BoolVar(new(bool), "i", "initiate")
   c.market = hboMax.Markets[0]
   market := maya.StringVar(&c.market, "m", fmt.Sprint(hboMax.Markets))
   // 3
   login := maya.BoolVar(new(bool), "l", "login")
   // 4
   address := maya.StringVar(&c.address, "a", "address")
   season := maya.IntVar(&c.season, "s", "season")
   // 5
   edit := maya.StringVar(&c.edit, "e", "edit ID")
   // 6
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case set[playReady]:
      return cache.Write(c)
   case set[initiate]:
      return c.do_initiate()
   case set[login]:
      return with_cache(c.do_login)
   case set[address]:
      return with_cache(c.do_address)
   case set[edit]:
      return with_cache(c.do_edit)
   case set[dash_id]:
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {playReady},
      {initiate, market},
      {login},
      {address, season},
      {edit},
      {dash_id},
   })
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playback.PlayReady,
   )
}

func (c *client) do_initiate() error {
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
   return cache.Write(c)
}

func (c *client) do_edit() error {
   var err error
   c.Playback, err = c.Login.PlayReady(c.edit)
   if err != nil {
      return err
   }
   c.Dash, err = c.Playback.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Dash     *hboMax.Dash
   Login    *hboMax.Login
   Playback *hboMax.Playback
   St       *http.Cookie
   // 1
   Job maya.Job
   // 2
   market string
   // 4
   address string
   season  int
   // 5
   edit string
   // 6
   dash_id string
}

func (c *client) do_login() error {
   var err error
   c.Login, err = hboMax.FetchLogin(c.St)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
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

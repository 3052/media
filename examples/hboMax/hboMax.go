package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "flag"
   "fmt"
   "log"
   "net/http"
)

func (c *client) do_dash_id() error {
   if cache.Error != nil {
      return cache.Error
   }
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playback.PlayReady,
   )
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

func (c *client) do() error {
   err := cache.Setup("rosso/hboMax.xml")
   if err != nil {
      return err
   }
   cache.Read(c)
   // 1
   flag.BoolVar(&c.initiate, "i", false, "device initiate")
   flag.StringVar(
      &c.market, "m", hboMax.Markets[0], fmt.Sprint(hboMax.Markets),
   )
   // 2
   flag.BoolVar(&c.login, "l", false, "device login")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   flag.IntVar(&c.season, "s", 0, "season")
   // 4
   flag.StringVar(&c.edit, "e", "", "edit ID")
   // 5
   flag.StringVar(&c.Job.PlayReady, "p", c.Job.PlayReady, "PlayReady")
   // 6
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   set := maya.Parse()
   if set["i"] {
      return c.do_initiate()
   }
   if set["l"] {
      return c.do_login()
   }
   if set["a"] {
      return c.do_address()
   }
   if set["e"] {
      return c.do_edit()
   }
   if set["p"] {
      return cache.Write(c)
   }
   if set["d"] {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"i", "m"},
      {"l"},
      {"a", "s"},
      {"e"},
      {"p"},
      {"d"},
   })
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

func (c *client) do_login() error {
   if cache.Error != nil {
      return cache.Error
   }
   c.Login = &hboMax.Login{}
   err := c.Login.Fetch(c.St)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
   if cache.Error != nil {
      return cache.Error
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

type client struct {
   Dash     *hboMax.Dash
   Login    *hboMax.Login
   Playback *hboMax.Playback
   St       *http.Cookie
   // 1
   initiate bool
   market   string
   // 2
   login bool
   // 3
   address string
   season  int
   // 4
   edit string
   // 5
   Job maya.Job
   // 6
   dash_id string
}

func (c *client) do_edit() error {
   if cache.Error != nil {
      return cache.Error
   }
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

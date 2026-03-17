package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "flag"
   "fmt"
   "log"
   "net/http"
)

func (c *client) do_dash_id(err error) error {
   if err != nil {
      return err
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
   err = cache.Read(c)
   // 1
   flag.StringVar(&c.Job.PlayReady, "p", c.Job.PlayReady, "PlayReady")
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
   flag.StringVar(&c.edit_id, "e", "", "edit ID")
   // 6
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   set := maya.Parse()
   switch {
   case set["p"]:
      return cache.Write(c)
   case set["i"]:
      return c.do_initiate()
   case set["l"]:
      return c.do_login(err)
   case set["a"]:
      return c.do_address(err)
   case set["e"]:
      return c.do_edit_id(err)
   case set["d"]:
      return c.do_dash_id(err)
   }
   return maya.Usage([][]string{
      {"p"},
      {"i", "m"},
      {"l"},
      {"a", "s"},
      {"e"},
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

func (c *client) do_login(err error) error {
   if err != nil {
      return err
   }
   c.Login, err = hboMax.FetchLogin(c.St)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address(err error) error {
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

type client struct {
   Dash     *hboMax.Dash
   Login    *hboMax.Login
   Playback *hboMax.Playback
   St       *http.Cookie
   // 1
   Job maya.Job
   // 2
   initiate bool
   market   string
   // 3
   login bool
   // 4
   address string
   season  int
   // 5
   edit_id string
   // 6
   dash_id string
}

func (c *client) do_edit_id(err error) error {
   if err != nil {
      return err
   }
   c.Playback, err = c.Login.PlayReady(c.edit_id)
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

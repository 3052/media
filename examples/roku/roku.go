package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "flag"
   "fmt"
   "log"
)

func (c *client) do_roku(err error) error {
   var code *roku.Code
   if c.get_code {
      if err != nil {
         return err
      }
      code = c.Code
   }
   c.Token, err = roku.FetchToken(code)
   if err != nil {
      return err
   }
   c.Playback, err = c.Token.Playback(c.roku)
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

func (c *client) do_set_code(err error) error {
   if err != nil {
      return err
   }
   c.Code, err = c.Token.Code(c.Activation)
   if err != nil {
      return err
   }
   return cache.Write(c)
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
   err := cache.Setup("rosso/roku.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c)
   // 1
   flag.StringVar(&c.Job.Widevine, "w", c.Job.Widevine, "Widevine")
   // 2
   flag.BoolVar(&c.token, "c", false, "token")
   // 3
   flag.BoolVar(&c.set_code, "s", false, "set code")
   // 4
   flag.StringVar(&c.roku, "r", "", "Roku ID")
   flag.BoolVar(&c.get_code, "g", false, "get code")
   // 5
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   set := maya.Parse()
   switch {
   case set["w"]:
      return cache.Write(c)
   case set["c"]:
      return c.do_token()
   case set["s"]:
      return c.do_set_code(err)
   case set["r"]:
      return c.do_roku(err)
   case set["d"]:
      return c.do_dash_id(err)
   }
   return maya.Usage([][]string{
      {"w"},
      {"c"},
      {"s"},
      {"r", "g"},
      {"d"},
   })
}

func (c *client) do_token() error {
   var err error
   c.Token, err = roku.FetchToken(nil)
   if err != nil {
      return err
   }
   c.Activation, err = c.Token.Activation()
   if err != nil {
      return err
   }
   fmt.Println(c.Activation)
   return cache.Write(c)
}

type client struct {
   Activation   *roku.Activation
   Code       *roku.Code
   Dash       *roku.Dash
   Playback   *roku.Playback
   Token *roku.Token
   // 1
   Job maya.Job
   // 2
   token bool
   // 3
   set_code bool
   // 4
   roku     string
   get_code bool
   // 5
   dash_id string
}

func (c *client) do_dash_id(err error) error {
   if err != nil {
      return err
   }
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playback.Widevine,
   )
}

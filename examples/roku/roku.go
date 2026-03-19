package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "flag"
   "fmt"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/roku.xml")
   if err != nil {
      return err
   }
   read_err := cache.Read(c)
   // 1
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   // 2
   token := maya.BoolVar(new(bool), "t", "token")
   // 3
   set_code := maya.BoolVar(new(bool), "s", "set code")
   // 4
   roku_id := maya.StringVar(&c.roku_id, "r", "Roku ID")
   get_code := maya.BoolVar(&c.get_code, "g", "get code")
   // 5
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case len(set) == 0:
      return maya.Usage([][]*flag.Flag{
         {widevine},
         {token},
         {set_code},
         {roku_id, get_code},
         {dash_id},
      })
   case set[widevine]:
      return cache.Write(c)
   case set[token]:
      return c.do_token()
   case set[roku_id] && !set[get_code]:
      return c.do_roku_id()
   case read_err != nil:
      return read_err
   case set[set_code]:
      return c.do_set_code()
   case set[roku_id]:
      return c.do_roku_id()
   case set[dash_id]:
      return c.do_dash_id()
   }
   return nil
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playback.Widevine,
   )
}

func (c *client) do_set_code() error {
   var err error
   c.Code, err = c.Token.Code(c.Activation)
   if err != nil {
      return err
   }
   return cache.Write(c)
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

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_roku_id() error {
   var code *roku.Code
   if c.get_code {
      code = c.Code
   }
   var err error
   c.Token, err = roku.FetchToken(code)
   if err != nil {
      return err
   }
   c.Playback, err = c.Token.Playback(c.roku_id)
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

type client struct {
   Activation *roku.Activation
   Code       *roku.Code
   Dash       *roku.Dash
   Playback   *roku.Playback
   Token      *roku.Token
   // 1
   Job maya.Job
   // 4
   roku_id string
   get_code bool
   // 5
   dash_id string
}

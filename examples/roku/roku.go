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
   with_cache := cache.Read(c)
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   token := maya.BoolVar(new(bool), "t", "token")
   //----------------------------------------------------------
   set_code := maya.BoolVar(new(bool), "s", "set code")
   //----------------------------------------------------------
   roku_id := maya.StringVar(&c.roku_id, "r", "Roku ID")
   get_code := maya.BoolVar(&c.get_code, "g", "get code")
   //----------------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   if set[widevine] {
      return cache.Write(c)
   }
   if set[token] {
      return c.do_token()
   }
   if set[set_code] {
      return with_cache(c.do_set_code)
   }
   if set[roku_id] {
      if set[get_code] {
         return with_cache(c.do_roku_id)
      }
      return c.do_roku_id()
   }
   if set[dash_id] {
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {token},
      {set_code},
      {roku_id, get_code},
      {dash_id},
   })
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
   //--------------------
   Job maya.Job
   //--------------------
   roku_id  string
   get_code bool
   //--------------------
   dash_id string
}

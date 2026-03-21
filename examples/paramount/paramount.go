package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "flag"
   "log"
   "net/http"
)

func (c *client) do_dash_id() error {
   at, err := paramount.GetAt(c.AppSecret)
   if err != nil {
      return err
   }
   if !c.get_cookie {
      c.Cookie = nil
   }
   token, err := paramount.PlayReady(at, c.ParamountId, c.Cookie)
   if err != nil {
      return err
   }
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, token.Send)
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s,*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   AppSecret string
   Cookie *http.Cookie
   Dash   *paramount.Dash
   //--------------------
   Job maya.Job
   //--------------------
   username string
   password string
   //--------------------
   ParamountId string
   //--------------------
   dash_id    string
   get_cookie bool
}

func (c *client) do_app_secret() error {
   var err error
   c.AppSecret, err = paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_username_password() error {
   at, err := paramount.GetAt(c.AppSecret)
   if err != nil {
      return err
   }
   c.Cookie, err = paramount.Login(at, c.username, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_paramount() error {
   at, err := paramount.GetAt(c.AppSecret)
   if err != nil {
      return err
   }
   if !c.get_cookie {
      c.Cookie = nil
   }
   video, err := paramount.FetchVideo(at, c.ParamountId, c.Cookie)
   if err != nil {
      return err
   }
   c.Dash, err = video.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do() error {
   err := cache.Setup("rosso/paramount.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   playReady := maya.StringVar(&c.Job.PlayReady, "PR", "PlayReady")
   //--------------------------------------------------------------
   app_secret := maya.BoolVar(new(bool), "a", "app secret")
   //--------------------------------------------------------------
   username := maya.StringVar(&c.username, "U", "username")
   password := maya.StringVar(&c.password, "P", "password")
   //--------------------------------------------------------------
   paramount_id := maya.StringVar(&c.ParamountId, "p", "paramount ID")
   get_cookie := maya.BoolVar(&c.get_cookie, "c", "get cookie")
   //--------------------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   if set[playReady] {
      return cache.Write(c)
   }
   if set[app_secret] {
      return c.do_app_secret()
   }
   if set[username] {
      if set[password] {
         return with_cache(c.do_username_password)
      }
   }
   if set[paramount_id] {
      return with_cache(c.do_paramount)
   }
   if set[dash_id] {
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {playReady},
      {app_secret},
      {username, password},
      
      {paramount_id, get_cookie},
      {dash_id, get_cookie},
   })
}

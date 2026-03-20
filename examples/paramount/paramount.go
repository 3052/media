package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "flag"
   "log"
   "net/http"
)

func (c *client) do_paramount() error {
   var err error
   c.Dash, err = paramount.FetchDash(c.paramount_id)
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
   username := maya.StringVar(&c.username, "U", "username")
   password := maya.StringVar(&c.password, "P", "password")
   //--------------------------------------------------------------
   paramount_id := maya.StringVar(&c.paramount_id, "p", "paramount ID")
   //--------------------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   get_cookie := maya.BoolVar(&c.get_cookie, "c", "get cookie")
   set := maya.Parse()
   if set[playReady] {
      return cache.Write(c)
   }
   if set[username] {
      if set[password] {
         return c.do_username_password()
      }
   }
   if set[paramount_id] {
      if set[get_cookie] {
         return with_cache(c.do_paramount)
      }
      return c.do_paramount()
   }
   if set[dash_id] {
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {playReady},
      {username, password},
      {paramount_id},
      {dash_id, get_cookie},
   })
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s,*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Cookie *http.Cookie
   Dash   *paramount.Dash
   //--------------------
   Job maya.Job
   //--------------------
   username string
   password string
   //--------------------
   paramount_id string
   //--------------------
   dash_id    string
   get_cookie bool
}

func (c *client) do_dash_id() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   if !c.get_cookie {
      c.Cookie = nil
   }
   token, err := paramount.PlayReady(at, c.paramount_id, c.Cookie)
   if err != nil {
      return err
   }
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, token.Send)
}

var cache maya.Cache

func (c *client) do_username_password() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   c.Cookie, err = paramount.Login(at, c.username, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

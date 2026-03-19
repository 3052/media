package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/kanopy"
   "flag"
   "fmt"
   "log"
)

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id,
      func(data []byte) ([]byte, error) {
         return c.Login.Widevine(c.PlayManifest, data)
      },
   )
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

type client struct {
   Dash         *kanopy.Dash
   Login        *kanopy.Login
   PlayManifest *kanopy.PlayManifest
   // 1
   Job maya.Job
   // 2
   email    string
   password string
   // 3
   address string
   // 4
   dash_id string
}

func (c *client) do() error {
   err := cache.Setup("rosso/kanopy.xml")
   if err != nil {
      return err
   }
   read_err := cache.Read(c)
   // 1
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   // 2
   email := maya.StringVar(&c.email, "e", "email")
   password := maya.StringVar(&c.password, "p", "password")
   // 3
   address := maya.StringVar(&c.address, "a", "address")
   // 4
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case len(set) == 0:
      return maya.Usage([][]*flag.Flag{
         {widevine},
         {email, password},
         {address},
         {dash_id},
      })
   case set[email.Name] && set[password.Name]:
      return c.do_email_password()
   case read_err != nil:
      return read_err
   case set[address.Name]:
      return c.do_address()
   case set[dash_id.Name]:
      return c.do_dash_id()
   }
   return nil
}

func (c *client) do_email_password() error {
   var err error
   c.Login, err = kanopy.FetchLogin(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
   video, err := kanopy.ParseVideo(c.address)
   if err != nil {
      return err
   }
   if video.VideoId == 0 {
      video, err = c.Login.Video(video.Alias)
      if err != nil {
         return err
      }
   }
   member, err := c.Login.Membership()
   if err != nil {
      return err
   }
   plays, err := c.Login.Plays(member, video.VideoId)
   if err != nil {
      return err
   }
   for _, caption := range plays.Captions {
      for _, file := range caption.Files {
         fmt.Println(file.Url)
      }
   }
   c.PlayManifest, err = plays.Dash()
   if err != nil {
      return err
   }
   c.Dash, err = c.PlayManifest.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

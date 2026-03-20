package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/plex"
   "flag"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/plex.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   // 1
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   // 2
   address := maya.StringVar(&c.address, "a", "address")
   xff := maya.StringVar(&c.xff, "x", "x-forwarded-for")
   // 3
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case set[widevine]:
      return cache.Write(c)
   case set[address]:
      return c.do_address()
   case set[dash_id]:
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {address, xff},
      {dash_id},
   })
}

func (c *client) do_address() error {
   var err error
   c.User, err = plex.FetchUser()
   if err != nil {
      return err
   }
   address, err := plex.GetPath(c.address)
   if err != nil {
      return err
   }
   metadata, err := c.User.RatingKey(address)
   if err != nil {
      return err
   }
   metadata, err = c.User.Media(metadata, c.xff)
   if err != nil {
      return err
   }
   c.MediaPart, err = metadata.Dash()
   if err != nil {
      return err
   }
   c.Dash, err = c.User.Dash(c.MediaPart, c.xff)
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
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

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id,
      func(data []byte) ([]byte, error) {
         return c.User.Widevine(c.MediaPart, data)
      },
   )
}

type client struct {
   Dash      *plex.Dash
   MediaPart *plex.MediaPart
   User      *plex.User
   // 1
   Job maya.Job
   // 2
   address         string
   xff string
   // 3
   dash_id string
}

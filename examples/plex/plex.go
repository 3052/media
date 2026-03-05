package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/plex"
   "flag"
   "log"
)

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = func(data []byte) ([]byte, error) {
      return c.User.Widevine(c.MediaPart, data)
   }
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
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

var job maya.WidevineJob

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := cache.Setup("rosso/plex.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.x_forwarded_for, "x", "", "x-forwarded-for")
   // 2
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "c", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "p", job.PrivateKey, "private key")
   flag.Parse()
   if c.address != "" {
      return c.do_address()
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"a", "x"},
      {"d", "c", "p"},
   })
}

type client struct {
   Dash      *plex.Dash
   MediaPart *plex.MediaPart
   User      *plex.User
   // 1
   address         string
   x_forwarded_for string
   // 2
   dash_id string
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
   metadata, err = c.User.Media(metadata, c.x_forwarded_for)
   if err != nil {
      return err
   }
   c.MediaPart, err = metadata.Dash()
   if err != nil {
      return err
   }
   c.Dash, err = c.User.Dash(c.MediaPart, c.x_forwarded_for)
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "flag"
   "fmt"
   "log"
)

func BoolVar(value *bool, name, usage string) *flag.Flag {
   flag.BoolVar(value, name, *value, usage)
   return flag.Lookup(name)
}

func IntVar(value *int, name, usage string) *flag.Flag {
   flag.IntVar(value, name, *value, usage)
   return flag.Lookup(name)
}

func StringVar(value *string, name, usage string) *flag.Flag {
   flag.StringVar(value, name, *value, usage)
   return flag.Lookup(name)
}

func IsSet(f *flag.Flag) bool {
   var set bool
   flag.Visit(func(g *flag.Flag) {
      if g.Name == f.Name {
         set = true
      }
   })
   return set
}

func (c *client) do_email() error {
   var err error
   c.Token, err = disney.RegisterDevice()
   if err != nil {
      return err
   }
   request_otp, err := c.Token.RequestOtp(c.Email)
   if err != nil {
      return err
   }
   fmt.Println(request_otp)
   return cache.Write(c)
}

func (c *client) do_passcode() error {
   otp, err := c.Token.AuthenticateWithOtp(c.Email, c.passcode)
   if err != nil {
      return err
   }
   login, err := c.Token.LoginWithActionGrant(otp.ActionGrant)
   if err != nil {
      return err
   }
   for i, profile := range login.Account.Profiles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&profile)
   }
   return cache.Write(c)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4,*.mp4a")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

func (c *client) do_profile_id() error {
   err := c.Token.SwitchProfile(c.profile_id)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
   entity, err := disney.GetEntity(c.address)
   if err != nil {
      return err
   }
   page, err := c.Token.Page(entity)
   if err != nil {
      return err
   }
   fmt.Println(page)
   return nil
}

func (c *client) do_season_id() error {
   season, err := c.Token.Season(c.season_id)
   if err != nil {
      return err
   }
   fmt.Println(season)
   return nil
}

func (c *client) do_media_id() error {
   stream, err := c.Token.Stream(c.media_id)
   if err != nil {
      return err
   }
   c.Hls, err = stream.Hls()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListHls(c.Hls.Body, c.Hls.Url)
}

func (c *client) do_refresh() error {
   err := disney.RefreshToken(c.Token)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

type client struct {
   Email      string
   Hls        *disney.Hls
   Job        maya.Job
   Token      *disney.Token
   address    string
   hls_id     int
   media_id   string
   passcode   string
   profile_id string
   season_id  string
}

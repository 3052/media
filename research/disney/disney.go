package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "flag"
   "fmt"
   "log"
)

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
func (c *client) do() error {
   err := cache.Setup("rosso/disney.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c)
   // 1
   playReady := StringVar(&c.Job.PlayReady, "PR", "PlayReady")
   // 2
   email := StringVar(&c.Email, "e", "email")
   // 3
   passcode := StringVar(&c.passcode, "p", "passcode")
   // 4
   profile := StringVar(&c.profile_id, "P", "profile ID")
   // 5
   refresh := BoolVar(new(bool), "r", "refresh")
   // 6
   address := StringVar(&c.address, "a", "address")
   // 7
   season := StringVar(&c.season_id, "s", "season ID")
   // 8
   media := StringVar(&c.media_id, "m", "media")
   // 9
   hls := IntVar(&c.hls_id, "h", "HLS ID")
   flag.Parse()
   switch {
   case IsSet(playReady):
      return cache.Write(c)
   case IsSet(email):
      return c.do_email()
   }
   if err != nil {
      return err
   }
   switch {
   case IsSet(passcode):
      return c.do_passcode()
   case IsSet(profile):
      return c.do_profile_id()
   case IsSet(refresh):
      return c.do_refresh()
   case IsSet(address):
      return c.do_address()
   case IsSet(season):
      return c.do_season_id()
   case IsSet(media):
      return c.do_media_id()
   case IsSet(hls):
      return c.Job.DownloadHls(
         c.Hls.Body, c.Hls.Url, c.hls_id, c.Token.PlayReady,
      )
   }
   return Usage([][]*flag.Flag{
      {playReady},
      {email},
      {passcode},
      {profile},
      {refresh},
      {address},
      {season},
      {media},
      {hls},
   })
}

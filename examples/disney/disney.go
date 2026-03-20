package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "flag"
   "fmt"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/disney.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   // 1
   playReady := maya.StringVar(&c.Job.PlayReady, "PR", "PlayReady")
   // 2
   email := maya.StringVar(&c.Email, "e", "email")
   // 3
   passcode := maya.StringVar(&c.passcode, "p", "passcode")
   // 4
   profile := maya.StringVar(&c.profile, "P", "profile ID")
   // 5
   refresh := maya.BoolVar(new(bool), "r", "refresh")
   // 6
   address := maya.StringVar(&c.address, "a", "address")
   // 7
   season := maya.StringVar(&c.season, "s", "season ID")
   // 8
   media := maya.StringVar(&c.media, "m", "media ID")
   // 9
   hls_id := maya.IntVar(&c.hls_id, "h", "HLS ID")
   set := maya.Parse()
   switch {
   case set[playReady]:
      return cache.Write(c)
   case set[email]:
      return c.do_email()
   case set[passcode]:
      return with_cache(c.do_passcode)
   case set[profile]:
      return with_cache(c.do_profile_id)
   case set[refresh]:
      return with_cache(c.do_refresh)
   case set[address]:
      return with_cache(c.do_address)
   case set[season]:
      return with_cache(c.do_season_id)
   case set[media]:
      return with_cache(c.do_media_id)
   case set[hls_id]:
      return with_cache(c.do_hls_id)
   }
   return maya.Usage([][]*flag.Flag{
      {playReady},
      {email},
      {passcode},
      {profile},
      {refresh},
      {address},
      {season},
      {media},
      {hls_id},
   })
}

func (c *client) do_hls_id() error {
   return c.Job.DownloadHls(c.Hls.Body, c.Hls.Url, c.hls_id, c.Token.PlayReady)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4,*.mp4a")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Hls   *disney.Hls
   Token *disney.Token
   // 1
   Job maya.Job
   // 2
   Email string
   // 3
   passcode string
   // 4
   profile string
   // 6
   address string
   // 7
   season string
   // 8
   media string
   // 9
   hls_id int
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

var cache maya.Cache

func (c *client) do_profile_id() error {
   err := c.Token.SwitchProfile(c.profile)
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
   season, err := c.Token.Season(c.season)
   if err != nil {
      return err
   }
   fmt.Println(season)
   return nil
}

func (c *client) do_media_id() error {
   stream, err := c.Token.Stream(c.media)
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

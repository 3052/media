package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "flag"
   "fmt"
   "log"
)

type client struct {
   Hls   *disney.Hls
   Token *disney.Token
   // 1
   Email string
   // 2
   passcode string
   // 3
   profile_id string
   // 4
   refresh bool
   // 5
   address string
   // 6
   season_id string
   // 7
   media_id string
   // 8
   hls_id int
}

func main() {
   log.SetFlags(log.Ltime)
   // MP4 ARE GEO BLOCKED SO JUST USE VPN
   maya.SetProxy("", "*.mp4,*.mp4a")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   job.CertificateChain, _ = maya.ResolveCache("SL3000/CertificateChain")
   job.EncryptSignKey, _ = maya.ResolveCache("SL3000/EncryptSignKey")
   err := cache.Setup("rosso/disney.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.Email, "e", "", "email")
   // 2
   flag.StringVar(&c.passcode, "p", "", "passcode")
   // 3
   flag.StringVar(&c.profile_id, "P", "", "profile ID")
   // 4
   flag.BoolVar(&c.refresh, "r", false, "refresh")
   // 5
   flag.StringVar(&c.address, "a", "", "address")
   // 6
   flag.StringVar(&c.season_id, "s", "", "season ID")
   // 7
   flag.StringVar(&c.media_id, "m", "", "media ID")
   // 8
   flag.IntVar(&c.hls_id, "h", -1, "HLS ID")
   flag.IntVar(&job.Threads, "t", 2, "threads")
   flag.StringVar(&job.CertificateChain, "C", job.CertificateChain, "certificate chain")
   flag.StringVar(&job.EncryptSignKey, "E", job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   if c.Email != "" {
      return c.do_email()
   }
   if c.passcode != "" {
      return c.do_passcode()
   }
   if c.profile_id != "" {
      return c.do_profile_id()
   }
   if c.refresh {
      return c.do_refresh()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.season_id != "" {
      return c.do_season_id()
   }
   if c.media_id != "" {
      return c.do_media_id()
   }
   if c.hls_id >= 0 {
      return c.do_hls_id()
   }
   return maya.Usage([][]string{
      {"e"},
      {"p"},
      {"P"},
      {"r"},
      {"a"},
      {"s"},
      {"m"},
      {"h", "t", "C", "E"},
   })
}

func (c *client) do_profile_id() error {
   return cache.Update(c, func() error {
      return c.Token.SwitchProfile(c.profile_id)
   })
}

func (c *client) do_passcode() error {
   return cache.Update(c, func() error {
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
      return nil
   })
}

var job maya.PlayReadyJob

var cache maya.Cache

func (c *client) do_email() error {
   c.Token = &disney.Token{}
   err := c.Token.RegisterDevice()
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

func (c *client) do_hls_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = c.Token.PlayReady
   return job.DownloadHls(c.Hls.Body, c.Hls.Url, c.hls_id)
}

func (c *client) do_media_id() error {
   err := cache.Update(c, func() error {
      stream, err := c.Token.Stream(c.media_id)
      if err != nil {
         return err
      }
      c.Hls, err = stream.Hls()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListHls(c.Hls.Body, c.Hls.Url)
}

func (c *client) do_season_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   season, err := c.Token.Season(c.season_id)
   if err != nil {
      return err
   }
   fmt.Println(season)
   return nil
}

func (c *client) do_address() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
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

func (c *client) do_refresh() error {
   return cache.Update(c, func() error {
      return disney.RefreshToken(c.Token)
   })
}

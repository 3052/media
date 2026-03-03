package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "flag"
   "fmt"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4,*.mp4a")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_hls() error {
   var state saved_state
   err := c.cache.Read(&state)
   if err != nil {
      return err
   }
   c.job.Send = state.Account.PlayReady
   return c.job.DownloadHls(state.Hls.Body, state.Hls.Url, c.hls)
}

func (c *client) do() error {
   c.job.CertificateChain, _ = maya.ResolveCache("SL3000/CertificateChain")
   c.job.EncryptSignKey, _ = maya.ResolveCache("SL3000/EncryptSignKey")
   err := c.cache.Setup("rosso/disney.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.profile_id, "P", "", "profile ID")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   // 4
   flag.StringVar(&c.season_id, "s", "", "season ID")
   // 5
   flag.StringVar(&c.media_id, "m", "", "media ID")
   // 6
   flag.IntVar(&c.hls, "h", -1, "HLS ID")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.profile_id != "" {
      return c.do_profile_id()
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
   if c.hls >= 0 {
      return c.do_hls()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"P"},
      {"a"},
      {"s"},
      {"m"},
      {"h", "C", "E"},
   })
}

func (c *client) do_email_password() error {
   device, err := disney.RegisterDevice()
   if err != nil {
      return err
   }
   inactive_account, err := device.Login(c.email, c.password)
   if err != nil {
      return err
   }
   for i, profile := range inactive_account.Data.Login.Account.Profiles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&profile)
   }
   return c.cache.Write(saved_state{InactiveAccount: inactive_account})
}

type saved_state struct {
   Account         *disney.Account
   Hls             *disney.Hls
   InactiveAccount *disney.InactiveAccount
}

func (c *client) do_profile_id() error {
   var state saved_state
   return c.cache.Update(&state, func() error {
      var err error
      state.Account, err = state.InactiveAccount.SwitchProfile(c.profile_id)
      return err
   })
}

func (c *client) do_address() error {
   var state saved_state
   err := c.cache.Update(&state, func() error {
      return state.Account.RefreshToken()
   })
   if err != nil {
      return err
   }
   entity, err := disney.GetEntity(c.address)
   if err != nil {
      return err
   }
   page, err := state.Account.Page(entity)
   if err != nil {
      return err
   }
   fmt.Println(page)
   return nil
}

func (c *client) do_season_id() error {
   var state saved_state
   err := c.cache.Read(&state)
   if err != nil {
      return err
   }
   season, err := state.Account.Season(c.season_id)
   if err != nil {
      return err
   }
   fmt.Println(season)
   return nil
}

type client struct {
   cache maya.Cache
   // 1
   email    string
   password string
   // 2
   profile_id string
   // 3
   address string
   // 4
   season_id string
   // 5
   media_id string
   // 6
   hls int
   job maya.PlayReadyJob
}

func (c *client) do_media_id() error {
   var state saved_state
   err := c.cache.Update(&state, func() error {
      stream, err := state.Account.Stream(c.media_id)
      if err != nil {
         return err
      }
      state.Hls, err = stream.Hls()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListHls(state.Hls.Body, state.Hls.Url)
}

package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "flag"
   "fmt"
   "log"
)

func (c *client) do_passcode() error {
   device, err := disney.RegisterDevice()
   if err != nil {
      return err
   }
   c.InactiveAccount, err = device.Login(c.email, c.password)
   if err != nil {
      return err
   }
   for i, profile := range c.InactiveAccount.Data.Login.Account.Profiles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&profile)
   }
   return cache.Write(c)
}

func (c *client) do_hls_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = c.Account.PlayReady
   return job.DownloadHls(c.Hls.Body, c.Hls.Url, c.hls_id)
}

func (c *client) do_media_id() error {
   err := cache.Update(c, func() error {
      stream, err := c.Account.Stream(c.media_id)
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
   season, err := c.Account.Season(c.season_id)
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
   page, err := c.Account.Page(entity)
   if err != nil {
      return err
   }
   fmt.Println(page)
   return nil
}

func (c *client) do_refresh() error {
   return cache.Update(c, func() error {
      return c.Account.RefreshToken()
   })
}

func (c *client) do_profile_id() error {
   return cache.Update(c, func() error {
      var err error
      c.Account, err = c.InactiveAccount.SwitchProfile(c.profile_id)
      return err
   })
}

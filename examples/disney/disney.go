package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "flag"
   "fmt"
   "os"
   "path/filepath"
)

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.job.CertificateChain = cache + "/SL3000/CertificateChain"
   c.job.EncryptSignKey = cache + "/SL3000/EncryptSignKey"
   c.name = cache + "/rosso/disney.xml"
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

type user_cache struct {
   Account *disney.Account
   AccountWithoutActiveProfile *disney.AccountWithoutActiveProfile
   Hls     *disney.Hls
}

func (c *command) do_email_password() error {
   device, err := disney.RegisterDevice()
   if err != nil {
      return err
   }
   var cache user_cache
   cache.AccountWithoutActiveProfile, err = device.Login(c.email, c.password)
   if err != nil {
      return err
   }
   for i, profile := range cache.AccountWithoutActiveProfile.Data.Login.Account.Profiles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&profile)
   }
   return maya.Write(c.name, cache)
}

func (c *command) do_profile_id() error {
   var cache user_cache
   err := maya.Read(c.name, &cache)
   if err != nil {
      return err
   }
   cache.Account, err = cache.AccountWithoutActiveProfile.SwitchProfile(
      c.profile_id,
   )
   if err != nil {
      return err
   }
   return maya.Write(c.name, cache)
}

func (c *command) do_hls() error {
   var cache user_cache
   err := maya.Read(c.name, &cache)
   if err != nil {
      return err
   }
   c.job.Send = cache.Account.PlayReady
   return c.job.DownloadHls(cache.Hls.Body, cache.Hls.Url, c.hls)
}

func (c *command) do_address() error {
   var cache user_cache
   err := maya.Read(c.name, &cache)
   if err != nil {
      return err
   }
   entity, err := disney.GetEntity(c.address)
   if err != nil {
      return err
   }
   page, err := cache.Account.Page(entity)
   if err != nil {
      return err
   }
   fmt.Println(page)
   return nil
}

func (c *command) do_season_id() error {
   var cache user_cache
   err := maya.Read(c.name, &cache)
   if err != nil {
      return err
   }
   season, err := cache.Account.Season(c.season_id)
   if err != nil {
      return err
   }
   fmt.Println(season)
   return nil
}

func (c *command) do_media_id() error {
   var cache user_cache
   err := maya.Read(c.name, &cache)
   if err != nil {
      return err
   }
   stream, err := cache.Account.Stream(c.media_id)
   if err != nil {
      return err
   }
   cache.Hls, err = stream.Hls()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListHls(cache.Hls.Body, cache.Hls.Url)
}

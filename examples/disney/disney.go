package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "flag"
   "fmt"
   "log"
)

func (c *client) do() error {
   job.CertificateChain, _ = maya.ResolveCache("SL3000/CertificateChain")
   job.EncryptSignKey, _ = maya.ResolveCache("SL3000/EncryptSignKey")
   err := cache.Setup("rosso/disney.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 2
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
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
   flag.IntVar(&c.hls, "h", -1, "HLS ID")
   flag.StringVar(&job.CertificateChain, "C", job.CertificateChain, "certificate chain")
   flag.StringVar(&job.EncryptSignKey, "E", job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   err = maya.SetProxy(c.proxy, "*.mp4,*.mp4a")
   if err != nil {
      return err
   }
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
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
   if c.hls >= 0 {
      return c.do_hls()
   }
   return maya.Usage([][]string{
      {"x"},
      {"e", "p"},
      {"P"},
      {"r"},
      {"a"},
      {"s"},
      {"m"},
      {"h", "C", "E"},
   })
}

func (c *client) do_hls() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = c.Account.PlayReady
   return job.DownloadHls(c.Hls.Body, c.Hls.Url, c.hls)
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

func (c *client) do_email_password() error {
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

var job maya.PlayReadyJob

var cache maya.Cache
func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Account         *disney.Account
   Hls             *disney.Hls
   InactiveAccount *disney.InactiveAccount
   // 1
   proxy string
   // 2
   email    string
   password string
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
   hls int
}

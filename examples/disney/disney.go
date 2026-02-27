package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (c *command) run() error {
   err := c.cache.Init("SL3000")
   if err != nil {
      return err
   }
   c.job.CertificateChain = c.cache.Join("CertificateChain")
   c.job.EncryptSignKey = c.cache.Join("EncryptSignKey")
   err = c.cache.Init("disney")
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
   flag.IntVar(&c.job.Threads, "t", 2, "threads")
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

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      switch path.Ext(req.URL.Path) {
      case ".mp4", ".mp4a":
         return "", false
      }
      return "", true
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
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

func (c *command) do_email_password() error {
   device, err := disney.RegisterDevice()
   if err != nil {
      return err
   }
   account_without, err := device.Login(c.email, c.password)
   if err != nil {
      return err
   }
   for i, profile := range account_without.Data.Login.Account.Profiles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&profile)
   }
   return c.cache.Set("AccountWithoutActiveProfile", account_without)
}

func (c *command) do_profile_id() error {
   var account_without disney.AccountWithoutActiveProfile
   err := c.cache.Get("AccountWithoutActiveProfile", &account_without)
   if err != nil {
      return err
   }
   account, err := account_without.SwitchProfile(c.profile_id)
   if err != nil {
      return err
   }
   return c.cache.Set("Account", account)
}

func (c *command) do_address() error {
   var account disney.Account
   err := c.cache.Get("Account", &account)
   if err != nil {
      return err
   }
   err = account.RefreshToken()
   if err != nil {
      return err
   }
   err = c.cache.Set("Account", account)
   if err != nil {
      return err
   }
   entity, err := disney.GetEntity(c.address)
   if err != nil {
      return err
   }
   page, err := account.Page(entity)
   if err != nil {
      return err
   }
   fmt.Println(page)
   return nil
}

func (c *command) do_season_id() error {
   var account disney.Account
   err := c.cache.Get("Account", &account)
   if err != nil {
      return err
   }
   season, err := account.Season(c.season_id)
   if err != nil {
      return err
   }
   fmt.Println(season)
   return nil
}
func (c *command) do_hls() error {
   var account disney.Account
   err := c.cache.Get("Account", &account)
   if err != nil {
      return err
   }
   var hls disney.Hls
   err = c.cache.Get("Hls", &hls)
   if err != nil {
      return err
   }
   c.job.Send = account.PlayReady
   return c.job.DownloadHls(hls.Body, hls.Url, c.hls)
}

func (c *command) do_media_id() error {
   var account disney.Account
   err := c.cache.Get("Account", &account)
   if err != nil {
      return err
   }
   stream, err := account.Stream(c.media_id)
   if err != nil {
      return err
   }
   hls, err := stream.Hls()
   if err != nil {
      return err
   }
   err = c.cache.Set("Hls", hls)
   if err != nil {
      return err
   }
   return maya.ListHls(hls.Body, hls.Url)
}

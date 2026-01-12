package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/disney"
   "encoding/xml"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      switch path.Ext(req.URL.Path) {
      case ".mp4", ".mp4a":
         return ""
      }
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type user_cache struct {
   Account *disney.Account
   Hls     *disney.Hls
}

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func read(name string) (*user_cache, error) {
   data, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   cache := &user_cache{}
   err = xml.Unmarshal(data, cache)
   if err != nil {
      return nil, err
   }
   return cache, nil
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.name = cache + "/disney/userCache.xml"
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.season, "s", "", "season")
   // 4
   flag.StringVar(&c.media_id, "m", "", "media ID")
   // 5
   flag.StringVar(&c.hls, "h", "", "HLS ID")
   flag.Parse()
   // 1
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   // 2
   if c.address != "" {
      return c.do_address()
   }
   // 3
   if c.season != "" {
      return c.do_season()
   }
   // 4
   if c.media_id != "" {
      return c.do_media_id()
   }
   // 5
   if c.hls != "" {
      return c.do_hls()
   }
   flag.Usage()
   return nil
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
   var cache user_cache
   cache.Account, err = account_without.SwitchProfile()
   if err != nil {
      return err
   }
   return write(c.name, &cache)
}

func (c *command) do_address() error {
   cache, err := read(c.name)
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

func (c *command) do_season() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   season, err := cache.Account.Season(c.season)
   if err != nil {
      return err
   }
   fmt.Println(season)
   return nil
}

type command struct {
   job      maya.WidevineJob
   name     string
   // 1
   email    string
   password string
   // 2
   address  string
   // 3
   season string
   // 4
   media_id string
   // 5
   hls      string
}

func (c *command) do_hls() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Account.Widevine(data)
   }
   return c.job.DownloadHls(cache.Hls.Body, cache.Hls.Url, c.hls)
}

func (c *command) do_media_id() error {
   cache, err := read(c.name)
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
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListHls(cache.Hls.Body, cache.Hls.Url)
}

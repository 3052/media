package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/disney"
   "encoding/xml"
   "errors"
   "flag"
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
      return "L"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.name = cache + "/disney/userCache.xml"
   c.job.ClientID = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"

   flag.StringVar(&c.job.ClientID, "C", c.job.ClientID, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.hls, "h", "", "HLS ID")
   flag.StringVar(&c.password, "p", "", "password")
   flag.Parse()

   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.hls != "" {
      return c.do_hls()
   }
   flag.Usage()
   return nil
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

func (c *command) do_email_password() error {
   var device disney.Device
   err := device.Register()
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

type user_cache struct {
   Account *disney.Account
   Hls     *disney.Hls
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
   explore, err := cache.Account.Explore(entity)
   if err != nil {
      return err
   }
   playback_id, ok := explore.PlaybackId()
   if !ok {
      return errors.New(".PlaybackId()")
   }
   playback, err := cache.Account.Playback(playback_id)
   if err != nil {
      return err
   }
   cache.Hls, err = playback.Hls()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListHls(cache.Hls.Body, cache.Hls.Url)
}

type command struct {
   address  string
   hls      string
   email    string
   job      maya.WidevineJob
   name     string
   password string
}

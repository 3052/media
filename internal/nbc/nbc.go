package main

import (
   "41.neocities.org/media/nbc"
   "41.neocities.org/net"
   "encoding/xml"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/nbc/mpd.xml"

   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.config.ClientId, "c", c.config.ClientId, "client ID")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.config.PrivateKey, "p", c.config.PrivateKey, "private key")
   flag.IntVar(&c.config.Threads, "t", 2, "threads")
   flag.Parse()

   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
}

type mpd struct {
   Body []byte
   Url  *url.URL
}

func (c *command) do_address() error {
   name, err := nbc.GetName(c.address)
   if err != nil {
      return err
   }
   metadata, err := nbc.FetchMetadata(name)
   if err != nil {
      return err
   }
   stream, err := metadata.Stream()
   if err != nil {
      return err
   }
   var cache mpd
   cache.Url, cache.Body, err = stream.Mpd()
   if err != nil {
      return err
   }
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", c.name)
   err = os.WriteFile(c.name, data, os.ModePerm)
   if err != nil {
      return err
   }
   return net.Representations(cache.Url, cache.Body)
}

type command struct {
   config  net.Config
   name    string
   // 1
   address string
   // 2
   dash    string
}
func (c *command) do_dash() error {
   data, err := os.ReadFile(c.name)
   if err != nil {
      return err
   }
   var cache mpd
   err = xml.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   c.config.Send = nbc.Widevine
   return c.config.Download(cache.Url, cache.Body, c.dash)
}

func main() {
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".mp4" {
         return ""
      }
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

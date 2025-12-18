package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/hboMax"
   "encoding/xml"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

type user_cache struct {
   Login    *hboMax.Login
   Mpd      *url.URL
   MpdBody  []byte
   Playback *hboMax.Playback
   St       *hboMax.St
}

func (c *command) do_initiate() error {
   var st hboMax.St
   err := st.Fetch()
   if err != nil {
      return err
   }
   err = write(c.name, &user_cache{St: &st})
   if err != nil {
      return err
   }
   initiate, err := st.Initiate(c.market)
   if err != nil {
      return err
   }
   fmt.Println(initiate)
   return nil
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

func (c *command) do_login() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   cache.Login, err = cache.St.Login()
   if err != nil {
      return err
   }
   return write(c.name, cache)
}

func (c *command) do_address() error {
   show_id, err := hboMax.ExtractId(c.address)
   if err != nil {
      return err
   }
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   var videos *hboMax.Videos
   if c.season >= 1 {
      videos, err = cache.Login.Season(show_id, c.season)
   } else {
      videos, err = cache.Login.Movie(show_id)
   }
   if err != nil {
      return err
   }
   videos.FilterAndSort()
   for i, video := range videos.Included {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(video)
   }
   return nil
}

func (c *command) do_edit() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   cache.Playback, err = cache.Login.PlayReady(c.edit)
   if err != nil {
      return err
   }
   cache.Mpd, cache.MpdBody, err = cache.Playback.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.Representations(cache.Mpd, cache.MpdBody)
}

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      return cache.Playback.PlayReady(data)
   }
   return c.config.Download(cache.Mpd, cache.MpdBody, c.dash)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
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

type command struct {
   address  string
   config   maya.Config
   dash     string
   edit     string
   initiate bool
   login    bool
   name     string
   season   int
   market   string
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.config.CertificateChain = cache + "/SL3000/CertificateChain"
   c.config.EncryptSignKey = cache + "/SL3000/EncryptSignKey"
   c.name = cache + "/hboMax/userCache.xml"

   flag.StringVar(&c.config.CertificateChain, "C", c.config.CertificateChain, "certificate chain")
   flag.StringVar(&c.config.EncryptSignKey, "E", c.config.EncryptSignKey, "encrypt sign key")
   flag.StringVar(&c.address, "a", "", "address")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.edit, "e", "", "edit ID")
   flag.BoolVar(&c.initiate, "i", false, "device initiate")
   flag.BoolVar(&c.login, "l", false, "device login")
   flag.StringVar(
      &c.market, "m", hboMax.Markets[0], fmt.Sprint(hboMax.Markets[1:]),
   )
   flag.IntVar(&c.season, "s", 0, "season")
   flag.Parse()
   if c.initiate {
      return c.do_initiate()
   }
   if c.login {
      return c.do_login()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.edit != "" {
      return c.do_edit()
   }
   if c.dash != "" {
      return c.do_dash()
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

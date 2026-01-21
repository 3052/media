package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/hboMax"
   "encoding/xml"
   "flag"
   "fmt"
   "log"
   "net/http"
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
   c.name = cache + "/hboMax/userCache.xml"
   c.job.CertificateChain = cache + "/SL3000/CertificateChain"
   c.job.EncryptSignKey = cache + "/SL3000/EncryptSignKey"
   // 1
   flag.BoolVar(&c.initiate, "i", false, "device initiate")
   flag.StringVar(
      &c.market, "m", hboMax.Markets[0], fmt.Sprint(hboMax.Markets[1:]),
   )
   // 2
   flag.BoolVar(&c.login, "l", false, "device login")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   flag.IntVar(&c.season, "s", 0, "season")
   // 4
   flag.StringVar(&c.edit, "e", "", "edit ID")
   // 5
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   // 1
   if c.initiate {
      return c.do_initiate()
   }
   // 2
   if c.login {
      return c.do_login()
   }
   // 3
   if c.address != "" {
      return c.do_address()
   }
   // 4
   if c.edit != "" {
      return c.do_edit()
   }
   // 5
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
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
   var show hboMax.ShowKey
   err := show.Parse(c.address)
   if err != nil {
      return err
   }
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   var videos *hboMax.Videos
   if c.season >= 1 {
      videos, err = cache.Login.Season(&show, c.season)
   } else {
      videos, err = cache.Login.Movie(&show)
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
   cache.Mpd, err = cache.Playback.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Mpd.Body, cache.Mpd.Url)
}

type command struct {
   job    maya.PlayReadyJob
   market string
   name   string
   season int
   // 1
   initiate bool
   // 2
   login bool
   // 3
   address string
   // 4
   edit string
   // 5
   dash string
}

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Playback.PlayReady(data)
   }
   return c.job.DownloadDash(cache.Mpd.Body, cache.Mpd.Url, c.dash)
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

type user_cache struct {
   Login    *hboMax.Login
   Mpd      *hboMax.Mpd
   Playback *hboMax.Playback
   St       *hboMax.St
}

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
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

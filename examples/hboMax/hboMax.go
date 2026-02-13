package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/hboMax"
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
      &c.market, "m", hboMax.Markets[0], fmt.Sprint(hboMax.Markets),
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
   flag.IntVar(&c.job.Threads, "t", 2, "threads")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
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
   maya.Usage([][]string{
      {"i", "m"},
      {"l"},
      {"a", "s"},
      {"e"},
      {"d", "t", "C", "E"},
   })
   return nil
}

type user_cache struct {
   Login    *hboMax.Login
   Dash     *hboMax.Dash
   Playback *hboMax.Playback
   St       *hboMax.St
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
   name string
   // 1
   initiate bool
   market   string
   // 2
   login bool
   // 3
   address string
   season  int
   // 4
   edit string
   // 5
   dash string
   job  maya.PlayReadyJob
}

func (c *command) do_initiate() error {
   var st hboMax.St
   err := st.Fetch()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, &user_cache{St: &st})
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
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   cache.Login, err = cache.St.Login()
   if err != nil {
      return err
   }
   return maya.Write(c.name, cache)
}

func (c *command) do_address() error {
   var show hboMax.ShowKey
   err := show.Parse(c.address)
   if err != nil {
      return err
   }
   cache, err := maya.Read[user_cache](c.name)
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
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   cache.Playback, err = cache.Login.PlayReady(c.edit)
   if err != nil {
      return err
   }
   cache.Dash, err = cache.Playback.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Playback.PlayReady(data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

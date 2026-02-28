package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (c *command) run() error {
   c.cache.Init("SL3000")
   c.job.CertificateChain = c.cache.Join("CertificateChain")
   c.job.EncryptSignKey = c.cache.Join("EncryptSignKey")
   c.cache.Init("hboMax")
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
   return maya.Usage([][]string{
      {"i", "m"},
      {"l"},
      {"a", "s"},
      {"e"},
      {"d", "C", "E"},
   })
}

func (c *command) do_address() error {
   var show hboMax.ShowKey
   err := show.Parse(c.address)
   if err != nil {
      return err
   }
   var login hboMax.Login
   err = c.cache.Get("Login", &login)
   if err != nil {
      return err
   }
   var videos *hboMax.Videos
   if c.season >= 1 {
      videos, err = login.Season(&show, c.season)
   } else {
      videos, err = login.Movie(&show)
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
   var login hboMax.Login
   err := c.cache.Get("Login", &login)
   if err != nil {
      return err
   }
   playback, err := login.PlayReady(c.edit)
   if err != nil {
      return err
   }
   err = c.cache.Set("Playback", playback)
   if err != nil {
      return err
   }
   dash, err := playback.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func (c *command) do_initiate() error {
   var st hboMax.St
   err := st.Fetch()
   if err != nil {
      return err
   }
   err = c.cache.Set("St", st)
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

func (c *command) do_dash() error {
   var playback hboMax.Playback
   err := c.cache.Get("Playback", &playback)
   if err != nil {
      return err
   }
   c.job.Send = playback.PlayReady
   var dash hboMax.Dash
   err = c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func (c *command) do_login() error {
   var st hboMax.St
   err := c.cache.Get("St", &st)
   if err != nil {
      return err
   }
   login, err := st.Login()
   if err != nil {
      return err
   }
   return c.cache.Set("Login", login)
}

type command struct {
   cache maya.Cache
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

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".mp4"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/hboMax"
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
   c.job.CertificateChain = filepath.Join(cache, "/SL3000/CertificateChain")
   c.job.EncryptSignKey = filepath.Join(cache, "/SL3000/EncryptSignKey")
   c.name = cache + "/hboMax/userCache.xml"
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
   // 1
   usage("i", "m")
   // 2
   usage("l")
   // 3
   usage("a", "s")
   // 4
   usage("e")
   // 5
   usage("d", "t", "C", "E")
   return nil
}

// 1
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

// 2
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

// 3
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

// 4
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

// 5
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

package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/hboMax"
   "fmt"
   "log"
)

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

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = cache.Playback.PlayReady
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
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

func main() {
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
   proxy string
   // 4
   address string
   season  int
   // 4, 5
   get_proxy bool
   // 5
   edit string
   // 6
   dash string
   job  maya.PlayReadyJob
}

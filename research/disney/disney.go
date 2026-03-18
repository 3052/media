package main

import (
   "41.neocities.org/maya"
   "flag"
)

func (c *client) do() error {
   err := cache.Setup("rosso/disney.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c)
   // 1
   playReady := StringVar(&c.Job.PlayReady, "PR", "PlayReady")
   // 2
   email := StringVar(&c.Email, "e", "email")
   // 3
   passcode := StringVar(&c.passcode, "p", "passcode")
   // 4
   profile := StringVar(&c.profile_id, "P", "profile ID")
   // 5
   refresh := BoolVar(new(bool), "r", "refresh")
   // 6
   address := StringVar(&c.address, "a", "address")
   // 7
   season := StringVar(&c.season_id, "s", "season ID")
   // 8
   media := StringVar(&c.media_id, "m", "media")
   // 9
   hls := IntVar(&c.hls_id, "h", "HLS ID")
   flag.Parse()
   switch {
   case IsSet(playReady):
      return cache.Write(c)
   case IsSet(email):
      return c.do_email()
   }
   if err != nil {
      return err
   }
   switch {
   case IsSet(passcode):
      return c.do_passcode()
   case IsSet(profile):
      return c.do_profile_id()
   case IsSet(refresh):
      return c.do_refresh()
   case IsSet(address):
      return c.do_address()
   case IsSet(season):
      return c.do_season_id()
   case IsSet(media):
      return c.do_media_id()
   case IsSet(hls):
      return c.Job.DownloadHls(
         c.Hls.Body, c.Hls.Url, c.hls_id, c.Token.PlayReady,
      )
   }
   return maya.Usage([][]string{
      {"PR"},
      {"e"},
      {"p"},
      {"P"},
      {"r"},
      {"a"},
      {"s"},
      {"m"},
      {"h"},
   })
}

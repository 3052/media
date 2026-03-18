package main

import (
   "41.neocities.org/maya"
   "flag"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4,*.mp4a")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   err := cache.Setup("rosso/disney.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c)
   // 1
   playReady := maya.StringVar(&c.Job.PlayReady, "PR", "PlayReady")
   // 2
   email := maya.StringVar(&c.Email, "e", "email")
   // 3
   passcode := maya.StringVar(&c.passcode, "p", "passcode")
   // 4
   profile_id := maya.StringVar(&c.profile_id, "P", "profile ID")
   // 5
   refresh := maya.BoolVar(new(bool), "r", "refresh")
   // 6
   address := maya.StringVar(&c.address, "a", "address")
   // 7
   season_id := maya.StringVar(&c.season_id, "s", "season ID")
   // 8
   media_id := maya.StringVar(&c.media_id, "m", "media ID")
   // 9
   hls_id := maya.IntVar(&c.hls_id, "h", "HLS ID")
   flag.Parse()
   switch {
   case maya.IsSet(playReady):
      return cache.Write(c)
   case maya.IsSet(email):
      return c.do_email()
   }
   if err != nil {
      return err
   }
   switch {
   case maya.IsSet(passcode):
      return c.do_passcode()
   case maya.IsSet(profile_id):
      return c.do_profile_id()
   case maya.IsSet(refresh):
      return c.do_refresh()
   case maya.IsSet(address):
      return c.do_address()
   case maya.IsSet(season_id):
      return c.do_season_id()
   case maya.IsSet(media_id):
      return c.do_media_id()
   case maya.IsSet(hls_id):
      return c.Job.DownloadHls(
         c.Hls.Body, c.Hls.Url, c.hls_id, c.Token.PlayReady,
      )
   }
   return maya.Usage([][]*flag.Flag{
      {playReady},
      {email},
      {passcode},
      {profile_id},
      {refresh},
      {address},
      {season_id},
      {media_id},
      {hls_id},
   })
}

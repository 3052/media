package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "flag"
   "fmt"
   "log"
)

type client struct {
   Hls   *disney.Hls
   Token *disney.Token
   // 1
   Email string
   // 2
   passcode string
   // 3
   profile_id string
   // 4
   refresh bool
   // 5
   address string
   // 6
   season_id string
   // 7
   media_id string
   // 8
   job maya.Job
   // 9
   hls_id int
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4,*.mp4a")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email() error {
   c.Token = &disney.Token{}
   err := c.Token.RegisterDevice()
   if err != nil {
      return err
   }
   request_otp, err := c.Token.RequestOtp(c.Email)
   if err != nil {
      return err
   }
   fmt.Println(request_otp)
   return cache.Write(c)
}

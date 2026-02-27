package main

import (
   "41.neocities.org/maya"
   "log"
   "net/http"
   "path"
)

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      switch path.Ext(req.URL.Path) {
      case ".mp4", ".mp4a":
         return "", false
      }
      return "", true
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   name string
   // 1
   email    string
   password string
   // 2
   profile_id string
   // 3
   address string
   // 4
   season_id string
   // 5
   media_id string
   // 6
   hls int
   job maya.PlayReadyJob
}

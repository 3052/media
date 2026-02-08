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
)

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

type command struct {
   name   string
   // 1
   initiate bool
   market string
   // 2
   login bool
   // 3
   address string
   season int
   // 4
   edit string
   // 5
   dash string
   job    maya.PlayReadyJob
}

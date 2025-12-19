package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/rakuten"
   "encoding/xml"
   "log"
   "net/http"
   "os"
   "path"
)

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
      switch path.Ext(req.URL.Path) {
      case ".isma", ".ismv":
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
   config maya.Config
   name   string
   language string
   // 1
   show string
   // 2
   season string
   // 3
   episode string
   // 4
   movie string
   // 5
   dash     string
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
   Movie   *rakuten.Movie
   Mpd *rakuten.Mpd
   TvShow  *rakuten.TvShow
}

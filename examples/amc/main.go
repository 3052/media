package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/amc"
   "encoding/xml"
   "log"
   "net/http"
   "os"
   "path"
)

type user_cache struct {
   Client *amc.Client
   Header http.Header
   Dash    *amc.Dash
   Source *amc.Source
}

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
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

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4f" {
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
   name     string
   // 1
   email    string
   password string
   // 2
   refresh  bool
   // 3
   series   int
   // 4
   season   int
   // 5
   episode  int
   // 6
   dash     string
   job   maya.WidevineJob
}

package main

import (
   "41.neocities.org/maya"
   "log"
)

type client struct {
   cache maya.Cache
   // 1
   proxy string
   // 2
   initiate bool
   market   string
   // 3
   login bool
   // 4
   address string
   season  int
   // 5
   edit string
   // 6
   dash string
   job  maya.PlayReadyJob
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

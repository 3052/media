package main

import (
   "154.pages.dev/log"
   "154.pages.dev/rosso"
   "flag"
)

type flags struct {
   address string
   email string
   h rosso.HttpStream
   hls_index int
   password string
   v log.Level
}

func main() {
   var f flags
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.email, "e", "", "email")
   flag.IntVar(&f.hls_index, "h", -1, "HLS index")
   flag.StringVar(&f.password, "p", "", "password")
   flag.TextVar(&f.v.Level, "v", f.v.Level, "log level")
   flag.Parse()
   log.TransportInfo()
   log.Handler(f.v)
   if f.email != "" {
      err := f.profile()
      if err != nil {
         panic(err)
      }
   } else if f.address != "" {
      err := f.download()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

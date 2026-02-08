package main

import (
   "flag"
   "fmt"
)

func usage(names ...string) {
   for _, name := range names {
      look := flag.Lookup(name)
      fmt.Printf("-%v %v\n", look.Name, look.Usage)
      if look.DefValue != "" {
         fmt.Printf("\tdefault %v\n", look.DefValue)
      }
   }
   fmt.Println()
}

package main

import (
   "41.neocities.org/media/canal"
   "41.neocities.org/media/internal"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) do_season() error {
   assets1, err := assets(test1.id, 1)
   if err != nil {
      t.Fatal(err)
   }
   for _, asset1 := range assets1 {
      fmt.Print("\n", &asset1, "\n")
   }
}

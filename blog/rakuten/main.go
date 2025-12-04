// main.go
package main

import (
   "flag"
   "log"
)

func main() {
   // 1. Setup Flags
   var movie, show, season, language, item, dash bool

   flag.BoolVar(&movie, "movie", false, "Process movie")
   flag.BoolVar(&show, "show", false, "Process show")
   flag.BoolVar(&season, "season", false, "Process season")
   flag.BoolVar(&language, "language", false, "Process language")
   flag.BoolVar(&item, "item", false, "Is item")
   flag.BoolVar(&dash, "dash", false, "Is dash")

   flag.Parse()

   // 2. Run Logic (defined in actions.go)
   didRun, err := runLogic(movie, show, season, language, item, dash)

   // STATE 2: Flag provided, error occurred (Exit Code 1)
   if err != nil {
      log.Fatalf("Error detected: %v", err)
   }

   // STATE 3: No valid flags provided (Exit Code 0)
   if !didRun {
      flag.Usage()
   }

   // STATE 1: Success (Exit Code 0)
   // Falls through to the end naturally
}

package main

import (
   "errors"
   "flag"
   "fmt"
   "log"
)

type Config struct {
   Movie    bool
   Show     bool
   Season   bool
   Language bool
   Item     bool
   Dash     bool
}

// Run defines flags, parses them, and executes logic.
// Returns immediately after the first matching action is found.
func (c *Config) Run() (bool, error) {
   // --- 1. SETUP FLAGS ---
   flag.BoolVar(&c.Movie, "movie", false, "Process movie")
   flag.BoolVar(&c.Show, "show", false, "Process show")
   flag.BoolVar(&c.Season, "season", false, "Process season")
   flag.BoolVar(&c.Language, "language", false, "Process language")
   flag.BoolVar(&c.Item, "item", false, "Is item")
   flag.BoolVar(&c.Dash, "dash", false, "Is dash")

   flag.Parse()

   // --- 2. EXECUTE LOGIC ---

   // Movie Logic
   if c.Movie {
      // Run logic and return result immediately
      return true, do_movie()
   }

   // Show Logic
   if c.Show {
      return true, do_show()
   }

   // Season Logic
   if c.Season {
      return true, do_season()
   }

   // Language Logic
   if c.Language {
      if c.Dash {
         return true, do_download(c.Item)
      }
      return true, do_representations(c.Item)
   }

   // State 3: No flags triggered any action
   return false, nil
}

// --- ACTION FUNCTIONS ---

func do_movie() error {
   fmt.Println("Action: Processing Movie")
   return nil
}

func do_show() error {
   fmt.Println("Action: Processing Show")
   return nil
}

func do_season() error {
   return errors.New("database connection failed")
}

func do_download(isItem bool) error {
   target := "Movie"
   if isItem {
      target = "Episode"
   }
   fmt.Printf("Action: Downloading MP4 (%s)\n", target)
   return nil
}

func do_representations(isItem bool) error {
   target := "Movie"
   if isItem {
      target = "Episode"
   }
   fmt.Printf("Action: Processing MPD Representations (%s)\n", target)
   return nil
}
func main() {
   var cfg Config

   // Calls Setup AND Execution
   didRun, err := cfg.Run()

   // STATE 2: Error occurred (Exit Code 1)
   if err != nil {
      log.Fatalf("Error detected: %v", err)
   }

   // STATE 3: No flags provided (Exit Code 0)
   if !didRun {
      flag.Usage()
   }

   // STATE 1: Success (Exit Code 0)
}

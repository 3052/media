// actions.go
package main

import (
   "errors"
   "fmt"
)

// --- LOGIC CONTROLLER ---
// Returns:
// 1. bool: true if an action was attempted
// 2. error: non-nil if that action failed
func runLogic(movie, show, season, language, item, dash bool) (bool, error) {
   actionPerformed := false
   var err error

   // 1. Movie Logic
   if movie {
      actionPerformed = true
      if err = do_movie(); err != nil {
         return true, err
      }
   }

   // 2. Show Logic
   if show {
      actionPerformed = true
      if err = do_show(); err != nil {
         return true, err
      }
   }

   // 3. Season Logic
   if season {
      actionPerformed = true
      if err = do_season(); err != nil {
         return true, err
      }
   }

   // 4. Language Logic
   if language {
      actionPerformed = true
      if item {
         if dash {
            err = episode_MP4()
         } else {
            err = episode_MPD()
         }
      } else {
         if dash {
            err = movie_MP4()
         } else {
            err = movie_MPD()
         }
      }

      if err != nil {
         return true, err
      }
   }

   return actionPerformed, nil
}

// --- MOCK ACTION FUNCTIONS ---

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

func episode_MP4() error { fmt.Println("Action: Episode MP4"); return nil }
func episode_MPD() error { fmt.Println("Action: Episode MPD"); return nil }
func movie_MP4() error   { fmt.Println("Action: Movie MP4"); return nil }
func movie_MPD() error   { fmt.Println("Action: Movie MPD"); return nil }

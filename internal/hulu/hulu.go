package main

import (
   "41.neocities.org/dash"
   "41.neocities.org/media/hulu"
   "fmt"
   "io"
   "net/http"
   "os"
   "sort"
)

func (f *flags) download() error {
   data, err := os.ReadFile(f.home + "/hulu.txt")
   if err != nil {
      return err
   }
   var auth hulu.Authenticate
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   deep, err := auth.DeepLink(&f.entity)
   if err != nil {
      return err
   }
   play, err := auth.Playlist(deep)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.StreamUrl)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   reps, err := dash.Unmarshal(data, resp.Request.URL)
   if err != nil {
      return err
   }
   sort.Slice(reps, func(i, j int) bool {
      return reps[i].Bandwidth < reps[j].Bandwidth
   })
   for _, rep := range reps {
      if rep.GetAdaptationSet().GetPeriod().Id != "content-0" {
         continue
      }
      if rep.Width < f.min_width {
         if rep.GetMimeType() == "video/mp4" {
            continue
         }
      }
      switch f.representation {
      case "":
         fmt.Print(&rep, "\n\n")
      case rep.Id:
         f.s.Namer, err = auth.Details(deep)
         if err != nil {
            return err
         }
         f.s.Wrapper = play
         return f.s.Download(rep)
      }
   }
   return nil
}

func (f *flags) authenticate() error {
   data, err := hulu.Authenticate{}.Marshal(f.email, f.password)
   if err != nil {
      return err
   }
   return os.WriteFile(f.home + "/hulu.txt", data, os.ModePerm)
}

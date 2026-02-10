package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/itv"
   "errors"
   "fmt"
   "path"
)

func (f *flags) download() error {
   var id itv.LegacyId
   err := id.Set(path.Base(f.address))
   if err != nil {
      return err
   }
   play, err := id.Playlist()
   if err != nil {
      return err
   }
   file, ok := play.Resolution1080()
   if !ok {
      return errors.New("resolution 1080")
   }
   represents, err := internal.Mpd(file.Href)
   if err != nil {
      return err
   }
   for _, represent := range represents {
      switch f.representation {
      case "":
         fmt.Print(&represent, "\n\n")
      case represent.Id:
         f.s.Client = file
         return f.s.Download(&represent)
      }
   }
   return nil
}

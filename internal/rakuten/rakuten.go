package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/rakuten"
   "errors"
   "fmt"
)

func (f *flags) download() error {
   class, ok := f.address.ClassificationId()
   if !ok {
      return errors.New("Address.ClassificationId")
   }
   var content *rakuten.Content
   if f.address.SeasonId != "" {
      season, err := f.address.Season(class)
      if err != nil {
         return err
      }
      content, ok = season.Content(&f.address)
      if !ok {
         return errors.New("Season.Content")
      }
   } else {
      var err error
      content, err = f.address.Movie(class)
      if err != nil {
         return err
      }
   }
   stream := content.Streamings()
   stream.Fhd()
   info, err := stream.Info(f.language, class)
   if err != nil {
      return err
   }
   represents, err := internal.Mpd(info)
   if err != nil {
      return err
   }
   for _, represent := range represents {
      switch f.representation {
      case "":
         fmt.Print(&represent, "\n\n")
      case represent.Id:
         stream.Hd()
         info, err = stream.Info(f.language, class)
         if err != nil {
            return err
         }
         f.s.Client = info
         return f.s.Download(&represent)
      }
   }
   return nil
}

func (f *flags) do_language() error {
   class, ok := f.address.ClassificationId()
   if !ok {
      return errors.New("Address.ClassificationId")
   }
   var content *rakuten.Content
   if f.address.SeasonId != "" {
      season, err := f.address.Season(class)
      if err != nil {
         return err
      }
      content, ok = season.Content(&f.address)
      if !ok {
         return errors.New("Season.Content")
      }
   } else {
      var err error
      content, err = f.address.Movie(class)
      if err != nil {
         return err
      }
   }
   fmt.Println(content)
   return nil
}

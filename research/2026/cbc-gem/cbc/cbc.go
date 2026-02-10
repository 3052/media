package main

import (
   "154.pages.dev/media/cbc"
   "os"
)

func (f flags) download() error {
   home, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   var profile cbc.GemProfile
   profile.Raw, err = os.ReadFile(home + "/cbc/profile.json")
   if err != nil {
      return err
   }
   var catalog cbc.GemCatalog
   catalog.New(f.address)
   f.h.Name = catalog.StructuredMetadata
   item, ok := catalog.Item()
   if ok {
      media, err := profile.Media(item)
      if err != nil {
         return err
      }
      master, err := f.h.HlsMaster(media.URL)
      if err != nil {
         return err
      }
      return f.h.HLS(master, f.hls_index)
   }
   return nil
}

func (f flags) profile() error {
   var token cbc.LoginToken
   token.New(f.email, f.password)
   profile, err := token.Profile()
   if err != nil {
      return err
   }
   home, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   return os.WriteFile(home + "/cbc/profile.json", profile.Raw, 0666)
}

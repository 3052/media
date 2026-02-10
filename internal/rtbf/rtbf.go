package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/rtbf"
   "errors"
   "fmt"
   "os"
)

func (f *flags) authenticate() error {
   data, err := rtbf.Login{}.Marshal(f.email, f.password)
   if err != nil {
      return err
   }
   return os.WriteFile(f.home+"/rtbf.txt", data, os.ModePerm)
}

func (f *flags) download() error {
   data, err := os.ReadFile(f.home + "/rtbf.txt")
   if err != nil {
      return err
   }
   var login rtbf.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   jwt, err := login.Jwt()
   if err != nil {
      return err
   }
   gigya, err := jwt.Login()
   if err != nil {
      return err
   }
   content, err := f.address.Content()
   if err != nil {
      return err
   }
   asset_id, ok := content.GetAssetId()
   if !ok {
      return errors.New("Content.GetAssetId")
   }
   title, err := gigya.Entitlement(asset_id)
   if err != nil {
      return err
   }
   format, ok := title.Dash()
   if !ok {
      return errors.New("Entitlement.Dash")
   }
   represents, err := internal.Mpd(format)
   if err != nil {
      return err
   }
   for _, represent := range represents {
      switch f.representation {
      case "":
         fmt.Print(&represent, "\n\n")
      case represent.Id:
         f.s.Client = title
         return f.s.Download(&represent)
      }
   }
   return nil
}

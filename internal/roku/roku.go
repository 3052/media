package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/roku"
   "41.neocities.org/x/http"
   "flag"
   "fmt"
   "log"
   "os"
   "path/filepath"
)

func (f *flags) download() error {
   if f.representation != "" {
      f.e.Widevine = play.Widevine()
      return f.e.Download(f.home, f.representation)
   }
   var code *roku.Code
   if f.token_read {
      data, err := os.ReadFile(f.home + "/roku.txt")
      if err != nil {
         return err
      }
      code = &roku.Code{}
      err = code.Unmarshal(data)
      if err != nil {
         return err
      }
   }
   var token roku.Token
   data, err := token.Marshal(code)
   if err != nil {
      return err
   }
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   play, err := token.Playback(f.roku)
   if err != nil {
      return err
   }
   resp, err := play.Mpd()
   if err != nil {
      return err
   }
   return internal.Mpd(resp, f.home)
}

func (f *flags) write_token() error {
   data, err := os.ReadFile("activation.txt")
   if err != nil {
      return err
   }
   var activation roku.Activation
   err = activation.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = os.ReadFile("token.txt")
   if err != nil {
      return err
   }
   var token roku.Token
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = roku.Code{}.Marshal(&activation, &token)
   if err != nil {
      return err
   }
   return os.WriteFile(f.home+"/roku.txt", data, os.ModePerm)
}

func write_code() error {
   var token roku.Token
   data, err := token.Marshal(nil)
   if err != nil {
      return err
   }
   err = os.WriteFile("token.txt", data, os.ModePerm)
   if err != nil {
      return err
   }
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   var activation roku.Activation
   data, err = activation.Marshal(&token)
   if err != nil {
      return err
   }
   err = os.WriteFile("activation.txt", data, os.ModePerm)
   if err != nil {
      return err
   }
   err = activation.Unmarshal(data)
   if err != nil {
      return err
   }
   fmt.Println(activation)
   return nil
}

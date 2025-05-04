package main

import (
   "41.neocities.org/media/cineMember"
   "41.neocities.org/media/internal"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) do_address() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/cineMember/Play")
      if err != nil {
         return err
      }
      var play cineMember.Play
      err = play.Unmarshal(data)
      if err != nil {
         return err
      }
      title, _ := play.Dash()
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return title.Widevine(data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/cineMember/User")
   if err != nil {
      return err
   }
   var user cineMember.User
   err = user.Unmarshal(data)
   if err != nil {
      return err
   }
   article, err := f.address.Article()
   if err != nil {
      return err
   }
   asset, ok := article.Film()
   if !ok {
      return errors.New(".Film()")
   }
   data, err = user.Play(article, asset)
   if err != nil {
      return err
   }
   var play cineMember.Play
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.media + "/cineMember/Play", data)
   if err != nil {
      return err
   }
   title, ok := play.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(title.Manifest)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/cineMember/Play")
      if err != nil {
         return err
      }
      var play cineMember.Play
      err = play.Unmarshal(data)
      if err != nil {
         return err
      }
      title, _ := play.Dash()
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return title.Widevine(data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/cineMember/User")
   if err != nil {
      return err
   }
   var user cineMember.User
   err = user.Unmarshal(data)
   if err != nil {
      return err
   }
   article, err := f.address.Article()
   if err != nil {
      return err
   }
   asset, ok := article.Film()
   if !ok {
      return errors.New(".Film()")
   }
   data, err = user.Play(article, asset)
   if err != nil {
      return err
   }
   var play cineMember.Play
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.media + "/cineMember/Play", data)
   if err != nil {
      return err
   }
   title, ok := play.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(title.Manifest)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}

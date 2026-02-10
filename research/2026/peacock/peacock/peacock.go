package main

import (
   "154.pages.dev/media/peacock"
   "errors"
   "fmt"
   "net/http"
   "os"
)

func (f flags) download() error {
   text, err := os.ReadFile(f.home + "/peacock.json")
   if err != nil {
      return err
   }
   var sign peacock.SignIn
   sign.Unmarshal(text)
   auth, err := sign.Auth()
   if err != nil {
      return err
   }
   video, err := auth.Video(f.peacock)
   if err != nil {
      return err
   }
   akamai, ok := video.Akamai()
   if !ok {
      return errors.New("peacock.VideoPlayout.Akamai")
   }
   req, err := http.NewRequest("", akamai, nil)
   if err != nil {
      return err
   }
   media, err := f.s.DASH(req)
   if err != nil {
      return err
   }
   for _, medium := range media {
      if medium.ID == f.representation {
         var node peacock.QueryNode
         err := node.New(f.peacock)
         if err != nil {
            return err
         }
         f.s.Name = node
         f.s.Poster = video
         return f.s.Download(medium)
      }
   }
   // 2 MPD all
   for i, medium := range media {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(medium)
   }
   return nil
}

func (f flags) authenticate() error {
   var sign peacock.SignIn
   err := sign.New(f.email, f.password)
   if err != nil {
      return err
   }
   text, err := sign.Marshal()
   if err != nil {
      return err
   }
   return os.WriteFile(f.home + "/peacock.json", text, 0666)
}

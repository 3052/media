package main

import (
   "154.pages.dev/media/criterion"
   "154.pages.dev/media/internal"
   "errors"
   "fmt"
   "net/http"
   "os"
   "path"
)

func (f flags) authenticate() error {
   var token criterion.AuthToken
   err := token.New(f.email, f.password)
   if err != nil {
      return err
   }
   return os.WriteFile(f.home + "/criterion.txt", token.Marshal(), 0666)
}

func (f flags) download() error {
   text, err := os.ReadFile(f.home + "/criterion.txt")
   if err != nil {
      return err
   }
   var token criterion.AuthToken
   err = token.Unmarshal(text)
   if err != nil {
      return err
   }
   item, err := token.Video(path.Base(f.address))
   if err != nil {
      return err
   }
   files, err := token.Files(item)
   if err != nil {
      return err
   }
   file, ok := files.Dash()
   if !ok {
      return errors.New("VideoFiles.Dash")
   }
   req, err := http.NewRequest("", file.Links.Source.Href, nil)
   if err != nil {
      return err
   }
   reps, err := internal.Dash(req)
   if err != nil {
      return err
   }
   for _, rep := range reps {
      switch f.representation {
      case "":
         fmt.Print(rep, "\n\n")
      case rep.Id:
         f.s.Name = item
         f.s.Poster = file
         return f.s.Download(rep)
      }
   }
   return nil
}

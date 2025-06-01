package main

import (
   "41.neocities.org/media/rakuten"
   "41.neocities.org/net"
   "errors"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flag_set) do_address() error {
   var path rakuten.Path
   path.New(f.address)
   class, ok := path.ClassificationId()
   if !ok {
      return errors.New(".ClassificationId()")
   }
   var content *rakuten.Content
   if path.SeasonId != "" {
      data, err := path.Season(class)
      if err != nil {
         return err
      }
      var season rakuten.Season
      err = season.Unmarshal(data)
      if err != nil {
         return err
      }
      err = write_file(f.media+"/rakuten/Season", data)
      if err != nil {
         return err
      }
      content, ok = season.Content(&path)
      if !ok {
         return errors.New(".Content")
      }
   } else {
      data, err := path.Movie(class)
      if err != nil {
         return err
      }
      content = &rakuten.Content{}
      err = content.Unmarshal(data)
      if err != nil {
         return err
      }
      err = write_file(f.media+"/rakuten/Content", data)
      if err != nil {
         return err
      }
   }
   fmt.Println(content)
   return nil
}

func (f *flag_set) do_language() error {
   var path rakuten.Path
   path.New(f.address)
   class, ok := path.ClassificationId()
   if !ok {
      return errors.New(".ClassificationId()")
   }
   var content *rakuten.Content
   if path.SeasonId != "" {
      data, err := os.ReadFile(f.media + "/rakuten/Season")
      if err != nil {
         return err
      }
      var season rakuten.Season
      err = season.Unmarshal(data)
      if err != nil {
         return err
      }
      content, ok = season.Content(&path)
      if !ok {
         return errors.New(".Content")
      }
   } else {
      data, err := os.ReadFile(f.media + "/rakuten/Content")
      if err != nil {
         return err
      }
      content = &rakuten.Content{}
      err = content.Unmarshal(data)
      if err != nil {
         return err
      }
   }
   streaming := content.Streamings()
   streaming.Fhd()
   info, err := streaming.Info(f.language, class)
   if err != nil {
      return err
   }
   resp, err := http.Get(info.Url)
   if err != nil {
      return err
   }
   streaming.Hd()
   info, err = streaming.Info(f.language, class)
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return info.License(data)
   }
   return f.filters.Filter(resp, &f.cdm)
}

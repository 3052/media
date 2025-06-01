package main

import (
   "41.neocities.org/media/rakuten"
   "fmt"
   "net/http"
   "os"
)

func (f *flag_set) do_content() error {
   data, err := os.ReadFile(f.media + "/rakuten/Address")
   if err != nil {
      return err
   }
   var address rakuten.Address
   err = address.Set(string(data))
   if err != nil {
      return err
   }
   info, err := address.Info(f.content, f.language, rakuten.Wvm, rakuten.Fhd)
   if err != nil {
      return err
   }
   resp, err := http.Get(info.Url)
   if err != nil {
      return err
   }
   info, err = address.Info(f.content, f.language, rakuten.Wvm, rakuten.Hd)
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return info.License(data)
   }
   return f.filters.Filter(resp, &f.cdm)
}

// print movie
func (f *flag_set) do_movie() error {
   var address rakuten.Address
   err := address.Set(f.movie)
   if err != nil {
      return err
   }
   content, err := address.Movie()
   if err != nil {
      return err
   }
   fmt.Println(content)
   return nil
}

// print seasons
func (f *flag_set) do_show() error {
   var address rakuten.Address
   err := address.Set(f.show)
   if err != nil {
      return err
   }
   err = write_file(f.media+"/rakuten/Address", []byte(f.show))
   if err != nil {
      return err
   }
   seasons, err := address.Seasons()
   if err != nil {
      return err
   }
   for i, season := range seasons {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&season)
   }
   return nil
}

// print episodes
func (f *flag_set) do_season() error {
   data, err := os.ReadFile(f.media + "/rakuten/Address")
   if err != nil {
      return err
   }
   var address rakuten.Address
   err = address.Set(string(data))
   if err != nil {
      return err
   }
   contents, err := address.Episodes(f.season)
   if err != nil {
      return err
   }
   for i, content := range contents {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&content)
   }
   return nil
}

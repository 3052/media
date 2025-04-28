package main

import (
   "41.neocities.org/media/amc"
   "41.neocities.org/media/internal"
   "errors"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) do_refresh() error {
   data, err := os.ReadFile(f.media + "/amc/Auth")
   if err != nil {
      return err
   }
   var auth amc.Auth
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = auth.Refresh()
   if err != nil {
      return err
   }
   err = write_file(f.media+"/amc/Auth", data)
   if err != nil {
      return err
   }
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = auth.Playback(f.episode)
   if err != nil {
      return err
   }
   err = write_file(f.media+"/amc/Playback", data)
   if err != nil {
      return err
   }
   var play amc.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   source, ok := play.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(source.Src)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_episode() error {
   data, err := os.ReadFile(f.media + "/amc/Auth")
   if err != nil {
      return err
   }
   var auth amc.Auth
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = auth.Refresh()
   if err != nil {
      return err
   }
   err = write_file(f.media+"/amc/Auth", data)
   if err != nil {
      return err
   }
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = auth.Playback(f.episode)
   if err != nil {
      return err
   }
   err = write_file(f.media+"/amc/Playback", data)
   if err != nil {
      return err
   }
   var play amc.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   source, ok := play.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(source.Src)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
   data, err := os.ReadFile(f.media + "/amc/Playback")
   if err != nil {
      return err
   }
   var play amc.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   source, _ := play.Dash()
   f.e.Widevine = func(data []byte) ([]byte, error) {
      return play.Widevine(source, data)
   }
   return f.e.Download(f.media+"/Mpd", f.dash)
}

///

func (f *flags) do_season() error {
   season, err := amc.SeasonEpisodes(f.season)
   if err != nil {
      return err
   }
   var line bool
   for episode := range season.Episodes() {
      if line {
         fmt.Println()
      } else {
         line = true
      }
      fmt.Println(&episode.Properties.Metadata)
   }
   return nil
}

func (f *flags) do_series() error {
   series, err := amc.SeriesDetail(f.series)
   if err != nil {
      return err
   }
   var line bool
   for season := range series.Seasons() {
      if line {
         fmt.Println()
      } else {
         line = true
      }
      fmt.Println(&season.Properties.Metadata)
   }
   return nil
}

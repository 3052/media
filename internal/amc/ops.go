package main

import (
   "encoding/json"
   "errors"
   "fmt"
   "log"
   "os"

   "41.neocities.org/media/amc"
   "41.neocities.org/net"
)

func (r *runner) read() (*amc.Cache, error) {
   data, err := os.ReadFile(r.cache + "/amc/Cache")
   if err != nil {
      return nil, err
   }
   var cache amc.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return nil, err
   }
   return &cache, nil
}

func (r *runner) write(cache *amc.Cache) error {
   data, err := json.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", r.cache + "/amc/Cache")
   return os.WriteFile(r.cache + "/amc/Cache", data, os.ModePerm)
}

func (r *runner) do_auth() error {
   var client amc.Client
   err := client.Unauth()
   if err != nil {
      return err
   }
   err = client.Login(r.email, r.password)
   if err != nil {
      return err
   }
   return r.write(&amc.Cache{Client: &client})
}

func (r *runner) do_refresh() error {
   cache, err := r.read()
   if err != nil {
      return err
   }
   err = cache.Client.Refresh()
   if err != nil {
      return err
   }
   return r.write(cache)
}

func (r *runner) do_series() error {
   cache, err := r.read()
   if err != nil {
      return err
   }
   series, err := cache.Client.SeriesDetail(r.series)
   if err != nil {
      return err
   }
   seasons, err := series.ExtractSeasons()
   if err != nil {
      return err
   }
   for i, season := range seasons {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(season)
   }
   return nil
}

func (r *runner) do_season() error {
   cache, err := r.read()
   if err != nil {
      return err
   }
   season, err := cache.Client.SeasonEpisodes(r.season)
   if err != nil {
      return err
   }
   episodes, err := season.ExtractEpisodes()
   if err != nil {
      return err
   }
   for i, episode := range episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(episode)
   }
   return nil
}

func (r *runner) do_episode() error {
   cache, err := r.read()
   if err != nil {
      return err
   }
   cache.Header, cache.Source, err = cache.Client.Playback(r.episode)
   if err != nil {
      return err
   }
   source, ok := amc.Dash(cache.Source)
   if !ok {
      return errors.New("amc.Dash")
   }
   err = source.Mpd(cache)
   if err != nil {
      return err
   }
   err = r.write(cache)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}

func (r *runner) do_dash() error {
   cache, err := r.read()
   if err != nil {
      return err
   }
   r.config.Send = func(data []byte) ([]byte, error) {
      source, _ := amc.Dash(cache.Source)
      return source.Widevine(cache.Header, data)
   }
   return r.config.Download(cache.MpdBody, cache.Mpd, r.dash)
}

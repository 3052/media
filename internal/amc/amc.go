package main

import (
   "41.neocities.org/media/amc"
   "41.neocities.org/net"
   "errors"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
)

func main() {
   http.DefaultTransport = &http.Transport{
      Proxy: func(req *http.Request) (*url.URL, error) {
         log.Println(req.Method, req.URL)
         return nil, nil
      },
   }
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.email != "" {
      if set.password != "" {
         err = set.do_email()
      }
   } else if set.refresh {
      err = set.do_refresh()
   } else if set.series >= 1 {
      err = set.do_series()
   } else if set.season >= 1 {
      err = set.do_season()
   } else if set.episode >= 1 {
      err = set.do_episode()
   } else {
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

type flag_set struct {
   cdm     net.Cdm
   filters net.Filters
   media   string
   ///////////////////
   email    string
   password string
   ///////////////
   refresh bool
   ////////////
   series int64
   ////////////
   season int64
   /////////////
   episode int64
}

func (f *flag_set) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.cdm.ClientId = f.media + "/client_id.bin"
   f.cdm.PrivateKey = f.media + "/private_key.pem"
   f.filters = net.Filters{
      {BitrateStart: 100_000, BitrateEnd: 200_000},
      {BitrateStart: 2_000_000, BitrateEnd: 4_000_000},
   }
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.cdm.ClientId, "client", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.cdm.PrivateKey, "key", f.cdm.PrivateKey, "private key")
   /////////////////////////////////////////////////////////////////////////
   flag.StringVar(&f.email, "email", "", "email")
   flag.StringVar(&f.password, "password", "", "password")
   ///////////////////////////////////////////////////////
   flag.BoolVar(&f.refresh, "r", false, "refresh")
   //////////////////////////////////////////////////
   flag.Int64Var(&f.series, "series", 0, "series ID")
   //////////////////////////////////////////////////
   flag.Int64Var(&f.season, "s", 0, "season ID")
   ////////////////////////////////////////////////////////
   flag.Int64Var(&f.episode, "e", 0, "episode or movie ID")
   flag.Parse()
   return nil
}

func (f *flag_set) do_episode() error {
   data, err := os.ReadFile(f.media + "/amc/Auth")
   if err != nil {
      return err
   }
   var auth amc.Auth
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = auth.Playback(f.episode)
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
   f.cdm.License = func(data []byte) ([]byte, error) {
      return play.License(source, data)
   }
   return f.filters.Filter(resp, &f.cdm)
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) do_refresh() error {
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
   return write_file(f.media+"/amc/Auth", data)
}

func (f *flag_set) do_email() error {
   var auth amc.Auth
   err := auth.Unauth()
   if err != nil {
      return err
   }
   data, err := auth.Login(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.media+"/amc/Auth", data)
}

func (f *flag_set) do_season() error {
   data, err := os.ReadFile(f.media + "/amc/Auth")
   if err != nil {
      return err
   }
   var auth amc.Auth
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   season, err := auth.SeasonEpisodes(f.season)
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

func (f *flag_set) do_series() error {
   data, err := os.ReadFile(f.media + "/amc/Auth")
   if err != nil {
      return err
   }
   var auth amc.Auth
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   series, err := auth.SeriesDetail(f.series)
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

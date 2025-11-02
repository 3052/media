package main

import (
   "41.neocities.org/media/amc"
   "41.neocities.org/net"
   "errors"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flag_set) do_series() error {
   data, err := os.ReadFile(f.cache + "/amc/Client")
   if err != nil {
      return err
   }
   var client amc.Client
   err = client.Unmarshal(data)
   if err != nil {
      return err
   }
   series, err := client.SeriesDetail(f.series)
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

type flag_set struct {
   cache    string
   config   net.Config
   email    string
   episode  int64
   filters  net.Filters
   password string
   refresh  bool
   season   int64
   series   int64
}
func (f *flag_set) do_episode() error {
   data, err := os.ReadFile(f.cache + "/amc/Client")
   if err != nil {
      return err
   }
   var client amc.Client
   err = client.Unmarshal(data)
   if err != nil {
      return err
   }
   play, err := client.Playback(f.episode)
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
   f.config.Send = func(data []byte) ([]byte, error) {
      return play.Widevine(source, data)
   }
   return f.filters.Filter(resp, &f.config)
}

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.config.ClientId = f.cache + "/L3/client_id.bin"
   f.config.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.StringVar(&f.email, "E", "", "email")
   flag.StringVar(&f.password, "P", "", "password")
   flag.Int64Var(&f.series, "S", 0, "series ID")
   flag.StringVar(&f.config.ClientId, "c", f.config.ClientId, "client ID")
   flag.Int64Var(&f.episode, "e", 0, "episode or movie ID")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.config.PrivateKey, "p", f.config.PrivateKey, "private key")
   flag.BoolVar(&f.refresh, "r", false, "refresh")
   flag.Int64Var(&f.season, "s", 0, "season ID")
   flag.Parse()
   return nil
}
func main() {
   http.DefaultTransport = &amc.Transport
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   switch {
   case set.email_password():
      err = set.do_auth()
   case set.episode >= 1:
      err = set.do_episode()
   case set.refresh:
      err = set.do_refresh()
   case set.season >= 1:
      err = set.do_season()
   case set.series >= 1:
      err = set.do_series()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) do_refresh() error {
   data, err := os.ReadFile(f.cache + "/amc/Client")
   if err != nil {
      return err
   }
   var client amc.Client
   err = client.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = client.Refresh()
   if err != nil {
      return err
   }
   return write_file(f.cache+"/amc/Client", data)
}

func (f *flag_set) do_auth() error {
   var client amc.Client
   err := client.Unauth()
   if err != nil {
      return err
   }
   data, err := client.Login(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/amc/Client", data)
}

func (f *flag_set) email_password() bool {
   if f.email != "" {
      if f.password != "" {
         return true
      }
   }
   return false
}

func (f *flag_set) do_season() error {
   data, err := os.ReadFile(f.cache + "/amc/Client")
   if err != nil {
      return err
   }
   var client amc.Client
   err = client.Unmarshal(data)
   if err != nil {
      return err
   }
   season, err := client.SeasonEpisodes(f.season)
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

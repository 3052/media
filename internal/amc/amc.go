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

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func main() {
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) != ".m4f" {
         return ""
      }
      return "L"
   })
   var program runner
   err := program.run()
   if err != nil {
      log.Fatal(err)
   }
   
}

///

func (r *runner) do_episode() error {
   data, err := os.ReadFile(r.cache + "/amc/Client")
   if err != nil {
      return err
   }
   var client amc.Client
   err = client.Unmarshal(data)
   if err != nil {
      return err
   }
   header, sources, err := client.Playback(r.episode)
   if err != nil {
      return err
   }
   source, ok := amc.Dash(sources)
   if !ok {
      return errors.New("amc.Dash")
   }
   resp, err := http.Get(source.Src)
   if err != nil {
      return err
   }
   r.config.Send = func(data []byte) ([]byte, error) {
      return amc.Widevine(header, source, data)
   }
   return r.filters.Filter(resp, &r.config)
}

func (r *runner) run() error {
   var err error
   r.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   r.cache = filepath.ToSlash(r.cache)
   r.config.ClientId = r.cache + "/L3/client_id.bin"
   r.config.PrivateKey = r.cache + "/L3/private_key.pem"
   flag.StringVar(&r.email, "E", "", "email")
   flag.StringVar(&r.password, "P", "", "password")
   flag.Int64Var(&r.series, "S", 0, "series ID")
   flag.StringVar(&r.config.ClientId, "c", r.config.ClientId, "client ID")
   flag.Int64Var(&r.episode, "e", 0, "episode or movie ID")
   flag.Var(&r.filters, "f", net.FilterUsage)
   flag.StringVar(&r.config.PrivateKey, "p", r.config.PrivateKey, "private key")
   flag.BoolVar(&r.refresh, "r", false, "refresh")
   flag.Int64Var(&r.season, "s", 0, "season ID")
   flag.Parse()
   return nil
   switch {
   case program.email_password():
      err = program.do_auth()
   case program.episode >= 1:
      err = program.do_episode()
   case program.refresh:
      err = program.do_refresh()
   case program.season >= 1:
      err = program.do_season()
   case program.series >= 1:
      err = program.do_series()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (r *runner) do_refresh() error {
   data, err := os.ReadFile(r.cache + "/amc/Client")
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
   return write_file(r.cache+"/amc/Client", data)
}

func (r *runner) do_auth() error {
   var client amc.Client
   err := client.Unauth()
   if err != nil {
      return err
   }
   data, err := client.Login(r.email, r.password)
   if err != nil {
      return err
   }
   return write_file(r.cache+"/amc/Client", data)
}

func (r *runner) email_password() bool {
   if r.email != "" {
      if r.password != "" {
         return true
      }
   }
   return false
}

func (r *runner) do_season() error {
   data, err := os.ReadFile(r.cache + "/amc/Client")
   if err != nil {
      return err
   }
   var client amc.Client
   err = client.Unmarshal(data)
   if err != nil {
      return err
   }
   season, err := client.SeasonEpisodes(r.season)
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

func (r *runner) do_series() error {
   data, err := os.ReadFile(r.cache + "/amc/Client")
   if err != nil {
      return err
   }
   var client amc.Client
   err = client.Unmarshal(data)
   if err != nil {
      return err
   }
   series, err := client.SeriesDetail(r.series)
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

type runner struct {
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

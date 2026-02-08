package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/amc"
   "encoding/xml"
   "errors"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

type user_cache struct {
   Client *amc.Client
   Header http.Header
   Dash    *amc.Dash
   Source *amc.Source
}

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func read(name string) (*user_cache, error) {
   data, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   cache := &user_cache{}
   err = xml.Unmarshal(data, cache)
   if err != nil {
      return nil, err
   }
   return cache, nil
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4f" {
         return ""
      }
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   name     string
   // 1
   email    string
   password string
   // 2
   refresh  bool
   // 3
   series   int
   // 4
   season   int
   // 5
   episode  int
   // 6
   dash     string
   job   maya.WidevineJob
}

func (c *command) do_episode() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   header, sources, err := cache.Client.Playback(c.episode)
   if err != nil {
      return err
   }
   cache.Header = header
   source, ok := amc.GetDash(sources)
   if !ok {
      return errors.New("amc.Dash")
   }
   cache.Source = source
   cache.Dash, err = source.Dash()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Source.Widevine(cache.Header, data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   c.name = cache + "/amc/userCache.xml"
   c.job.ClientId = filepath.Join(cache, "/L3/client_id.bin")
   c.job.PrivateKey = filepath.Join(cache, "/L3/private_key.pem")
   // 1
   flag.StringVar(&c.email, "E", "", "email")
   flag.StringVar(&c.password, "P", "", "password")
   // 2
   flag.BoolVar(&c.refresh, "r", false, "refresh")
   // 3
   flag.IntVar(&c.series, "S", 0, "series ID")
   // 4
   flag.IntVar(&c.season, "s", 0, "season ID")
   // 5
   flag.IntVar(&c.episode, "e", 0, "episode or movie ID")
   // 6
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   // 1
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   // 2
   if c.refresh {
      return c.do_refresh()
   }
   // 3
   if c.series >= 1 {
      return c.do_series()
   }
   // 4
   if c.season >= 1 {
      return c.do_season()
   }
   // 5
   if c.episode >= 1 {
      return c.do_episode()
   }
   // 6
   if c.dash != "" {
      return c.do_dash()
   }
   usage("E", "P")
   usage("r")
   usage("S")
   usage("s")
   usage("e")
   usage("d", "c", "p")
   return nil
}

func (c *command) do_email_password() error {
   var client amc.Client
   err := client.Unauth()
   if err != nil {
      return err
   }
   err = client.Login(c.email, c.password)
   if err != nil {
      return err
   }
   return write(c.name, &user_cache{Client: &client})
}

func (c *command) do_refresh() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   err = cache.Client.Refresh()
   if err != nil {
      return err
   }
   return write(c.name, cache)
}

func (c *command) do_series() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   series, err := cache.Client.SeriesDetail(c.series)
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

func (c *command) do_season() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   season, err := cache.Client.SeasonEpisodes(c.season)
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

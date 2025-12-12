package main

import (
   "41.neocities.org/media/amc"
   "41.neocities.org/net"
   "encoding/json"
   "errors"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

type user_cache struct {
   Header http.Header
   Mpd    struct {
      Body []byte
      Url  *url.URL
   }
   Source []amc.Source
   Client amc.Client
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
   return write(c.name, &user_cache{Client: client})
}

func write(name string, cache *user_cache) error {
   data, err := json.Marshal(cache)
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
   err = json.Unmarshal(data, cache)
   if err != nil {
      return nil, err
   }
   return cache, nil
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
func main() {
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
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

func (c *command) do_episode() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   cache.Header, cache.Source, err = cache.Client.Playback(c.episode)
   if err != nil {
      return err
   }
   source, ok := amc.Dash(cache.Source)
   if !ok {
      return errors.New("amc.Dash")
   }
   cache.Mpd.Url, cache.Mpd.Body, err = source.Mpd()
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return net.Representations(cache.Mpd.Url, cache.Mpd.Body)
}

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.config.Send = func(data []byte) ([]byte, error) {
      source, _ := amc.Dash(cache.Source)
      return source.Widevine(cache.Header, data)
   }
   return c.config.Download(cache.Mpd.Url, cache.Mpd.Body, c.dash)
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/amc/user_cache.json"

   flag.StringVar(&c.email, "E", "", "email")
   flag.StringVar(&c.password, "P", "", "password")
   flag.Int64Var(&c.series, "S", 0, "series ID")
   flag.StringVar(&c.config.ClientId, "c", c.config.ClientId, "client ID")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.Int64Var(&c.episode, "e", 0, "episode or movie ID")
   flag.StringVar(&c.config.PrivateKey, "p", c.config.PrivateKey, "private key")
   flag.BoolVar(&c.refresh, "r", false, "refresh")
   flag.Int64Var(&c.season, "s", 0, "season ID")
   flag.Parse()

   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.refresh {
      return c.do_refresh()
   }
   if c.series >= 1 {
      return c.do_series()
   }
   if c.season >= 1 {
      return c.do_season()
   }
   if c.episode >= 1 {
      return c.do_episode()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
}

type command struct {
   config   net.Config
   dash     string
   email    string
   episode  int64
   name     string
   password string
   refresh  bool
   season   int64
   series   int64
}

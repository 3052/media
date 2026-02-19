package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/amc"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func main() {
   maya.SetTransport(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".m4f"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   name string
   // 1
   email    string
   password string
   // 2
   refresh bool
   // 3
   series int
   // 4
   season int
   // 5
   episode int
   // 6
   dash string
   job  maya.WidevineJob
}

func (c *command) do_episode() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   sources, header, err := cache.Client.Playback(c.episode)
   if err != nil {
      return err
   }
   cache.DataSource, err = amc.GetDash(sources)
   if err != nil {
      return err
   }
   cache.BcJwt = amc.BcJwt(header)
   cache.Dash, err = cache.DataSource.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.name = cache + "/amc/userCache.xml"
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
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
   return maya.Usage([][]string{
      {"E", "P"},
      {"r"},
      {"S"},
      {"s"},
      {"e"},
      {"d", "c", "p"},
   })
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
   return maya.Write(c.name, &user_cache{Client: &client})
}

func (c *command) do_refresh() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   err = cache.Client.Refresh()
   if err != nil {
      return err
   }
   return maya.Write(c.name, cache)
}

func (c *command) do_series() error {
   cache, err := maya.Read[user_cache](c.name)
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

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.DataSource.Widevine(cache.BcJwt, data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

type user_cache struct {
   BcJwt      string
   Client     *amc.Client
   Dash       *amc.Dash
   DataSource *amc.DataSource
}

func (c *command) do_season() error {
   cache, err := maya.Read[user_cache](c.name)
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

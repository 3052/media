package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (c *command) do_episode() error {
   var client amc.Client
   err := c.cache.Get("Client", &client)
   if err != nil {
      return err
   }
   sources, header, err := client.Playback(c.episode)
   if err != nil {
      return err
   }
   data_source, err := amc.GetDash(sources)
   if err != nil {
      return err
   }
   err = c.cache.Set("DataSource", data_source)
   if err != nil {
      return err
   }
   err = c.cache.Set("BcJwt", amc.BcJwt(header))
   if err != nil {
      return err
   }
   dash, err := data_source.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("amc")
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
   return c.cache.Set("Client", client)
}

func (c *command) do_refresh() error {
   var client amc.Client
   err := c.cache.Get("Client", &client)
   if err != nil {
      return err
   }
   err = client.Refresh()
   if err != nil {
      return err
   }
   return c.cache.Set("Client", client)
}

func (c *command) do_series() error {
   var client amc.Client
   err := c.cache.Get("Client", &client)
   if err != nil {
      return err
   }
   series, err := client.SeriesDetail(c.series)
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
   var source amc.DataSource
   err := c.cache.Get("DataSource", &source)
   if err != nil {
      return err
   }
   var bc_jwt string
   err = c.cache.Get("BcJwt", &bc_jwt)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return source.Widevine(bc_jwt, data)
   }
   var dash amc.Dash
   err = c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func (c *command) do_season() error {
   var client amc.Client
   err := c.cache.Get("Client", &client)
   if err != nil {
      return err
   }
   season, err := client.SeasonEpisodes(c.season)
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
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".m4f"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   cache maya.Cache
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

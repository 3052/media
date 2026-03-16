package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "flag"
   "fmt"
   "log"
)

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4f")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   BcJwt      string
   Client     *amc.Client
   Dash       *amc.Dash
   DataSource *amc.DataSource
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
   Job maya.Job
   // 7
   dash_id string
}

///

func (c *client) do() error {
   err := cache.Setup("rosso/amc.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c, true)
   if err != nil {
      return err
   }
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
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "c", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "p", job.PrivateKey, "private key")
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
   if c.dash_id != "" {
      return c.do_dash_id()
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

func (c *client) do_email_password() error {
   var err error
   c.Client, err = amc.Unauth()
   if err != nil {
      return err
   }
   err = c.Client.Login(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_refresh() error {
   return cache.Update(c, func() error {
      return c.Client.Refresh()
   })
}

func (c *client) do_series() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   series, err := c.Client.SeriesDetail(c.series)
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

func (c *client) do_season() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   season, err := c.Client.SeasonEpisodes(c.season)
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
func (c *client) do_episode() error {
   err := cache.Update(c, func() error {
      sources, header, err := c.Client.Playback(c.episode)
      if err != nil {
         return err
      }
      c.DataSource, err = amc.GetDash(sources)
      if err != nil {
         return err
      }
      c.Dash, err = c.DataSource.Dash()
      if err != nil {
         return err
      }
      c.BcJwt = amc.BcJwt(header)
      return nil
   })
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = func(data []byte) ([]byte, error) {
      return c.DataSource.Widevine(c.BcJwt, data)
   }
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
}

package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "flag"
   "fmt"
   "log"
)

func (c *client) do_refresh(err error) error {
   if err != nil {
      return err
   }
   err = c.Client.Refresh()
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_series(err error) error {
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

func (c *client) do_season(err error) error {
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

func (c *client) do_episode(err error) error {
   if err != nil {
      return err
   }
   sources, header, err := c.Client.Playback(c.episode)
   if err != nil {
      return err
   }
   c.DataSource, err = amc.GetDash(sources)
   if err != nil {
      return err
   }
   c.BcJwt = amc.BcJwt(header)
   c.Dash, err = c.DataSource.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

type client struct {
   BcJwt      string
   Client     *amc.Client
   Dash       *amc.Dash
   DataSource *amc.DataSource
   // 1
   Job maya.Job
   // 2
   email    string
   password string
   // 3
   refresh bool
   // 4
   series int
   // 5
   season int
   // 6
   episode int
   // 7
   dash_id string
}

func (c *client) do_dash_id(err error) error {
   if err != nil {
      return err
   }
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, func(data []byte) ([]byte, error) {
         return c.DataSource.Widevine(c.BcJwt, data)
      },
   )
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4f")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   err := cache.Setup("rosso/amc.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c)
   // 1
   flag.StringVar(&c.Job.Widevine, "w", c.Job.Widevine, "Widevine")
   // 2
   flag.StringVar(&c.email, "E", "", "email")
   flag.StringVar(&c.password, "P", "", "password")
   // 3
   flag.BoolVar(&c.refresh, "r", false, "refresh")
   // 4
   flag.IntVar(&c.series, "s", 0, "series ID")
   // 5
   flag.IntVar(&c.season, "S", 0, "season ID")
   // 6
   flag.IntVar(&c.episode, "e", 0, "episode or movie ID")
   // 7
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   set := maya.Parse()
   switch {
   case set["w"]:
      return cache.Write(c)
   case set["E"] && set["P"]:
      return c.do_email_password()
   case set["r"]:
      return c.do_refresh(err)
   case set["s"]:
      return c.do_series(err)
   case set["S"]:
      return c.do_season(err)
   case set["e"]:
      return c.do_episode(err)
   case set["d"]:
      return c.do_dash_id(err)
   }
   return maya.Usage([][]string{
      {"w"},
      {"E", "P"},
      {"r"},
      {"s"},
      {"S"},
      {"e"},
      {"d"},
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

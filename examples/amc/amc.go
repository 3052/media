package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "flag"
   "fmt"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/amc.xml")
   if err != nil {
      return err
   }
   read_err := cache.Read(c)
   // 1
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   // 2
   email := maya.StringVar(&c.email, "E", "email")
   password := maya.StringVar(&c.password, "P", "password")
   // 3
   refresh := maya.BoolVar(new(bool), "r", "refresh")
   // 4
   series := maya.IntVar(&c.series, "s", "series ID")
   // 5
   season := maya.IntVar(&c.season, "S", "season ID")
   // 6
   episode := maya.IntVar(&c.episode, "e", "episode or movie ID")
   // 7
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case len(set) == 0:
      return maya.Usage([][]*flag.Flag{
         {widevine},
         {email, password},
         {refresh},
         {series},
         {season},
         {episode},
         {dash_id},
      })
   case set[widevine]:
      return cache.Write(c)
   case set[email] && set[password]:
      return c.do_email_password()
   case read_err != nil:
      return read_err
   case set[refresh]:
      return c.do_refresh()
   case set[series]:
      return c.do_series()
   case set[season]:
      return c.do_season()
   case set[episode]:
      return c.do_episode()
   case set[dash_id]:
      return c.do_dash_id()
   }
   return nil
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
   err := c.Client.Refresh()
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_series() error {
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
   // 4
   series int
   // 5
   season int
   // 6
   episode int
   // 7
   dash_id string
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id,
      func(data []byte) ([]byte, error) {
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

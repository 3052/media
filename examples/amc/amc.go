package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amc"
   "flag"
   "fmt"
   "log"
)

func (c *client) do() error {
   c.job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   c.job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Setup("rosso/amc.xml")
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

func (c *client) do_email_password() error {
   var client amc.Client
   err := client.Unauth()
   if err != nil {
      return err
   }
   err = client.Login(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Write(saved_state{Client: &client})
}

func (c *client) do_refresh() error {
   var state saved_state
   return c.cache.Update(&state, func() error {
      return state.Client.Refresh()
   })
}

func (c *client) do_series() error {
   var state saved_state
   err := c.cache.Read(&state)
   if err != nil {
      return err
   }
   series, err := state.Client.SeriesDetail(c.series)
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
   var state saved_state
   err := c.cache.Read(&state)
   if err != nil {
      return err
   }
   season, err := state.Client.SeasonEpisodes(c.season)
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

type client struct {
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
func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4f")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type saved_state struct {
   BcJwt      string
   Client     *amc.Client
   Dash       *amc.Dash
   DataSource *amc.DataSource
}

func (c *client) do_episode() error {
   var state saved_state
   err := c.cache.Update(&state, func() error {
      sources, header, err := state.Client.Playback(c.episode)
      if err != nil {
         return err
      }
      state.DataSource, err = amc.GetDash(sources)
      if err != nil {
         return err
      }
      state.Dash, err = state.DataSource.Dash()
      if err != nil {
         return err
      }
      state.BcJwt = amc.BcJwt(header)
      return nil
   })
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

func (c *client) do_dash() error {
   var state saved_state
   err := c.cache.Read(&state)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return state.DataSource.Widevine(state.BcJwt, data)
   }
   return c.job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash)
}

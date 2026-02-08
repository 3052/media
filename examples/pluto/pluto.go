package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/pluto"
   "encoding/xml"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

type user_cache struct {
   Dash    *pluto.Dash
   Series *pluto.Series
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
         return ""
      }
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
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

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   c.name = cache + "/pluto/userCache.xml"
   c.job.ClientId = filepath.Join(cache, "/L3/client_id.bin")
   c.job.PrivateKey = filepath.Join(cache, "/L3/private_key.pem")
   // 1
   flag.StringVar(&c.movie, "m", "", "movie ID")
   // 2
   flag.StringVar(&c.show, "s", "", "show ID")
   // 3
   flag.StringVar(&c.episode, "e", "", "episode ID")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   // 1
   if c.movie != "" {
      return c.do_movie()
   }
   // 2
   if c.show != "" {
      return c.do_show()
   }
   // 3
   if c.episode != "" {
      return c.do_episode()
   }
   // 4
   if c.dash != "" {
      return c.do_dash()
   }
   maya.Usage(
      []string{"m"},
      []string{"s"},
      []string{"e"},
      []string{"d", "c", "p"},
   )
   return nil
}

type command struct {
   name    string
   // 1
   movie   string
   // 2
   show    string
   // 3
   episode string
   // 4
   dash    string
   job  maya.WidevineJob
}

///

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.job.Send = pluto.Widevine
   return c.job.Download(cache.Dash.Url, cache.Dash.Body, c.dash)
}

func (c *command) do_show() error {
   var series pluto.Series
   err := series.Fetch(c.show)
   if err != nil {
      return err
   }
   fmt.Println(&series.Vod[0])
   return write(c.name, &user_cache{Series: &series})
}

func (c *command) do_movie() error {
   var series pluto.Series
   err := series.Fetch(c.movie)
   if err != nil {
      return err
   }
   var dash pluto.Dash
   err = dash.Fetch(series.GetMovieURL())
   if err != nil {
      return err
   }
   err = write(c.name, &user_cache{Dash: &dash})
   if err != nil {
      return err
   }
   return maya.Representations(dash.Url, dash.Body)
}

func (c *command) do_episode() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   link, err := cache.Series.GetEpisodeURL(c.episode)
   if err != nil {
      return err
   }
   var dash pluto.Dash
   err = dash.Fetch(link)
   if err != nil {
      return err
   }
   cache.Dash = &dash
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.Representations(dash.Url, dash.Body)
}

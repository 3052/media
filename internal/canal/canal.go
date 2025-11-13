package main

import (
   "41.neocities.org/media/canal"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.config.ClientId = f.cache + "/L3/client_id.bin"
   f.config.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.StringVar(&f.config.ClientId, "C", f.config.ClientId, "client ID")
   flag.StringVar(&f.config.PrivateKey, "P", f.config.PrivateKey, "private key")
   flag.IntVar(&f.config.Threads, "T", 2, "threads")
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.email, "e", "", "email")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.password, "p", "", "password")
   flag.BoolVar(&f.refresh, "r", false, "refresh")
   flag.Int64Var(&f.season, "s", 0, "season")
   flag.StringVar(&f.tracking_id, "t", "", "tracking ID")
   flag.Parse()
   return nil
}

func (f *flag_set) do_address() error {
   tracking_id, err := canal.TrackingId(f.address)
   if err != nil {
      return err
   }
   fmt.Println("tracking id =", tracking_id)
   return nil
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}
func (f *flag_set) do_session() error {
   var ticket canal.Ticket
   err := ticket.Fetch()
   if err != nil {
      return err
   }
   token, err := ticket.Token(f.email, f.password)
   if err != nil {
      return err
   }
   data, err := canal.FetchSession(token.SsoToken)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/canal/Session", data)
}

func (f *flag_set) do_refresh() error {
   data, err := os.ReadFile(f.cache + "/canal/Session")
   if err != nil {
      return err
   }
   var session canal.Session
   err = session.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = canal.FetchSession(session.SsoToken)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/canal/Session", data)
}

func main() {
   http.DefaultTransport = &canal.Transport
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.address != "" {
      err = set.do_address()
   } else if set.tracking_id != "" {
      if set.season >= 1 {
         err = set.do_season()
      } else {
         err = set.do_episode_movie()
      }
   } else if set.email_password() {
      err = set.do_session()
   } else if set.refresh {
      err = set.do_refresh()
   } else {
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flag_set) email_password() bool {
   if f.email != "" {
      if f.password != "" {
         return true
      }
   }
   return false
}

type flag_set struct {
   address     string
   cache       string
   config      net.Config
   email       string
   filters     net.Filters
   password    string
   refresh     bool
   season      int64
   tracking_id string
}

func (f *flag_set) do_season() error {
   data, err := os.ReadFile(f.cache + "/canal/Session")
   if err != nil {
      return err
   }
   var session canal.Session
   err = session.Unmarshal(data)
   if err != nil {
      return err
   }
   episodes, err := session.Episodes(f.tracking_id, f.season)
   if err != nil {
      return err
   }
   for i, episode := range episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&episode)
   }
   return nil
}

func (f *flag_set) do_episode_movie() error {
   data, err := os.ReadFile(f.cache + "/canal/Session")
   if err != nil {
      return err
   }
   var session canal.Session
   err = session.Unmarshal(data)
   if err != nil {
      return err
   }
   player, err := session.Player(f.tracking_id)
   if err != nil {
      return err
   }
   resp, err := http.Get(player.Url)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return player.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
}

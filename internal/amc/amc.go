package main

import (
   "41.neocities.org/media/amc"
   "41.neocities.org/net"
   "errors"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.e.ClientId = f.media + "/client_id.bin"
   f.e.PrivateKey = f.media + "/private_key.pem"
   f.bandwidth.Value = []int64{127000, 2588000}
   return nil
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flags) do_refresh() error {
   data, err := os.ReadFile(f.media + "/amc/Auth")
   if err != nil {
      return err
   }
   var auth amc.Auth
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = auth.Refresh()
   if err != nil {
      return err
   }
   return write_file(f.media+"/amc/Auth", data)
}

func (f *flags) do_email() error {
   var auth amc.Auth
   err := auth.Unauth()
   if err != nil {
      return err
   }
   data, err := auth.Login(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.media+"/amc/Auth", data)
}

func (f *flags) do_season() error {
   data, err := os.ReadFile(f.media + "/amc/Auth")
   if err != nil {
      return err
   }
   var auth amc.Auth
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   season, err := auth.SeasonEpisodes(f.season)
   if err != nil {
      return err
   }
   var line bool
   for episode := range season.Episodes() {
      if line {
         fmt.Println()
      } else {
         line = true
      }
      fmt.Println(&episode.Properties.Metadata)
   }
   return nil
}

func (f *flags) do_series() error {
   data, err := os.ReadFile(f.media + "/amc/Auth")
   if err != nil {
      return err
   }
   var auth amc.Auth
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   series, err := auth.SeriesDetail(f.series)
   if err != nil {
      return err
   }
   var line bool
   for season := range series.Seasons() {
      if line {
         fmt.Println()
      } else {
         line = true
      }
      fmt.Println(&season.Properties.Metadata)
   }
   return nil
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.e.ClientId, "client", f.e.ClientId, "client ID")
   flag.StringVar(&f.e.PrivateKey, "key", f.e.PrivateKey, "private key")
   //////////////////////////////////////////////////////////////////////////////
   flag.StringVar(&f.email, "email", "", "email")
   flag.StringVar(&f.password, "password", "", "password")
   //////////////////////////////////////////////////////////////////////////////
   flag.BoolVar(&f.refresh, "r", false, "refresh")
   //////////////////////////////////////////////////////////////////////////////
   flag.Int64Var(&f.series, "series", 0, "series ID")
   //////////////////////////////////////////////////////////////////////////////
   flag.Int64Var(&f.season, "s", 0, "season ID")
   //////////////////////////////////////////////////////////////////////////////
   flag.Int64Var(&f.episode, "e", 0, "episode or movie ID")
   flag.Float64Var(&f.tolerance, "t", 0.2, "tolerance")
   flag.Func("b", "bandwidth", func(data string) error {
      return f.bandwidth.Set(data)
   })
   flag.Parse()
   if f.email != "" {
      if f.password != "" {
         err = f.do_email()
      }
   } else if f.refresh {
      err = f.do_refresh()
   } else if f.series >= 1 {
      err = f.do_series()
   } else if f.season >= 1 {
      err = f.do_season()
   } else if f.episode >= 1 {
      err = f.do_episode()
   } else {
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flags) do_episode() error {
   data, err := os.ReadFile(f.media + "/amc/Auth")
   if err != nil {
      return err
   }
   var auth amc.Auth
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = auth.Playback(f.episode)
   if err != nil {
      return err
   }
   var play amc.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   source, ok := play.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(source.Src)
   if err != nil {
      return err
   }
   f.e.Widevine = func(data []byte) ([]byte, error) {
      return play.Widevine(source, data)
   }
   return f.e.Tolerance(resp, f.bandwidth.Value, f.tolerance)
}

type flags struct {
   e        net.License
   media    string
   /////////////////////////
   email    string
   password string
   /////////////////////////
   refresh  bool
   /////////////////////////
   series   int64
   /////////////////////////
   season   int64
   /////////////////////////
   episode  int64
   tolerance float64
   bandwidth net.Bandwidth
}

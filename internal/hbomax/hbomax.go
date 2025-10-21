package main

import (
   "41.neocities.org/media/hboMax"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
)

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.filters.Values = []net.Filter{
      { Bandwidth: 7_000_000 },
      { Bandwidth: 200_000 },
   }
   f.cache = filepath.ToSlash(f.cache)
   f.config.CertificateChain = f.cache + "/SL3000/CertificateChain"
   f.config.EncryptSignKey = f.cache + "/SL3000/EncryptSignKey"
   flag.StringVar(&f.config.CertificateChain, "C", f.config.CertificateChain, "certificate chain")
   flag.StringVar(&f.config.EncryptSignKey, "E", f.config.EncryptSignKey, "encrypt sign key")
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.edit, "e", "", "edit ID")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.BoolVar(&f.initiate, "i", false, "device initiate")
   flag.BoolVar(&f.login, "l", false, "device login")
   flag.IntVar(&f.season, "s", 0, "season")
   flag.Parse()
   return nil
}

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.cache + "/hboMax/Login")
   if err != nil {
      return err
   }
   var login hboMax.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   show_id, err := hboMax.ShowId(f.address)
   if err != nil {
      return err
   }
   var videos *hboMax.Videos
   if f.season >= 1 {
      videos, err = login.Season(show_id, f.season)
   } else {
      videos, err = login.Movie(show_id)
   }
   if err != nil {
      return err
   }
   videos.EpisodeMovie()
   for i, video := range videos.Included {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(video)
   }
   return nil
}

func main() {
   log.SetFlags(log.Ltime)
   http.DefaultTransport = &http.Transport{
      Proxy: func(req *http.Request) (*url.URL, error) {
         if filepath.Ext(req.URL.Path) != ".mp4" {
            log.Println(req.Method, req.URL)
         }
         return http.ProxyFromEnvironment(req)
      },
   }
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   switch {
   case set.address != "":
      err = set.do_address()
   case set.edit != "":
      err = set.do_edit()
   case set.initiate:
      err = set.do_initiate()
   case set.login:
      err = set.do_login()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flag_set) do_edit() error {
   data, err := os.ReadFile(f.cache + "/hboMax/Login")
   if err != nil {
      return err
   }
   var login hboMax.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = login.PlayReady(f.edit)
   if err != nil {
      return err
   }
   var play hboMax.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.Mpd())
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return play.PlayReady(data)
   }
   return f.filters.Filter(resp, &f.config)
}

type flag_set struct {
   address  string
   cache    string
   config   net.Config
   edit     string
   filters  net.Filters
   initiate bool
   login    bool
   season   int
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) do_initiate() error {
   var st hboMax.St
   err := st.New()
   if err != nil {
      return err
   }
   log.Println("Create", f.cache+"/hboMax/St")
   file, err := os.Create(f.cache + "/hboMax/St")
   if err != nil {
      return err
   }
   defer file.Close()
   _, err = fmt.Fprint(file, st)
   if err != nil {
      return err
   }
   initiate, err := st.Initiate()
   if err != nil {
      return err
   }
   fmt.Println(initiate)
   return nil
}

func (f *flag_set) do_login() error {
   data, err := os.ReadFile(f.cache + "/hboMax/St")
   if err != nil {
      return err
   }
   var st hboMax.St
   err = st.Set(string(data))
   if err != nil {
      return err
   }
   data, err = st.Login()
   if err != nil {
      return err
   }
   return write_file(f.cache+"/hboMax/Login", data)
}

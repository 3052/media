package main

import (
   "41.neocities.org/media/hboMax"
   "41.neocities.org/net"
   "encoding/json"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.config.CertificateChain = f.cache + "/SL3000/CertificateChain"
   f.config.EncryptSignKey = f.cache + "/SL3000/EncryptSignKey"
   flag.StringVar(&f.config.CertificateChain, "C", f.config.CertificateChain, "certificate chain")
   flag.StringVar(&f.config.EncryptSignKey, "E", f.config.EncryptSignKey, "encrypt sign key")
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.dash, "d", "", "DASH ID")
   flag.StringVar(&f.edit, "e", "", "edit ID")
   flag.BoolVar(&f.initiate, "i", false, "device initiate")
   flag.BoolVar(&f.login, "l", false, "device login")
   flag.IntVar(&f.season, "s", 0, "season")
   flag.Parse()
   return nil
}

func main() {
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".mp4" {
         return ""
      }
      return "LP"
   })
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   switch {
   case set.initiate:
      err = set.do_initiate()
   case set.login:
      err = set.do_login()
   case set.address != "":
      err = set.do_address()
   case set.edit != "":
      err = set.do_edit()
   case set.dash != "":
      err = set.do_dash()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

type flag_set struct {
   cache    string
   config   net.Config
   // 1
   initiate bool
   // 2
   login    bool
   // 3
   address  string
   season   int
   // 4
   edit     string
   // 5
   dash     string
}

func (f *flag_set) do_initiate() error {
   var st hboMax.St
   err := st.Fetch()
   if err != nil {
      return err
   }
   data, err := json.Marshal(hboMax.Cache{St: &st})
   if err != nil {
      return err
   }
   err = write_file(f.cache + "/hboMax/Cache", data)
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
   data, err := os.ReadFile(f.cache + "/hboMax/Cache")
   if err != nil {
      return err
   }
   var cache hboMax.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   cache.Login, err = cache.St.Login()
   if err != nil {
      return err
   }
   data, err = json.Marshal(cache)
   if err != nil {
      return err
   }
   return write_file(f.cache + "/hboMax/Cache", data)
}

func (f *flag_set) do_address() error {
   show_id, err := hboMax.ExtractId(f.address)
   if err != nil {
      return err
   }
   data, err := os.ReadFile(f.cache + "/hboMax/Cache")
   if err != nil {
      return err
   }
   var cache hboMax.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   var videos *hboMax.Videos
   if f.season >= 1 {
      videos, err = cache.Login.Season(show_id, f.season)
   } else {
      videos, err = cache.Login.Movie(show_id)
   }
   if err != nil {
      return err
   }
   videos.FilterAndSort()
   for i, video := range videos.Included {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(video)
   }
   return nil
}

func (f *flag_set) do_edit() error {
   data, err := os.ReadFile(f.cache + "/hboMax/Cache")
   if err != nil {
      return err
   }
   var cache hboMax.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   cache.Playback, err = cache.Login.PlayReady(f.edit)
   if err != nil {
      return err
   }
   err = cache.Playback.FetchManifest(&cache)
   if err != nil {
      return err
   }
   data, err = json.Marshal(cache)
   if err != nil {
      return err
   }
   err = write_file(f.cache + "/hboMax/Cache", data)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}

func (f *flag_set) do_dash() error {
   data, err := os.ReadFile(f.cache + "/hboMax/Cache")
   if err != nil {
      return err
   }
   var cache hboMax.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return cache.Playback.PlayReady(data)
   }
   return f.config.Download(cache.MpdBody, cache.Mpd, f.dash)
}

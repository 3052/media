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

func main() {
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".mp4" {
         return ""
      }
      return "LP"
   })
   var program runner
   err := program.run()
   if err != nil {
      log.Fatal(err)
   }
}

func (r *runner) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   r.config.CertificateChain = cache + "/SL3000/CertificateChain"
   r.config.EncryptSignKey = cache + "/SL3000/EncryptSignKey"
   r.cache = cache + "/hboMax/Cache.json"

   flag.StringVar(&r.config.CertificateChain, "C", r.config.CertificateChain, "certificate chain")
   flag.StringVar(&r.config.EncryptSignKey, "E", r.config.EncryptSignKey, "encrypt sign key")
   flag.StringVar(&r.address, "a", "", "address")
   flag.StringVar(&r.dash, "d", "", "DASH ID")
   flag.StringVar(&r.edit, "e", "", "edit ID")
   flag.BoolVar(&r.initiate, "i", false, "device initiate")
   flag.BoolVar(&r.login, "l", false, "device login")
   flag.IntVar(&r.season, "s", 0, "season")
   flag.Parse()
   if r.initiate {
      return r.do_initiate()
   }
   if r.login {
      return r.do_login()
   }
   if r.address != "" {
      return r.do_address()
   }
   if r.edit != "" {
      return r.do_edit()
   }
   if r.dash != "" {
      return r.do_dash()
   }
   flag.Usage()
   return nil
}

func (r *runner) write(storage *cache) error {
   data, err := json.Marshal(storage)
   if err != nil {
      return err
   }
   log.Println("WriteFile", r.cache)
   return os.WriteFile(r.cache, data, os.ModePerm)
}

func (r *runner) do_initiate() error {
   var st hboMax.St
   err := st.Fetch()
   if err != nil {
      return err
   }
   err = r.write(&cache{st: &st})
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

func (r *runner) read(storage *cache) error {
   data, err := os.ReadFile(r.cache)
   if err != nil {
      return err
   }
   return json.Unmarshal(data, storage)
}

func (r *runner) do_login() error {
   var storage cache
   err := r.read(&storage)
   if err != nil {
      return err
   }
   storage.login, err = storage.st.Login()
   if err != nil {
      return err
   }
   return r.write(&storage)
}

func (r *runner) do_address() error {
   show_id, err := hboMax.ExtractId(r.address)
   if err != nil {
      return err
   }
   var storage cache
   err = r.read(&storage)
   if err != nil {
      return err
   }
   var videos *hboMax.Videos
   if r.season >= 1 {
      videos, err = storage.Login.Season(show_id, r.season)
   } else {
      videos, err = storage.Login.Movie(show_id)
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

type runner struct {
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

///

func (r *runner) do_edit() error {
   var storage cache
   err = r.read(&storage)
   storage.playback, err = storage.login.PlayReady(r.edit)
   if err != nil {
      return err
   }
   storage.mpd, storage.mpd_body, err = storage.playback.Dash()
   if err != nil {
      return err
   }
   err = r.write(&storage)
   if err != nil {
      return err
   }
   return net.Representations(storage.MpdBody, storage.Mpd)
}

type cache struct {
   login *Login
   mpd *url.URL
   mpd_body []byte
   playback *Playback
   st *hboMax.St
}

func (r *runner) do_dash() error {
   data, err := os.ReadFile(r.cache + "/hboMax/Cache")
   if err != nil {
      return err
   }
   var cache hboMax.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   r.config.Send = func(data []byte) ([]byte, error) {
      return cache.Playback.PlayReady(data)
   }
   return r.config.Download(cache.MpdBody, cache.Mpd, r.dash)
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

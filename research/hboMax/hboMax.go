package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hboMax"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func main() {
   err := new(program).run()
   if err != nil {
      log.Fatal(err)
   }
}

type program struct {
   cache_file string
   // 1
   proxy *string
   // 2
   initiate bool
   market   string
   // 3
   login bool
   // 4
   address string
   season  int
   // 5
   edit string
   // 6
   dash string
   job  maya.PlayReadyJob
}

type cache_data struct {
   Dash     *hboMax.Dash
   Login    *hboMax.Login
   Playback *hboMax.Playback
   Proxy    string
   St       *hboMax.St
}

///

func (p *program) run_proxy() error {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      if path.Ext(req.URL.Path) == ".mp4" {
         return "", false
      }
      return *p.proxy, true
   })
   return nil
}

func (p *program) run() error {
   cache_dir, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache_dir = filepath.ToSlash(cache_dir)
   p.cache_file = cache_dir + "/rosso/hboMax.xml"
   p.job.CertificateChain = cache_dir + "/SL3000/CertificateChain"
   p.job.EncryptSignKey = cache_dir + "/SL3000/EncryptSignKey"
   // 1
   flag.Func("x", "proxy", func(data string) error {
      p.proxy = &data
      return nil
   })
   // 2
   flag.BoolVar(&p.initiate, "i", false, "device initiate")
   flag.StringVar(
      &p.market, "m", hboMax.Markets[0], fmt.Sprint(hboMax.Markets),
   )
   // 3
   flag.BoolVar(&p.login, "l", false, "device login")
   // 4
   flag.StringVar(&p.address, "a", "", "address")
   flag.IntVar(&p.season, "s", 0, "season")
   // 5
   flag.StringVar(&p.edit, "e", "", "edit ID")
   // 6
   flag.StringVar(&p.dash, "d", "", "DASH ID")
   flag.StringVar(&p.job.CertificateChain, "C", p.job.CertificateChain, "certificate chain")
   flag.StringVar(&p.job.EncryptSignKey, "E", p.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   err = p.run_proxy()
   if err != nil {
      return err
   }
   if p.initiate {
      return p.run_initiate()
   }
   if p.login {
      return p.run_login()
   }
   if p.address != "" {
      return p.run_address()
   }
   if p.edit != "" {
      return p.run_edit()
   }
   if p.dash != "" {
      return p.run_dash()
   }
   return maya.Usage([][]string{
      {"x"},
      {"i", "m"},
      {"l"},
      {"a", "s"},
      {"e"},
      {"d", "C", "E"},
   })
}

func (p *program) run_address() error {
   var show hboMax.ShowKey
   err := show.Parse(p.address)
   if err != nil {
      return err
   }
   cache, err := maya.Read[cache_data](p.cache_file)
   if err != nil {
      return err
   }
   var videos *hboMax.Videos
   if p.season >= 1 {
      videos, err = cache.Login.Season(&show, p.season)
   } else {
      videos, err = cache.Login.Movie(&show)
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

func (p *program) run_login() error {
   cache, err := maya.Read[cache_data](p.cache_file)
   if err != nil {
      return err
   }
   cache.Login, err = cache.St.Login()
   if err != nil {
      return err
   }
   return maya.Write(p.cache_file, cache)
}

func (p *program) run_dash() error {
   cache, err := maya.Read[cache_data](p.cache_file)
   if err != nil {
      return err
   }
   p.job.Send = cache.Playback.PlayReady
   return p.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, p.dash)
}

func (p *program) run_edit() error {
   cache, err := maya.Read[cache_data](p.cache_file)
   if err != nil {
      return err
   }
   cache.Playback, err = cache.Login.PlayReady(p.edit)
   if err != nil {
      return err
   }
   cache.Dash, err = cache.Playback.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(p.cache_file, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

func (p *program) run_initiate() error {
   var st hboMax.St
   err := st.Fetch()
   if err != nil {
      return err
   }
   err = maya.Write(p.cache_file, &cache_data{St: &st})
   if err != nil {
      return err
   }
   initiate, err := st.Initiate(p.market)
   if err != nil {
      return err
   }
   fmt.Println(initiate)
   return nil
}

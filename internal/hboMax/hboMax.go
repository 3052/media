package main

import (
   "41.neocities.org/media/hboMax"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

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
   playback, err := login.PlayReady(f.edit)
   if err != nil {
      return err
   }
   data, address, err := playback.Dash()
   if err != nil {
      return err
   }
   if f.dash != "" {
      f.config.Send = func(data []byte) ([]byte, error) {
         return playback.PlayReady(data)
      }
      return f.config.Download(data, address, f.dash)
   }
   return f.config.PrintRepresentations(data, address)
}

func (f *flag_set) do_initiate() error {
   var st hboMax.St
   err := st.Token()
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
   log.Println("WriteFile", f.cache+"/hboMax/Login")
   return os.WriteFile(f.cache+"/hboMax/Login", data, os.ModePerm)
}

type flag_set struct {
   address  string
   cache    string
   config   net.Config
   edit     string
   initiate bool
   login    bool
   season   int
   dash     string
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
   show_id, err := hboMax.ExtractId(f.address)
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
   videos.FilterAndSort()
   for i, video := range videos.Included {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(video)
   }
   return nil
}

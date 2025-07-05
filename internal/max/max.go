package main

import (
   "41.neocities.org/media/max"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   http.DefaultTransport = &http.Transport{
      Proxy: func(req *http.Request) (*url.URL, error) {
         log.Println(req.Method, req.URL)
         return nil, nil
      },
   }
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
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.media + "/max/Login")
   if err != nil {
      return err
   }
   var login max.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   show_id, err := max.ShowId(f.address)
   if err != nil {
      return err
   }
   var videos *max.Videos
   if f.season >= 1 {
      videos, err = login.Season(show_id, f.season)
   } else {
      videos, err = login.Movie(show_id)
   }
   if err != nil {
      return err
   }
   for video := range videos.Seq() {
      fmt.Println(&video)
   }
   return nil
}
func (f *flag_set) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.cdm.ClientId = f.media + "/client_id.bin"
   f.cdm.PrivateKey = f.media + "/private_key.pem"
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.edit, "e", "", "edit ID")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.BoolVar(&f.initiate, "i", false, "device initiate")
   flag.BoolVar(&f.login, "l", false, "device login")
   flag.StringVar(&f.cdm.PrivateKey, "p", f.cdm.PrivateKey, "private key")
   flag.IntVar(&f.season, "s", 0, "season")
   flag.Parse()
   return nil
}

func (f *flag_set) do_edit() error {
   data, err := os.ReadFile(f.media + "/max/Login")
   if err != nil {
      return err
   }
   var login max.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = login.Widevine(f.edit)
   if err != nil {
      return err
   }
   var play max.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.Mpd())
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return play.License(data)
   }
   return f.filters.Filter(resp, &f.cdm)
}

type flag_set struct {
   address  string
   cdm        net.Cdm
   filters        net.Filters
   edit     string
   initiate bool
   login    bool
   media    string
   season   int
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) do_initiate() error {
   var st max.St
   err := st.New()
   if err != nil {
      return err
   }
   log.Println("Create", f.media+"/max/St")
   file, err := os.Create(f.media + "/max/St")
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
   data, err := os.ReadFile(f.media + "/max/St")
   if err != nil {
      return err
   }
   var st max.St
   err = st.Set(string(data))
   if err != nil {
      return err
   }
   data, err = st.Login()
   if err != nil {
      return err
   }
   return write_file(f.media+"/max/Login", data)
}

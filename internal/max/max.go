package main

import (
   "41.neocities.org/media/max"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
   "slices"
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
   return nil
}

type flags struct {
   address  string
   dash     string
   e        net.License
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

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "d", "", "DASH ID")
   flag.StringVar(&f.edit, "e", "", "edit ID")
   flag.BoolVar(&f.initiate, "i", false, "device initiate")
   flag.BoolVar(&f.login, "l", false, "device login")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.IntVar(&f.season, "s", 0, "season")
   flag.IntVar(&net.ThreadCount, "t", 1, "thread count")
   flag.Parse()
   switch {
   case f.initiate:
      err = f.do_initiate()
   case f.login:
      err = f.do_login()
   case f.address != "":
      err = f.do_address()
   case f.edit != "":
      err = f.do_edit()
   case f.dash != "":
      err = f.do_dash()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flags) do_initiate() error {
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

func (f *flags) do_login() error {
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

func (f *flags) do_address() error {
   data, err := os.ReadFile(f.media + "/max/Login")
   if err != nil {
      return err
   }
   var login max.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   var id max.ShowId
   err = id.Set(f.address)
   if err != nil {
      return err
   }
   var videos *max.Videos
   if f.season >= 1 {
      videos, err = login.Season(id, f.season)
   } else {
      videos, err = login.Movie(id)
   }
   if err != nil {
      return err
   }
   sorted := slices.SortedFunc(videos.Seq(), func(a, b max.Video) int {
      return a.Attributes.EpisodeNumber - b.Attributes.EpisodeNumber
   })
   for i, video := range sorted {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&video)
   }
   return nil
}

func (f *flags) do_edit() error {
   data, err := os.ReadFile(f.media + "/max/Login")
   if err != nil {
      return err
   }
   var login max.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = login.Playback(f.edit)
   if err != nil {
      return err
   }
   var play max.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.media+"/max/Playback", data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.Mpd())
   if err != nil {
      return err
   }
   return net.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
   data, err := os.ReadFile(f.media + "/max/Playback")
   if err != nil {
      return err
   }
   var play max.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   f.e.Widevine = func(data []byte) ([]byte, error) {
      return play.Widevine(data)
   }
   return f.e.Download(f.media+"/Mpd", f.dash)
}

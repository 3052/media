package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/max"
   "fmt"
   "log"
   "os"
)

func create(name string) (*os.File, error) {
   log.Println("Create", name)
   return os.Create(name)
}

func (f *flags) do_initiate() error {
   var st max.St
   err := st.New()
   if err != nil {
      return err
   }
   file, err := create("st.txt")
   if err != nil {
      return err
   }
   defer file.Close()
   fmt.Fprint(file, st)
   initiate, err := st.Initiate()
   if err != nil {
      return err
   }
   fmt.Printf("%+v\n", initiate)
   return nil
}

func (f *flags) download() error {
   data, err := os.ReadFile(f.home + "/max.txt")
   if err != nil {
      return err
   }
   var login max.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   play, err := login.Playback(&f.url)
   if err != nil {
      return err
   }
   represents, err := internal.Mpd(play)
   if err != nil {
      return err
   }
   for _, represent := range represents {
      switch f.representation {
      case "":
         fmt.Print(&represent, "\n\n")
      case represent.Id:
         f.s.Client = play
         return f.s.Download(&represent)
      }
   }
   return nil
}

func (f *flags) do_login() error {
   data, err := os.ReadFile("st.txt")
   if err != nil {
      return err
   }
   var st max.St
   err = st.Set(string(data))
   if err != nil {
      return err
   }
   data, err = max.Login{}.Marshal(st)
   if err != nil {
      return err
   }
   return os.WriteFile(f.home+"/max.txt", data, os.ModePerm)
}

package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/plex"
   "errors"
   "flag"
   "fmt"
   "log"
   "os"
   "path/filepath"
)

func (f *flags) write_file(name string, data []byte) error {
   log.Println("WriteFile", f.media + name)
   return os.WriteFile(f.media + name, data, os.ModePerm)
}

func (f *flags) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/plex/User")
      if err != nil {
         return err
      }
      var user plex.User
      err = user.Unmarshal(data)
      if err != nil {
         return err
      }
      data, err = os.ReadFile(f.media + "/plex/Metadata")
      if err != nil {
         return err
      }
      var metadata plex.Metadata
      err = metadata.Unmarshal(data)
      if err != nil {
         return err
      }
      part, _ := metadata.Dash()
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return user.Widevine(part, data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   data, err := plex.NewUser()
   if err != nil {
      return err
   }
   err = f.write_file("/plex/User", data)
   if err != nil {
      return err
   }
   var user plex.User
   err = user.Unmarshal(data)
   if err != nil {
      return err
   }
   match, err := user.Match(f.address)
   if err != nil {
      return err
   }
   data1, err := user.Metadata(match)
   if err != nil {
      return err
   }
   err = f.write_file("/plex/Metadata", data1)
   if err != nil {
      return err
   }
   var metadata plex.Metadata
   err = metadata.Unmarshal(data1)
   if err != nil {
      return err
   }
   part, ok := metadata.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := user.Mpd(part)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}

type flags struct {
   address        plex.Address
   e              internal.License
   get_forward    bool
   dash string
   media          string
}

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

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.Var(&f.address, "a", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.BoolVar(&f.get_forward, "g", false, "get forward")
   flag.StringVar(&plex.ForwardedFor, "e", "", "set forward")
   flag.Parse()
   switch {
   case f.get_forward:
      for _, forward := range internal.Forward {
         fmt.Println(forward.Country, forward.Ip)
      }
   case f.address[0] != "":
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}

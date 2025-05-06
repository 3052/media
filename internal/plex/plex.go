package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/plex"
   "errors"
   "flag"
   "log"
   "os"
   "path/filepath"
)

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
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

type flags struct {
   media   string
   e       internal.License
   address string
   dash    string
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "d", "", "dash ID")
   flag.StringVar(&plex.ForwardedFor, "f", "", "set forward")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.Parse()
   switch {
   case f.address != "":
      err = f.do_address()
   case f.dash != "":
      err = f.do_dash()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flags) do_address() error {
   data, err := plex.NewUser()
   if err != nil {
      return err
   }
   err = write_file(f.media + "/plex/User", data)
   if err != nil {
      return err
   }
   var user plex.User
   err = user.Unmarshal(data)
   if err != nil {
      return err
   }
   var url plex.Url
   url.New(f.address)
   match, err := user.Match(url)
   if err != nil {
      return err
   }
   data1, err := user.Metadata(match)
   if err != nil {
      return err
   }
   err = write_file(f.media + "/plex/Metadata", data1)
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

func (f *flags) do_dash() error {
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

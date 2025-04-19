package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/tubi"
   "41.neocities.org/platform/mullvad"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
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

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.IntVar(&f.tubi, "b", 0, "Tubi ID")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.BoolVar(&f.mullvad, "m", false, "Mullvad")
   flag.Parse()
   if f.mullvad {
      http.DefaultClient.Transport = &mullvad.Transport{}
      defer mullvad.Disconnect()
   }
   switch {
   case f.tubi >= 1:
      err := f.do_tubi()
      if err != nil {
         panic(err)
      }
   case f.dash != "":
      err := f.do_dash()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}

func (f *flags) do_tubi() error {
   data, err := tubi.NewContent(f.tubi)
   if err != nil {
      return err
   }
   log.Println("WriteFile", f.media+"/tubi/Content")
   err = os.WriteFile(f.media+"/tubi/Content", data, os.ModePerm)
   if err != nil {
      return err
   }
   var content tubi.Content
   err = content.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(content.VideoResources[0].Manifest.Url)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}

type flags struct {
   e       internal.License
   media   string
   mullvad bool
   
   tubi    int
   dash    string
}

func (f *flags) do_dash() error {
   data, err := os.ReadFile(f.media + "/tubi/Content")
   if err != nil {
      return err
   }
   var content tubi.Content
   err = content.Unmarshal(data)
   if err != nil {
      return err
   }
   f.e.Widevine = func(data []byte) ([]byte, error) {
      return content.VideoResources[0].Widevine(data)
   }
   return f.e.Download(f.media+"/Mpd", f.dash)
}

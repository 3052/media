package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/roku"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/roku/Playback")
      if err != nil {
         return err
      }
      var play roku.Playback
      err = play.Unmarshal(data)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return play.Widevine(data)
      }
      return f.e.Download(f.media + "/Mpd", f.dash)
   }
   var code *roku.Code
   if f.token_read {
      data, err := os.ReadFile(f.media + "/roku/Code")
      if err != nil {
         return err
      }
      code = &roku.Code{}
      err = code.Unmarshal(data)
      if err != nil {
         return err
      }
   }
   data, err := code.AccountToken()
   if err != nil {
      return err
   }
   var token roku.AccountToken
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   data1, err := token.Playback(f.roku)
   if err != nil {
      return err
   }
   err = f.write_file("/roku/Playback", data1)
   if err != nil {
      return err
   }
   var play roku.Playback
   err = play.Unmarshal(data1)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.Url)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media + "/Mpd", resp)
}

type flags struct {
   code_write     bool
   e              internal.License
   media          string
   dash string
   roku           string
   token_read     bool
   token_write    bool
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
   flag.StringVar(&f.roku, "b", "", "Roku ID")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.BoolVar(&f.code_write, "code", false, "write code")
   flag.BoolVar(&f.token_write, "token", false, "write token")
   flag.BoolVar(&f.token_read, "t", false, "read token")
   flag.Parse()
   switch {
   case f.code_write:
      err := f.write_code()
      if err != nil {
         panic(err)
      }
   case f.token_write:
      err := f.write_token()
      if err != nil {
         panic(err)
      }
   case f.roku != "":
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}

func (f *flags) write_file(name string, data []byte) error {
   log.Println("WriteFile", f.media + name)
   return os.WriteFile(f.media + name, data, os.ModePerm)
}

func (f *flags) write_code() error {
   data, err := (*roku.Code).AccountToken(nil)
   if err != nil {
      return err
   }
   var token roku.AccountToken
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   err = f.write_file("/roku/AccountToken", data)
   if err != nil {
      return err
   }
   data1, err := token.Activation()
   if err != nil {
      return err
   }
   var activation roku.Activation
   err = activation.Unmarshal(data1)
   if err != nil {
      return err
   }
   fmt.Println(&activation)
   return f.write_file("/roku/Activation", data1)
}

func (f *flags) write_token() error {
   data, err := os.ReadFile(f.media + "/roku/AccountToken")
   if err != nil {
      return err
   }
   var token roku.AccountToken
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = os.ReadFile(f.media + "/roku/Activation")
   if err != nil {
      return err
   }
   var activation roku.Activation
   err = activation.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = token.Code(&activation)
   if err != nil {
      return err
   }
   return f.write_file("/roku/Code", data)
}

package main

import (
   "41.neocities.org/media/canal"
   "41.neocities.org/media/internal"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/canal/Token")
      if err != nil {
         return err
      }
      var token canal.Token
      err = token.Unmarshal(data)
      if err != nil {
         return err
      }
      session, err := token.Session()
      if err != nil {
         return err
      }
      var fields canal.Fields
      err = fields.New(f.address)
      if err != nil {
         return err
      }
      play, err := session.Play(fields.ObjectIds())
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return play.Widevine(data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/canal/Token")
   if err != nil {
      return err
   }
   var token canal.Token
   err = token.Unmarshal(data)
   if err != nil {
      return err
   }
   session, err := token.Session()
   if err != nil {
      return err
   }
   var fields canal.Fields
   err = fields.New(f.address)
   if err != nil {
      return err
   }
   play, err := session.Play(fields.ObjectIds())
   if err != nil {
      return err
   }
   resp, err := http.Get(play.Url)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
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
   e        internal.License
   email    string
   media    string
   password string
   dash     string
   address  string
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "p", "", "password")
   flag.StringVar(&f.address, "a", "", "address")
   flag.Parse()
   switch {
   case f.password != "":
      err := f.authenticate()
      if err != nil {
         panic(err)
      }
   case f.address != "":
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}

func (f *flags) write_file(name string, data []byte) error {
   log.Println("WriteFile", f.media+name)
   return os.WriteFile(f.media+name, data, os.ModePerm)
}

func (f *flags) authenticate() error {
   var ticket canal.Ticket
   err := ticket.New()
   if err != nil {
      return err
   }
   data, err := ticket.Token(f.email, f.password)
   if err != nil {
      return err
   }
   return f.write_file("/canal/Token", data)
}

package main

import (
   "41.neocities.org/media/canal"
   "41.neocities.org/media/internal"
   "41.neocities.org/platform/proxy"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

type flags struct {
   dash      string
   e         internal.License
   email     string
   media     string
   password  string
   proxy     bool
   object_id canal.ObjectId
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.Var(&f.object_id, "a", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "dash ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "password", "", "password")
   flag.StringVar(&f.email, "email", "", "email")
   flag.BoolVar(&f.proxy, "p", false, "proxy")
   flag.Parse()
   if f.proxy {
      http.DefaultClient.Transport = &proxy.Transport{
         Protocols: &http.Protocols{}, // github.com/golang/go/issues/25793
         Proxy:     http.ProxyFromEnvironment,
      }
   }
   switch {
   case f.password != "":
      err := f.authenticate()
      if err != nil {
         panic(err)
      }
   case f.object_id[0] != "":
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
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
   token, err := ticket.Token(f.email, f.password)
   if err != nil {
      return err
   }
   data, err := canal.NewSession(token.SsoToken)
   if err != nil {
      return err
   }
   return f.write_file("/canal/Session", data)
}

func (f *flags) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/canal/Play")
      if err != nil {
         return err
      }
      var play canal.Play
      err = play.Unmarshal(data)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return play.Widevine(data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   data, err := os.ReadFile(f.media + "/canal/Session")
   if err != nil {
      return err
   }
   var session canal.Session
   err = session.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = canal.NewSession(session.SsoToken)
   if err != nil {
      return err
   }
   err = session.Unmarshal(data)
   if err != nil {
      return err
   }
   err = f.write_file("/canal/Session", data)
   if err != nil {
      return err
   }
   data, err = session.Play(f.object_id)
   if err != nil {
      return err
   }
   var play canal.Play
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   err = f.write_file("/canal/Play", data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.Url)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}

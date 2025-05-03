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

func (f *flags) do_dash() error {
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

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flags) do_email() error {
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
   return write_file(f.media+"/canal/Session", data)
}

type flags struct {
   address  string
   dash     string
   e        internal.License
   email    string
   media    string
   password string
   proxy    bool
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.address, "a", "", "canal address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "d", "", "dash ID")
   flag.StringVar(&f.email, "email", "", "canal email")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "password", "", "canal password")
   flag.BoolVar(&f.proxy, "proxy", false, "proxy server")
   flag.Parse()
   if f.proxy {
      http.DefaultClient.Transport = &proxy.Transport{
         Protocols: &http.Protocols{}, // github.com/golang/go/issues/25793
         Proxy:     http.ProxyFromEnvironment,
      }
   }
   if f.email != "" {
      if f.password != "" {
         err := f.do_email()
         if err != nil {
            panic(err)
         }
      }
   } else if f.address != "" {
      err := f.do_address()
      if err != nil {
         panic(err)
      }
   } else if f.dash != "" {
      err := f.do_dash()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

func (f *flags) do_address() error {
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
   err = write_file(f.media+"/canal/Session", data)
   if err != nil {
      return err
   }
   var fields canal.Fields
   err = fields.New(f.address)
   if err != nil {
      return err
   }
   data, err = session.Play(fields)
   if err != nil {
      return err
   }
   var play canal.Play
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.media+"/canal/Play", data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.Url)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}

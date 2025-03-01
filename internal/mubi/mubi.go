package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/mubi"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func main() {
   // github.com/golang/go/issues/18639
   // we dont need this until later, but you have to call before the first
   // request in the program
   os.Setenv("GODEBUG", "http2client=0")
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.Var(&f.address, "a", "address")
   flag.BoolVar(&f.auth, "auth", false, "authenticate")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.BoolVar(&f.code, "code", false, "link code")
   flag.StringVar(&f.representation, "i", "", "representation")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.BoolVar(&f.text, "text", false, "text track")
   flag.Parse()
   switch {
   case f.code:
      err := f.do_code()
      if err != nil {
         panic(err)
      }
   case f.auth:
      err := f.do_auth()
      if err != nil {
         panic(err)
      }
   case f.address[0] != "":
      err := f.do_dash()
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

func (f *flags) do_code() error {
   data, err := mubi.NewLinkCode()
   if err != nil {
      return err
   }
   var code mubi.LinkCode
   err = code.Unmarshal(data)
   if err != nil {
      return err
   }
   fmt.Println(&code)
   return f.write_file("/mubi/LinkCode", data)
}

func (f *flags) do_auth() error {
   data, err := os.ReadFile(f.media + "/mubi/LinkCode")
   if err != nil {
      return err
   }
   var code mubi.LinkCode
   err = code.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = code.Authenticate()
   if err != nil {
      return err
   }
   return f.write_file("/mubi/Authenticate", data)
}

func (f *flags) do_dash() error {
   if f.representation != "" {
      data, err := os.ReadFile(f.media + "/mubi/Authenticate")
      if err != nil {
         return err
      }
      var auth mubi.Authenticate
      err = auth.Unmarshal(data)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return auth.Widevine(data)
      }
      return f.e.Download(f.media + "/Mpd", f.representation)
   }
   data, err := os.ReadFile(f.media + "/mubi/Authenticate")
   if err != nil {
      return err
   }
   var auth mubi.Authenticate
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   film, err := f.address.Film()
   if err != nil {
      return err
   }
   err = auth.Viewing(film)
   if err != nil {
      return err
   }
   secure, err := auth.SecureUrl(film)
   if err != nil {
      return err
   }
   if f.text {
      for _, text := range secure.TextTrackUrls {
         err = func() error {
            resp, err := http.Get(text.Url)
            if err != nil {
               return err
            }
            defer resp.Body.Close()
            file, err := os.Create(text.Base())
            if err != nil {
               return err
            }
            defer file.Close()
            _, err = file.ReadFrom(resp.Body)
            if err != nil {
               return err
            }
            return nil
         }()
         if err != nil {
            return err
         }
      }
      return nil
   }
   resp, err := http.Get(secure.Url)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media + "/Mpd", resp)
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
   address        mubi.Address
   auth           bool
   code           bool
   e              internal.License
   media          string
   representation string
   text           bool
}

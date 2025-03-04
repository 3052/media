package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/rakuten"
   "errors"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

type flags struct {
   address  rakuten.Address
   dash     string
   e        internal.License
   language string
   media    string
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
   flag.StringVar(&f.language, "b", "", "language")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.Parse()
   if f.address.MarketCode != "" {
      if f.language != "" {
         err := f.download()
         if err != nil {
            panic(err)
         }
      } else {
         err := f.do_language()
         if err != nil {
            panic(err)
         }
      }
   } else {
      flag.Usage()
   }
}

func (f *flags) write_file(name string, data []byte) error {
   log.Println("WriteFile", f.media + name)
   return os.WriteFile(f.media + name, data, os.ModePerm)
}

func (f *flags) do_language() error {
   class, ok := f.address.ClassificationId()
   if !ok {
      return errors.New(".ClassificationId()")
   }
   var content *rakuten.Content
   if f.address.SeasonId != "" {
      data, err := f.address.Season(class)
      if err != nil {
         return err
      }
      var season rakuten.Season
      err = season.Unmarshal(data)
      if err != nil {
         return err
      }
      err = f.write_file("/rakuten/Season", data)
      if err != nil {
         return err
      }
      content, ok = season.Content(&f.address)
      if !ok {
         return errors.New(".Content")
      }
   } else {
      data, err := f.address.Movie(class)
      if err != nil {
         return err
      }
      content = &rakuten.Content{}
      err = content.Unmarshal(data)
      if err != nil {
         return err
      }
      err = f.write_file("/rakuten/Content", data)
      if err != nil {
         return err
      }
   }
   fmt.Println(content)
   return nil
}

func (f *flags) download() error {
   class, ok := f.address.ClassificationId()
   if !ok {
      return errors.New(".ClassificationId()")
   }
   var content *rakuten.Content
   if f.address.SeasonId != "" {
      data, err := os.ReadFile(f.media + "/rakuten/Season")
      if err != nil {
         return err
      }
      var season rakuten.Season
      err = season.Unmarshal(data)
      if err != nil {
         return err
      }
      content, ok = season.Content(&f.address)
      if !ok {
         return errors.New(".Content")
      }
   } else {
      data, err := os.ReadFile(f.media + "/rakuten/Content")
      if err != nil {
         return err
      }
      content = &rakuten.Content{}
      err = content.Unmarshal(data)
      if err != nil {
         return err
      }
   }
   stream := content.Streamings()
   if f.dash != "" {
      stream.Hd()
      info, err := stream.Info(f.language, class)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return info.Widevine(data)
      }
      return f.e.Download(f.media + "/Mpd", f.dash)
   }
   stream.Fhd()
   info, err := stream.Info(f.language, class)
   if err != nil {
      return err
   }
   resp, err := http.Get(info.Url)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media + "/Mpd", resp)
}

package main

import (
   "41.neocities.org/media/rakuten"
   "41.neocities.org/stream"
   "errors"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

type flags struct {
   address string
   dash string
   e     stream.License
   language string
   media string
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.language, "b", "", "language")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.StringVar(&f.e.PrivateKey, "k", f.e.PrivateKey, "private key")
   flag.IntVar(&stream.ThreadCount, "t", 1, "thread count")
   flag.Parse()
   if f.address != "" {
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

func (f *flags) do_language() error {
   var address rakuten.Address
   address.Set(f.address)
   class, ok := address.ClassificationId()
   if !ok {
      return errors.New(".ClassificationId()")
   }
   var content *rakuten.Content
   if address.SeasonId != "" {
      data, err := address.Season(class)
      if err != nil {
         return err
      }
      var season rakuten.Season
      err = season.Unmarshal(data)
      if err != nil {
         return err
      }
      err = write_file(f.media+"/rakuten/Season", data)
      if err != nil {
         return err
      }
      content, ok = season.Content(&address)
      if !ok {
         return errors.New(".Content")
      }
   } else {
      data, err := address.Movie(class)
      if err != nil {
         return err
      }
      content = &rakuten.Content{}
      err = content.Unmarshal(data)
      if err != nil {
         return err
      }
      err = write_file(f.media+"/rakuten/Content", data)
      if err != nil {
         return err
      }
   }
   fmt.Println(content)
   return nil
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

func (f *flags) download() error {
   var address rakuten.Address
   address.Set(f.address)
   class, ok := address.ClassificationId()
   if !ok {
      return errors.New(".ClassificationId()")
   }
   var content *rakuten.Content
   if address.SeasonId != "" {
      data, err := os.ReadFile(f.media + "/rakuten/Season")
      if err != nil {
         return err
      }
      var season rakuten.Season
      err = season.Unmarshal(data)
      if err != nil {
         return err
      }
      content, ok = season.Content(&address)
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
   streaming := content.Streamings()
   if f.dash != "" {
      streaming.Hd()
      info, err := streaming.Info(f.language, class)
      if err != nil {
         return err
      }
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return info.Widevine(data)
      }
      return f.e.Download(f.media+"/Mpd", f.dash)
   }
   streaming.Fhd()
   info, err := streaming.Info(f.language, class)
   if err != nil {
      return err
   }
   resp, err := http.Get(info.Url)
   if err != nil {
      return err
   }
   return stream.Mpd(f.media+"/Mpd", resp)
}

package main

import (
   "41.neocities.org/media/rakuten"
   "41.neocities.org/net"
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

type flag_set struct {
   cdm      net.Cdm
   filters  net.Filters
   media    string
   address  string
   language string
}

/*
print movie:

rakuten -movie url

print seasons:

rakuten -show url

print episodes:

rakuten -season id

download:

rakuten -c content_id -a audio_language
*/

///

func (f *flag_set) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.cdm.ClientId = f.media + "/client_id.bin"
   f.cdm.PrivateKey = f.media + "/private_key.pem"
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.language, "b", "", "language")
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.cdm.PrivateKey, "k", f.cdm.PrivateKey, "private key")
   flag.IntVar(&net.Threads, "t", 2, "threads")
   flag.Parse()
   return nil
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.address != "" {
      if set.language != "" {
         // second
         // do we really need address for this?
         err = set.do_language()
      } else {
         // first
         err = set.do_address()
      }
   } else {
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flag_set) do_address() error {
   var path rakuten.Path
   path.New(f.address)
   class, ok := path.ClassificationId()
   if !ok {
      return errors.New(".ClassificationId()")
   }
   var content *rakuten.Content
   if path.SeasonId != "" {
      data, err := path.Season(class)
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
      content, ok = season.Content(&path)
      if !ok {
         return errors.New(".Content")
      }
   } else {
      data, err := path.Movie(class)
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

func (f *flag_set) do_language() error {
   var path rakuten.Path
   path.New(f.address)
   class, ok := path.ClassificationId()
   if !ok {
      return errors.New(".ClassificationId()")
   }
   var content *rakuten.Content
   if path.SeasonId != "" {
      data, err := os.ReadFile(f.media + "/rakuten/Season")
      if err != nil {
         return err
      }
      var season rakuten.Season
      err = season.Unmarshal(data)
      if err != nil {
         return err
      }
      content, ok = season.Content(&path)
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
   streaming.Fhd()
   info, err := streaming.Info(f.language, class)
   if err != nil {
      return err
   }
   resp, err := http.Get(info.Url)
   if err != nil {
      return err
   }
   streaming.Hd()
   info, err = streaming.Info(f.language, class)
   if err != nil {
      return err
   }
   f.cdm.License = func(data []byte) ([]byte, error) {
      return info.License(data)
   }
   return f.filters.Filter(resp, &f.cdm)
}

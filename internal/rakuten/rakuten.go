package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/rakuten"
   "errors"
   "flag"
   "fmt"
   "log"
   "os"
   "path/filepath"
)

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.Var(&f.address, "a", "address")
   flag.StringVar(&f.language, "b", "", "language")
   flag.StringVar(&f.s.ClientId, "c", f.s.ClientId, "client ID")
   flag.StringVar(&f.representation, "i", "", "representation")
   flag.StringVar(&f.s.PrivateKey, "k", f.s.PrivateKey, "private key")
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

func (f *flags) New() error {
   home, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   home = filepath.ToSlash(home)
   f.s.ClientId = home + "/widevine/client_id.bin"
   f.s.PrivateKey = home + "/widevine/private_key.pem"
   return nil
}

type flags struct {
   address        rakuten.Address
   representation string
   s              internal.Stream
   language       string
}
func (f *flags) download() error {
   class, ok := f.address.ClassificationId()
   if !ok {
      return errors.New("Address.ClassificationId")
   }
   var content *rakuten.Content
   if f.address.SeasonId != "" {
      season, err := f.address.Season(class)
      if err != nil {
         return err
      }
      content, ok = season.Content(&f.address)
      if !ok {
         return errors.New("Season.Content")
      }
   } else {
      var err error
      content, err = f.address.Movie(class)
      if err != nil {
         return err
      }
   }
   stream := content.Streamings()
   stream.Fhd()
   info, err := stream.Info(f.language, class)
   if err != nil {
      return err
   }
   represents, err := internal.Mpd(info)
   if err != nil {
      return err
   }
   for _, represent := range represents {
      switch f.representation {
      case "":
         fmt.Print(&represent, "\n\n")
      case represent.Id:
         stream.Hd()
         info, err = stream.Info(f.language, class)
         if err != nil {
            return err
         }
         f.s.Client = info
         return f.s.Download(&represent)
      }
   }
   return nil
}

func (f *flags) do_language() error {
   class, ok := f.address.ClassificationId()
   if !ok {
      return errors.New("Address.ClassificationId")
   }
   var content *rakuten.Content
   if f.address.SeasonId != "" {
      season, err := f.address.Season(class)
      if err != nil {
         return err
      }
      content, ok = season.Content(&f.address)
      if !ok {
         return errors.New("Season.Content")
      }
   } else {
      var err error
      content, err = f.address.Movie(class)
      if err != nil {
         return err
      }
   }
   fmt.Println(content)
   return nil
}

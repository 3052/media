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

func (f *flag_set) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.cdm.ClientId = f.media + "/client_id.bin"
   f.cdm.PrivateKey = f.media + "/private_key.pem"
   flag.StringVar(&f.cdm.ClientId, "C", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.cdm.PrivateKey, "P", f.cdm.PrivateKey, "private key")
   ///////////////////////////////////////////////////////////////////////
   flag.StringVar(&f.movie, "movie", "", "print movie")
   ////////////////////////////////////////////////////
   flag.StringVar(&f.show, "show", "", "print seasons")
   ////////////////////////////////////////////////////
   flag.StringVar(&f.season, "season", "", "print episodes")
   /////////////////////////////////////////////////////////
   flag.StringVar(&f.content, "c", "", "content ID")
   flag.StringVar(&f.language, "a", "", "audio language")
   flag.Parse()
   return nil
}

type flag_set struct {
   cdm      net.Cdm
   filters  net.Filters
   media    string
   ///////////////
   movie string
   ////////////
   show string
   /////////////
   season string
   /////////////
   content string
   language string
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.movie != "" {
      err = set.do_movie()
   } else if set.show != "" {
      err = set.do_show()
   } else if set.season != "" {
      err = set.do_season()
   } else if set.content != "" {
      if set.language != "" {
         err = set.do_content()
      }
   } else {
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

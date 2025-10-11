package main

import (
   "41.neocities.org/media/tubi"
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.e.ClientId = f.cache + "/L3/client_id.bin"
   f.e.PrivateKey = f.cache + "/L3/private_key.pem"
   f.bitrate.Value = [][2]int{{100_000, 200_000}, {2_200_000, 4_000_000}}
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.IntVar(&net.Threads, "threads", 1, "threads")
   flag.IntVar(&f.tubi, "t", 0, "Tubi ID")
   flag.Var(&f.bitrate, "b", "bitrate")
   flag.Parse()
   return nil
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.tubi >= 1 {
      err = set.do_tubi()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

type flag_set struct {
   e     net.License
   cache string
   tubi    int
   bitrate net.Bitrate
}

func (f *flag_set) do_tubi() error {
   data, err := tubi.NewContent(f.tubi)
   if err != nil {
      return err
   }
   var content tubi.Content
   err = content.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(content.VideoResources[0].Manifest.Url)
   if err != nil {
      return err
   }
   f.e.Widevine = func(data []byte) ([]byte, error) {
      return content.VideoResources[0].Widevine(data)
   }
   return f.e.Bitrate(resp, &f.bitrate)
}

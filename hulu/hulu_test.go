package hulu

import (
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
   "testing"
)

func (f *flag_set) do_address() error {
   data, err := os.ReadFile(f.media + "/hulu/Authenticate")
   if err != nil {
      return err
   }
   var auth hulu.Authenticate
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   err = auth.Refresh()
   if err != nil {
      return err
   }
   deep, err := auth.DeepLink(hulu.Id(f.address))
   if err != nil {
      return err
   }
   data, err = auth.Playlist(deep)
   if err != nil {
      return err
   }
   err = write_file(f.media+"/hulu/Playlist", data)
   if err != nil {
      return err
   }
   var play hulu.Playlist
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(play.StreamUrl)
   if err != nil {
      return err
   }
   return net.Mpd(f.media+"/Mpd", resp)
}

var tests = []struct {
   content string
   url     string
}{
   {
      content: "film",
      url:     "hulu.com/watch/f70dfd4d-dbfb-46b8-abb3-136c841bba11",
   },
   {
      content: "episode",
      url:     "hulu.com/watch/023c49bf-6a99-4c67-851c-4c9e7609cc1d",
   },
}

func TestDeepLink(t *testing.T) {
   fmt.Println(tests)
}

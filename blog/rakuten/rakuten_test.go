package rakuten

import (
   "41.neocities.org/net"
   "fmt"
   "net/http"
   "testing"
)

const test_address = "//rakuten.tv/uk?content_type=tv_shows&tv_show_id=hell-s-kitchen-usa"

func Test(t *testing.T) {
   http.DefaultTransport = net.Transport(nil)
   var web address
   err := web.Set(test_address)
   if err != nil {
      t.Fatal(err)
   }
   seasons, err := web.seasons()
   if err != nil {
      t.Fatal(err)
   }
   for _, season1 := range seasons {
      fmt.Println(season1.Id)
   }
}

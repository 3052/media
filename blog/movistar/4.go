package movistar

import (
   "encoding/json"
   "net/http"
   "strconv"
)

type vod_item struct {
   CasId string
   UrlVideo string // MPD mullvad
}

func new_vod_item(id int64) (*vod_item, error) {
   req, _ := http.NewRequest("", "https://ottcache.dof6.com", nil)
   req.URL.Path = func() string {
      b := []byte("/movistarplus/amazon.tv/contents/")
      b = strconv.AppendInt(b, id, 10)
      b = append(b, "/details"...)
      return string(b)
   }()
   req.URL.RawQuery = "mdrm=true"
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      VodItems []vod_item
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.VodItems[0], nil
}

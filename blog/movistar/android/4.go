package movistar

import (
   "encoding/json"
   "net/http"
   "strconv"
)

type details struct {
   VodItems []struct {
      // mullvad fail
      // nord pass
      UrlVideo string
   }
}

func (d *details) New(id int64) error {
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
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(d)
}

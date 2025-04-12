package movistar

import (
   "net/http"
   "strconv"
)

func details(id int64) (*http.Response, error) {
   req, _ := http.NewRequest("", "https://ottcache.dof6.com", nil)
   req.URL.Path = func() string {
      b := []byte("/movistarplus/amazon.tv/contents/")
      b = strconv.AppendInt(b, id, 10)
      b = append(b, "/details"...)
      return string(b)
   }()
   return http.DefaultClient.Do(req)
}

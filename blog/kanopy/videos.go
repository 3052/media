package kanopy

import (
   "encoding/json"
   "net/http"
   "strconv"
)

type videos_response struct {
   Video *struct {
      ProductionYear int
      Title string
   }
}

func (w *web_token) videos(id int64) (*videos_response, error) {
   req, err := http.NewRequest("", "https://www.kanopy.com", nil)
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/kapi/videos/" + strconv.FormatInt(id, 10)
   req.Header = http.Header{
      "authorization": {"Bearer " + w.Jwt},
      "user-agent": {user_agent},
      "x-version": {x_version},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   videos := &videos_response{}
   err = json.NewDecoder(resp.Body).Decode(videos)
   if err != nil {
      return nil, err
   }
   return videos, nil
}

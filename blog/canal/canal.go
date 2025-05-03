package canal

import (
   "encoding/json"
   "errors"
   "net/http"
   "strconv"
)

const bearer = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0di5zb2xvY29vLmF1dGgiOnsicyI6IiEhISEiLCJ1IjoicE5sZkdGQzFqa0o0dDhsT3h5Q0s0ZyIsImwiOiJlbl9VUyIsImQiOiJQQyIsIm9tIjoiTyIsImMiOiJ2aDhUOGhEM29ZSExNaVVBQkNZMk82RTcyUDlGMHFwQ3lCR2tXc3VGSjVjIiwic3QiOiJmdWxsIiwiZyI6ImV5SnZjQ0k2SWpFd01EUTJJaXdpWkdJaU9tWmhiSE5sTENKd2RDSTZabUZzYzJVc0luVndJam9pYlRkamNDSXNJbUp5SWpvaWJUZGpjQ0lzSW1SbElqb2lZbkpoYm1STllYQndhVzVuSW4wIiwiYiI6Im03Y3AifSwibmJmIjoxNzQ2MzA1NTUyLCJleHAiOjE3NDYzMDcxODksImlhdCI6MTc0NjMwNTU1MiwiYXVkIjoibTdjcCJ9.Wj5NxqVx21vu8WYjMnMwlMJ2EuHiu7POCrm_GifwfCE"

type asset struct {
   Params struct {
      SeriesSeason  string
      SeriesEpisode int
   }
   Title string
   Id    string
}

func assets(series_id string, season int64) ([]asset, error) {
   req, _ := http.NewRequest("", "https://tvapi-hlm2.solocoo.tv/v1/assets", nil)
   req.Header.Set("authorization", "Bearer "+bearer)
   req.URL.RawQuery = func() string {
      b := []byte("limit=99&query=episodes,")
      b = append(b, series_id...)
      b = append(b, ",season,"...)
      b = strconv.AppendInt(b, season, 10)
      return string(b)
   }()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Assets  []asset
      Message string
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if value.Message != "" {
      return nil, errors.New(value.Message)
   }
   return value.Assets, nil
}

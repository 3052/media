package canal

import "net/http"

const bearer = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0di5zb2xvY29vLmF1dGgiOnsicyI6IiEhISEiLCJ1IjoicE5sZkdGQzFqa0o0dDhsT3h5Q0s0ZyIsImwiOiJjc19DWiIsImQiOiJQQyIsIm9tIjoiTyIsImMiOiJyUVZXbEdZdHQ4ZjctWkoyeGRpR3A5YVRrQkk1R1Y3SkVvN3p6S2pzc2Q0Iiwic3QiOiJmdWxsIiwiZyI6ImV5SmtaU0k2SW1KeVlXNWtUV0Z3Y0dsdVp5SXNJbTl3SWpvaU1UQXdORFlpTENKaWNpSTZJbTAzWTNBaUxDSjFjQ0k2SW0wM1kzQWlMQ0prWWlJNlptRnNjMlVzSW5CMElqcG1ZV3h6WlgwIiwiYiI6Im03Y3AifSwibmJmIjoxNzQ2MzAwMjUwLCJleHAiOjE3NDYzMDIxNjMsImlhdCI6MTc0NjMwMDI1MCwiYXVkIjoibTdjcCJ9.-NZcgD43XBgAQY_XEA4V2i6ZBsBLRQXzTuNu94E4WqI"

func assets() (*http.Response, error) {
   req, _ := http.NewRequest("", "https://tvapi-hlm2.solocoo.tv/v1/assets", nil)
   req.Header.Set("authorization", "Bearer " + bearer)
   req.URL.RawQuery = "query=episodes,ZXkaWHVpx827Fz_4ZNtW5l8MoKD5_2lhv0nYe4m3,season,2"
   //value["sort"] = []string{"seasonepisode"}
   return http.DefaultClient.Do(req)
}

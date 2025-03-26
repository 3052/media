package main

import (
   "net/http"
   "net/url"
   "os"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "fapi.molotov.tv"
   req.URL.Path = "/v2/channels/531/programs/15082/view"
   req.URL.Scheme = "https"
   req.Header["X-Molotov-Agent"] = []string{"{\"app_id\":\"browser_app\",\"app_build\":4,\"app_version_name\":\"4.4.4\",\"browser_name\":\"unknown\",\"type\":\"desktop\",\"os_version\":\"5.0 (Windows)\",\"electron_version\":\"0.0.0\",\"os\":\"Win32\",\"manufacturer\":\"\",\"serial\":\"c9f8605d-0121-4830-96d9-fa297ed6cd9a\",\"model\":\"Firefox - Windows\",\"hasTouchbar\":false,\"brand\":\"Windows NT 10.0; Win64; x64; rv:128.0\",\"api_version\":8,\"features_supported\":[\"social\",\"new_button_conversion\",\"paywall\",\"channel_separator\",\"empty_view_v2\",\"store_offer_v2\",\"player_mplus_teasing\",\"embedded_player\",\"channels_classification\",\"new-post-registration\",\"appstart-d0-full-image\",\"payment_v2\",\"armageddon\",\"user_favorite\",\"parental_control_v3\",\"emptyview_v2\",\"before_pay_periodicity_selection\",\"player_midrolls\",\"cookie_wall\",\"reverse_epg\"],\"inner_app_version_name\":\"4.22.0\",\"qa\":false}"}
   value := url.Values{}
   value["access_token"] = []string{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiMjgxODQxMDgiLCJhbGxvd2VkX2NpZHJzIjpbIjAuMC4wLjAvMCJdLCJleHBpcmVzIjoxNzQzMDM2MDgwLCJwcm9maWxlX2lkIjoiMjgxMzc5NjQiLCJzY29wZXMiOm51bGwsInVzZXJfaWQiOiIyODE4NDEwOCIsInYiOjF9.091390wNyt1_Mwbz9FhZtNYpNa6uASc8RJ1fwTb5fKE"}
   req.URL.RawQuery = value.Encode()
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

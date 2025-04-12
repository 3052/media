package movistar

import (
   "io"
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "strings"
)

// XFF fail
// mullvad pass
// nord pass
func Zero() {
   data, err := exec.Command("password", "movistarplus.es").Output()
   if err != nil {
      panic(err)
   }
   username, password, _ := strings.Cut(string(data), ":")
   var body = strings.NewReader(url.Values{
      "grant_type":[]string{"password"},
      "password":[]string{password},
      "username":[]string{username},
   }.Encode())
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "auth.dof6.com"
   req.URL.Path = "/auth/oauth2/token"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(body)
   value := url.Values{}
   value["deviceClass"] = []string{"amazon.tv"}
   req.URL.RawQuery = value.Encode()
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

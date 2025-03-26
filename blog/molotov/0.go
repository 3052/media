package main

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "strings"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "fapi.molotov.tv"
   req.URL.Path = "/v3.1/auth/login"
   req.URL.Scheme = "https"
   req.Header.Set(
      "x-molotov-agent", `{ "app_build": 4, "app_id": "browser_app" }`,
   )
   data, err := exec.Command("password", "molotov.tv").Output()
   if err != nil {
      panic(err)
   }
   email, password, _ := strings.Cut(string(data), ":")
   data1 := fmt.Sprintf(`
   {
      "grant_type": "password",
      "email": %q,
      "password": %q
   }
   `, email, password)
   req.Body = io.NopCloser(strings.NewReader(data1))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

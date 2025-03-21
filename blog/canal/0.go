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
   req.URL.Host = "m7cplogin.solocoo.tv"
   req.URL.Path = "/login"
   req.URL.Scheme = "https"
   req.Header["Authorization"] = []string{"Client key=web.NhFyz4KsZ54,time=1742523689,sig=1JQMSqZg2dp_K_81udPMzy0GagRw4FVu2GH-j7OmpPc"}
   data1, err := exec.Command("password", "canalplus.cz").Output()
   if err != nil {
      panic(err)
   }
   username, password, _ := strings.Cut(string(data1), ":")
   req.Body = io.NopCloser(strings.NewReader(
      fmt.Sprintf(data, username, password),
   ))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

const data = `{"ticket":"eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2V5IjoibTcifQ..lUS_rO5lGmDeoMF5UPZKcQ.hxPO2rOnHHkv6M6ildYwi-_Z-gKZeBFntOKgQ-STd-li3Iz64TK0Dl_9-E_ndF9mv0jT7BuTunjBAkSrS32hvruJQrmKERrg7QWKl0Qo8_YCpQyFJe6mrYewiaqMp3hbtCBUXpmUgNfaBm4Rf4gaWjg4Bfe_7dqFyVSCsilHPFORKDtYUtb_S6ys4CVWacdfDluNbEVmtbOa2OfNQ3vpRJs9zqcN44usInmA-jb8NHhBmXz5Q6TXqIjQ1C9MoEoK5DMkCjrWecwxl3Cclpgt91pW5XI0nBuoWWkpY163CjHlVu7vH0xMSqTRrjRG3_68IUSaZC3F2IA5WbVUEADAHPn8I3Jur3ZQSR_okjnADD4.D8Ktg0VyrzZBaxgZ4Xizow","userInput":{"username":%q,"password":%q}}`

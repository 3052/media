package canal

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os/exec"
   "strings"
)

func One() (*http.Response, error) {
   const data = `
   {
     "ticket": "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2V5IjoibTcifQ..lUS_rO5lGmDeoMF5UPZKcQ.hxPO2rOnHHkv6M6ildYwi-_Z-gKZeBFntOKgQ-STd-li3Iz64TK0Dl_9-E_ndF9mv0jT7BuTunjBAkSrS32hvruJQrmKERrg7QWKl0Qo8_YCpQyFJe6mrYewiaqMp3hbtCBUXpmUgNfaBm4Rf4gaWjg4Bfe_7dqFyVSCsilHPFORKDtYUtb_S6ys4CVWacdfDluNbEVmtbOa2OfNQ3vpRJs9zqcN44usInmA-jb8NHhBmXz5Q6TXqIjQ1C9MoEoK5DMkCjrWecwxl3Cclpgt91pW5XI0nBuoWWkpY163CjHlVu7vH0xMSqTRrjRG3_68IUSaZC3F2IA5WbVUEADAHPn8I3Jur3ZQSR_okjnADD4.D8Ktg0VyrzZBaxgZ4Xizow",
     "userInput": {
       "username": %q,
       "password": %q
     }
   }
   `
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "m7cplogin.solocoo.tv"
   req.URL.Path = "/login"
   req.URL.Scheme = "https"
   data1, err := exec.Command("password", "canalplus.cz").Output()
   if err != nil {
      panic(err)
   }
   username, password, _ := strings.Cut(string(data1), ":")
   data2 := fmt.Sprintf(data, username, password)
   req.Body = io.NopCloser(strings.NewReader(data2))
   var client1 client
   err = client1.New(req.URL, []byte(data2))
   if err != nil {
      panic(err)
   }
   req.Header.Set("authorization", client1.String())
   return http.DefaultClient.Do(&req)
}

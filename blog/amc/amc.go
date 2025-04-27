package amc

import (
   "encoding/json"
   "iter"
   "net/http"
   "strconv"
   "strings"
)

func season_episodes(id int64) (*child, error) {
   req, _ := http.NewRequest("", "https://gw.cds.amcn.com", nil)
   req.URL.Path = func() string {
      b := []byte("/content-compiler-cr/api/v1/content/amcn/amcplus/")
      b = append(b, "type/series-detail/id/"...)
      b = strconv.AppendInt(b, id, 10)
      return string(b)
   }()
   req.Header = http.Header{
      "authorization":   {"Bearer eyJraWQiOiJwcm9kLTEiLCJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJlbnRpdGxlbWVudHMiOiJhbWNuLWF1dGgsb2Itc3ViLWFtY3BsdXMiLCJkZWZhdWx0X3Byb2ZpbGVzIjpbeyJwcm9maWxlSWQiOjEwMzgzMTYxLCJwcm9maWxlTmFtZSI6IkRlZmF1bHQiLCJzZXJ2aWNlSWQiOiJhbWNwbHVzIn1dLCJhdXRoX3R5cGUiOiJiZWFyZXIiLCJhbWNuLWFjY291bnQtY291bnRyeSI6InVzIiwicm9sZXMiOlsiYW1jbi1hdXRoIiwib2Itc3ViLWFtY3BsdXMiXSwiaXNzIjoiaXAtMTAtMi0xNjYtNDAuZWMyLmludGVybmFsIiwidG9rZW5fdHlwZSI6ImF1dGgiLCJhdWQiOiJyZXNvdXJjZV9zZXJ2ZXIiLCJhbWNuLWFjY291bnQtaWQiOiI2YjAwZGUxZS1jNmRiLTQ3NmEtYmZjOC04NDUyMTg2MTJiYzkiLCJmZWF0dXJlX2ZsYWdzIjoiZXlKaGJXTndiSFZ6TFcxMmNHUWlPbnNpWTJoaGNuUmxjaTF0ZG5BdFpXNWhZbXhsWkNJNmRISjFaWDBzSW1GdFkzQnNkWE10ZEdsNlpXNHRZM2NpT25zaVpXNWhZbXhsWkNJNlptRnNjMlY5TENKallXSmZZVzFqY0MxamIyNTBaVzUwTFdkaGRHVmtMVzlqZEMweU1ESTBJam9pUVNJc0ltTmhZbDloYldOd0xXTnZiblJsYm5RdFpYaDBjbUZ6TFcxaGNpMHlNREkxSWpvaVFTSXNJbUZ0WTNCc2RYTXRZV1F0ZEdsbGNpSTZleUpoWkMxMGFXVnlMWEIxY21Ob1lYTmxMVzl1SWpwMGNuVmxmU3dpWVcxamNHeDFjeTF6YTJsd0xYQnliMjF2TFdGa2N5STZleUp6YTJsd0xYQnliMjF2TFdGa2N5MWxibUZpYkdWa0lqcDBjblZsTENKMllYSnBZWFJwYjI0aU9pSkJUVU1ySUVGa0lFWnlaV1VpZlN3aVkyOXRZMkZ6ZEMxaFpDMWliRzlqYTJWeUxYTmpjbVZsYmlJNmV5SnpkV0owYVhSc1pTSTZJa1p2Y2lCb1pXeHdMQ0JsYldGcGJDQmpkWE4wYjIxbGNuTmxjblpwWTJWQVlXMWpjR3gxY3k1amIyMHVJaXdpZEdsMGJHVWlPaUpVYUdVZ1RXOXVkR2hzZVNCM2FYUm9JRUZrY3lCd2JHRnVJR2x6SUc1dmRDQmpkWEp5Wlc1MGJIa2djM1Z3Y0c5eWRHVmtJRzl1SUZobWFXNXBkSGtnWkdWMmFXTmxjeTRpTENKbGJtRmliR1ZrSWpwMGNuVmxmU3dpWVcxamNHeDFjeTEyYVhwcGJ5MXdjbTl0YjNScGIyNGlPbnNpWlc1aFlteGxaQ0k2Wm1Gc2MyVjlmUT09IiwiZXhwIjoxNzQ1Nzk3NTIxLCJpYXQiOjE3NDU3NzU5MjEsImFtY24tc2VydmljZS1ncm91cC1pZCI6IjEwIiwianRpIjoiOGFjYTY5YzQtZjJmZi00Y2YzLTlhYTktMzhmNmJkOTA4M2E4IiwiYW1jbi11c2UtYWNjb3VudC1jb3VudHJ5IjpmYWxzZX0.NneInIS7E-sOuSLnNxLto_VR7xAbV4gUiuh3cEXjb4PIvs_p-TUowydhBIOyb-n_RiKMyuJuRK9Gp5CW_5B35dXG254dCzA4UYUGnGUfc-sd1qz1N3tQWybG-MgyC_GiJ97pMNxY9HfmGgd6jd3LHaeRXR0nToIuKLIkbPgGeXsWvbaOAxhr3CN-a7z4bzkb3f9OperSDv1r6iNwZG8V9Ui36pCN_yqOXmL5Y5j4PoVCpkr2mVSGWCFV_v2NuROcS_1KdXLucTyA3z4wduFoxnffcM2jkSTiGdfNHMb4EIW3tEz5uGDyEzhYqLbPz9i7W8oxnbgIf1CuTZKwMM9AIg"},
      "x-amcn-network":  {"amcplus"},
      "x-amcn-platform": {"android"},
      "x-amcn-tenant":   {"amcn"},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data child
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Data, nil
}

func series_detail(id int64) (*child, error) {
   req, _ := http.NewRequest("", "https://gw.cds.amcn.com", nil)
   req.URL.Path = func() string {
      b := []byte("/content-compiler-cr/api/v1/content/amcn/amcplus/")
      b = append(b, "type/series-detail/id/"...)
      b = strconv.AppendInt(b, id, 10)
      return string(b)
   }()
   req.Header = http.Header{
      "authorization":   {"Bearer eyJraWQiOiJwcm9kLTEiLCJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJlbnRpdGxlbWVudHMiOiJhbWNuLWF1dGgsb2Itc3ViLWFtY3BsdXMiLCJkZWZhdWx0X3Byb2ZpbGVzIjpbeyJwcm9maWxlSWQiOjEwMzgzMTYxLCJwcm9maWxlTmFtZSI6IkRlZmF1bHQiLCJzZXJ2aWNlSWQiOiJhbWNwbHVzIn1dLCJhdXRoX3R5cGUiOiJiZWFyZXIiLCJhbWNuLWFjY291bnQtY291bnRyeSI6InVzIiwicm9sZXMiOlsiYW1jbi1hdXRoIiwib2Itc3ViLWFtY3BsdXMiXSwiaXNzIjoiaXAtMTAtMi0xNjYtNDAuZWMyLmludGVybmFsIiwidG9rZW5fdHlwZSI6ImF1dGgiLCJhdWQiOiJyZXNvdXJjZV9zZXJ2ZXIiLCJhbWNuLWFjY291bnQtaWQiOiI2YjAwZGUxZS1jNmRiLTQ3NmEtYmZjOC04NDUyMTg2MTJiYzkiLCJmZWF0dXJlX2ZsYWdzIjoiZXlKaGJXTndiSFZ6TFcxMmNHUWlPbnNpWTJoaGNuUmxjaTF0ZG5BdFpXNWhZbXhsWkNJNmRISjFaWDBzSW1GdFkzQnNkWE10ZEdsNlpXNHRZM2NpT25zaVpXNWhZbXhsWkNJNlptRnNjMlY5TENKallXSmZZVzFqY0MxamIyNTBaVzUwTFdkaGRHVmtMVzlqZEMweU1ESTBJam9pUVNJc0ltTmhZbDloYldOd0xXTnZiblJsYm5RdFpYaDBjbUZ6TFcxaGNpMHlNREkxSWpvaVFTSXNJbUZ0WTNCc2RYTXRZV1F0ZEdsbGNpSTZleUpoWkMxMGFXVnlMWEIxY21Ob1lYTmxMVzl1SWpwMGNuVmxmU3dpWVcxamNHeDFjeTF6YTJsd0xYQnliMjF2TFdGa2N5STZleUp6YTJsd0xYQnliMjF2TFdGa2N5MWxibUZpYkdWa0lqcDBjblZsTENKMllYSnBZWFJwYjI0aU9pSkJUVU1ySUVGa0lFWnlaV1VpZlN3aVkyOXRZMkZ6ZEMxaFpDMWliRzlqYTJWeUxYTmpjbVZsYmlJNmV5SnpkV0owYVhSc1pTSTZJa1p2Y2lCb1pXeHdMQ0JsYldGcGJDQmpkWE4wYjIxbGNuTmxjblpwWTJWQVlXMWpjR3gxY3k1amIyMHVJaXdpZEdsMGJHVWlPaUpVYUdVZ1RXOXVkR2hzZVNCM2FYUm9JRUZrY3lCd2JHRnVJR2x6SUc1dmRDQmpkWEp5Wlc1MGJIa2djM1Z3Y0c5eWRHVmtJRzl1SUZobWFXNXBkSGtnWkdWMmFXTmxjeTRpTENKbGJtRmliR1ZrSWpwMGNuVmxmU3dpWVcxamNHeDFjeTEyYVhwcGJ5MXdjbTl0YjNScGIyNGlPbnNpWlc1aFlteGxaQ0k2Wm1Gc2MyVjlmUT09IiwiZXhwIjoxNzQ1Nzk3NTIxLCJpYXQiOjE3NDU3NzU5MjEsImFtY24tc2VydmljZS1ncm91cC1pZCI6IjEwIiwianRpIjoiOGFjYTY5YzQtZjJmZi00Y2YzLTlhYTktMzhmNmJkOTA4M2E4IiwiYW1jbi11c2UtYWNjb3VudC1jb3VudHJ5IjpmYWxzZX0.NneInIS7E-sOuSLnNxLto_VR7xAbV4gUiuh3cEXjb4PIvs_p-TUowydhBIOyb-n_RiKMyuJuRK9Gp5CW_5B35dXG254dCzA4UYUGnGUfc-sd1qz1N3tQWybG-MgyC_GiJ97pMNxY9HfmGgd6jd3LHaeRXR0nToIuKLIkbPgGeXsWvbaOAxhr3CN-a7z4bzkb3f9OperSDv1r6iNwZG8V9Ui36pCN_yqOXmL5Y5j4PoVCpkr2mVSGWCFV_v2NuROcS_1KdXLucTyA3z4wduFoxnffcM2jkSTiGdfNHMb4EIW3tEz5uGDyEzhYqLbPz9i7W8oxnbgIf1CuTZKwMM9AIg"},
      "x-amcn-network":  {"amcplus"},
      "x-amcn-platform": {"android"},
      "x-amcn-tenant":   {"amcn"},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data child
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Data, nil
}

func (c *child) String() string {
   var b strings.Builder
   b.WriteString("text = ")
   b.WriteString(c.Properties.Metadata.ItemText)
   b.WriteString("\ntype = ")
   b.WriteString(c.Properties.Metadata.ElementType)
   b.WriteString("\nendpoint = ")
   b.WriteString(c.Callback.Endpoint)
   return b.String()
}

type child struct {
   Properties struct {
      Metadata struct {
         ItemText    string
         ElementType string
      }
   }
   Callback *struct {
      Endpoint string
   }
   Children []child
}

func (c *child) seasons() iter.Seq[*child] {
   return func(yield func(*child) bool) {
      for _, child1 := range c.Children {
         if child1.Properties.Metadata.ElementType == "tab_bar" {
            for _, child2 := range child1.Children {
               if child2.Properties.Metadata.ItemText == "Seasons" {
                  for _, child3 := range child2.Children {
                     for _, child4 := range child3.Children {
                        if !yield(&child4) {
                           return
                        }
                     }
                  }
                  return
               }
            }
         }
      }
   }
}

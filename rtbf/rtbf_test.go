package rtbf

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
   category string
   path     string
   url      string
}{
   {
      category: "films",
      url: "auvio.rtbf.be/media/sibyl-avec-virginie-efira-et-adele-exarchopoulos-3355182",
      path: "/media/sibyl-avec-virginie-efira-et-adele-exarchopoulos-3355182",
   },
   {
      category: "series",
      path: "/media/the-durrells-une-famille-anglaise-a-corfou-the-durrells-une-famille-anglaise-a-corfou-s01-3351856",
      url: "auvio.rtbf.be/media/the-durrells-une-famille-anglaise-a-corfou-the-durrells-une-famille-anglaise-a-corfou-s01-3351856",
   },
}

func Test(t *testing.T) {
   for _, testVar := range tests {
      contentVar, err := Address{testVar.path}.Content()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", contentVar)
      time.Sleep(time.Second)
   }
}

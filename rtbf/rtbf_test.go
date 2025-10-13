package rtbf

import (
   "log"
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
      url:      "auvio.rtbf.be/media/sibyl-avec-virginie-efira-et-adele-exarchopoulos-3355182",
      path:     "/media/sibyl-avec-virginie-efira-et-adele-exarchopoulos-3355182",
   },
   {
      category: "series",
      path:     "/media/the-durrells-une-famille-anglaise-a-corfou-the-durrells-une-famille-anglaise-a-corfou-s01-3351856",
      url:      "auvio.rtbf.be/media/the-durrells-une-famille-anglaise-a-corfou-the-durrells-une-famille-anglaise-a-corfou-s01-3351856",
   },
}

func Test(t *testing.T) {
   for _, testVar := range tests {
      asset_id, err := Address{testVar.path}.AssetId()
      if err != nil {
         t.Fatal(err)
      }
      log.Print(asset_id)
      time.Sleep(time.Second)
   }
}

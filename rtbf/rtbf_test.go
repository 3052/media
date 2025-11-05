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
      category: "series",
      path:     "/media/the-durrells-une-famille-anglaise-a-corfou-the-durrells-une-famille-anglaise-a-corfou-s01-3351856",
      url:      "https://auvio.rtbf.be/media/the-durrells-une-famille-anglaise-a-corfou-the-durrells-une-famille-anglaise-a-corfou-s01-3351856",
   },
   {
      category: "films",
      path: "/media/l-affaire-thomas-crown-avec-steve-mcqueen-et-faye-dunawa-3381405",
      url: "https://auvio.rtbf.be/media/l-affaire-thomas-crown-avec-steve-mcqueen-et-faye-dunawa-3381405",
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

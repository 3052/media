package molotov

import (
   "41.neocities.org/widevine"
   "bytes"
   "net/http"
   "os"
)

func Four() {
   home, err := os.UserHomeDir()
   if err != nil {
      panic(err)
   }
   private_key, err := os.ReadFile(home + "/media/private_key.pem")
   if err != nil {
      panic(err)
   }
   client_id, err := os.ReadFile(home + "/media/client_id.bin")
   if err != nil {
      panic(err)
   }
   var pssh widevine.Pssh
   pssh.KeyIds = [][]byte{
      []byte("\xc3\x1c\xd0+m\x17\x01\xee\xa1\xedp7\xa8~\xd8J"),
   }
   var cdm widevine.Cdm
   err = cdm.New(private_key, client_id, pssh.Marshal())
   if err != nil {
      panic(err)
   }
   data, err := cdm.RequestBody()
   if err != nil {
      panic(err)
   }
   req, err := http.NewRequest(
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(data),
   )
   if err != nil {
      panic(err)
   }
   req.Header["X-Dt-Auth-Token"] = []string{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcnQiOiJbe1wiYXNzZXRJZFwiOlwiY29jZjlyNmh2YW1zdWNzZm9mMDBcIixcInByb2ZpbGVcIjp7XCJyZW50YWxcIjp7XCJhYnNvbHV0ZUV4cGlyYXRpb25cIjpcIjIwMjUtMDMtMjdUMDQ6MTk6MTFaXCIsXCJwbGF5RHVyYXRpb25cIjoxNDQwMDAwMH19LFwibWVzc2FnZVwiOlwiTGljZW5zZSBHcmFudGVkIVwiLFwib3V0cHV0UHJvdGVjdGlvblwiOntcImRpZ2l0YWxcIjp0cnVlLFwiYW5hbG9ndWVcIjp0cnVlLFwiZW5mb3JjZVwiOmZhbHNlfSxcInN0b3JlTGljZW5zZVwiOmZhbHNlfV0iLCJpYXQiOjE3NDMwMzQ3NTEsImp0aSI6ImFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhIiwib3B0RGF0YSI6IntcInVzZXJJZFwiOlwiMjgxODQxMDhcIixcInNlc3Npb25JZFwiOlwiT0VfTkZncXpzekYwNmllYTBtQlM3dUdSVUlvPVwiLFwibWVyY2hhbnRcIjpcIm1vbG90b3ZcIn0ifQ.f6oSV7uO8epqqw1Zq_WlJfPFMdngtPadYBGVK8MvBAgLp-dTBb7DPMrc0lzDf-Xvhgq_9Y8VKtujFu4rIhr6Xw"}
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   resp.Write(os.Stdout)
}

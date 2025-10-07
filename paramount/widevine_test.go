package paramount

import (
   "41.neocities.org/drm/widevine"
   "bytes"
   "encoding/hex"
   "net/http"
   "os"
   "testing"
)

func TestWidevine(t *testing.T) {
   const content_id = "Ddx7cwK2iWCMANoD0Q2hQTR4FLETD_gj"
   var pssh widevine.Pssh
   key_id, err := hex.DecodeString("8992ab68697c476f832acfc7903ea9a5")
   if err != nil {
      t.Fatal(err)
   }
   pssh.KeyIds = [][]byte{key_id}
   pssh.ContentId = []byte(content_id)
   var module widevine.Cdm
   client_id, err := os.ReadFile(
      `C:\Users\Steven\media\7470\device_client_id_blob`,
   )
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := os.ReadFile(
      `C:\Users\Steven\media\7470\device_private_key`,
   )
   if err != nil {
      t.Fatal(err)
   }
   err = module.New(private_key, client_id, pssh.Marshal())
   if err != nil {
      t.Fatal(err)
   }
   data, err := module.RequestBody()
   if err != nil {
      t.Fatal(err)
   }
   atVar, err := ComCbsApp.At()
   if err != nil {
      t.Fatal(err)
   }
   sessionVar, err := atVar.Session(content_id)
   if err != nil {
      t.Fatal(err)
   }
   req, err := http.NewRequest("POST", sessionVar.Url, bytes.NewReader(data))
   if err != nil {
      t.Fatal(err)
   }
   req.Header.Set("authorization", "Bearer " + sessionVar.LsSession)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}

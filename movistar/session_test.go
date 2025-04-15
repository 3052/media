package movistar

import (
   "41.neocities.org/widevine"
   "encoding/base64"
   "log"
   "net/http"
   "os"
   "testing"
)

func TestSession(t *testing.T) {
   log.SetFlags(log.Ltime)
   http.DefaultClient.Transport = transport{}
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/movistar/token")
   if err != nil {
      t.Fatal(err)
   }
   var token1 Token
   err = token1.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   data, err = os.ReadFile(home + "/media/movistar/device")
   if err != nil {
      t.Fatal(err)
   }
   var device1 Device
   err = device1.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   oferta1, err := token1.Oferta()
   if err != nil {
      t.Fatal(err)
   }
   init1, err := oferta1.InitData(device1)
   if err != nil {
      t.Fatal(err)
   }
   var details1 Details
   err = details1.New(test.id)
   if err != nil {
      t.Fatal(err)
   }
   session1, err := device1.Session(init1, &details1)
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := os.ReadFile(home + "/media/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(home + "/media/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   key_id, err := base64.StdEncoding.DecodeString(test.key_id)
   if err != nil {
      t.Fatal(err)
   }
   var pssh widevine.Pssh
   pssh.KeyIds = [][]byte{key_id}
   pssh.ContentId, err = base64.StdEncoding.DecodeString(test.content_id)
   if err != nil {
      t.Fatal(err)
   }
   var cdm widevine.Cdm
   err = cdm.New(private_key, client_id, pssh.Marshal())
   if err != nil {
      t.Fatal(err)
   }
   data, err = cdm.RequestBody()
   if err != nil {
      t.Fatal(err)
   }
   _, err = session1.Widevine(data)
   if err != nil {
      t.Fatal(err)
   }
}

var test = struct {
   content_id string
   id         int64
   key_id     string
   url        string
}{
   content_id: "MTE3NjU2OA==",
   id:         3427440,
   key_id:     "Yc2mUFQwSrKc25rgupRzRQ==",
   url:        "movistarplus.es/cine/ficha?id=3427440",
}

func (transport) RoundTrip(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultTransport.RoundTrip(req)
}

type transport struct{}

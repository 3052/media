package paramount

import (
   "41.neocities.org/drm/widevine"
   "bytes"
   "encoding/hex"
   "fmt"
   "net/http"
   "os"
   "testing"
)

func TestPlayReady(t *testing.T) {
   token, err := ComCbsApp.At()
   if err != nil {
      t.Fatal(err)
   }
   sessionVar, err := token.playReady("wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q")
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache + "/paramount/PlayReady", []byte(sessionVar.LsSession), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

func TestWidevine(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(cache + "/L3/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := os.ReadFile(cache + "/L3/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   const content_id = "Ddx7cwK2iWCMANoD0Q2hQTR4FLETD_gj"
   var pssh widevine.Pssh
   key_id, err := hex.DecodeString("8992ab68697c476f832acfc7903ea9a5")
   if err != nil {
      t.Fatal(err)
   }
   pssh.KeyIds = [][]byte{key_id}
   pssh.ContentId = []byte(content_id)
   var module widevine.Cdm
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

var location_tests = []struct {
   content_id string
   location   []string
   url        string
   period int
}{
   {
      location:   []string{"USA"},
      url:        "paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
      content_id: "wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
   },
   {
      content_id: "rZ59lcp4i2fU4dAaZJ_iEgKqVg_ogrIf",
      location:   []string{"USA"},
      url:        "cbs.com/shows/video/rZ59lcp4i2fU4dAaZJ_iEgKqVg_ogrIf",
      period: 2,
   },
   {
      content_id: "3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ",
      url:        "paramountplus.com/movies/video/3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ",
      location: []string{
         "Brazil", "Canada", "Chile", "Colombia", "Mexico", "Peru",
      },
   },
   {
      content_id: "Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ",
      url:        "paramountplus.com/movies/video/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ",
      location:   []string{"Australia", "United Kingdom"},
   },
   {
      content_id: "WNujiS5PHkY5wN9doNY6MSo_7G8uBUcX",
      url:        "paramountplus.com/shows/video/WNujiS5PHkY5wN9doNY6MSo_7G8uBUcX",
      location:   []string{"Australia"},
   },
   {
      content_id: "esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
      location:   []string{"USA"},
      url:        "paramountplus.com/shows/video/esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
      period: 1,
   },
}

func TestLocation(t *testing.T) {
   fmt.Println(location_tests)
}

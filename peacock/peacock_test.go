package peacock

import (
   "41.neocities.org/drm/widevine"
   "encoding/hex"
   "fmt"
   "net/http"
   "os"
   "os/exec"
   "testing"
)

func TestSignRead(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/peacock/peacock.txt")
   if err != nil {
      t.Fatal(err)
   }
   id, err := http.ParseSetCookie(string(data))
   if err != nil {
      t.Fatal(err)
   }
   var token AuthToken
   err = token.Fetch(id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", token)
}

func TestSignWrite(t *testing.T) {
   user, err := output("credential", "-h=peacocktv.com", "-k=user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h=peacocktv.com")
   if err != nil {
      t.Fatal(err)
   }
   id, err := FetchIdSession(user, password)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache+"/peacock/peacock.txt", []byte(id.String()), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

var watch = struct {
   content_id string
   key_id     string
   url        string
}{
   content_id: "GMO_00000000091566_02_HDSDR",
   key_id:     "3016b29577190c2f1a4653203a2313f7",
   url:        "https://peacocktv.com/watch/playback/vod/GMO_00000000091566_02_HDSDR/6668f89a-b581-36ac-9895-7783aa16b471",
}

func TestLicense(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(cache + "/L3/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   pem_bytes, err := os.ReadFile(cache + "/L3/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := widevine.ParsePrivateKey(pem_bytes)
   if err != nil {
      t.Fatal(err)
   }
   key_id, err := hex.DecodeString(watch.key_id)
   if err != nil {
      t.Fatal(err)
   }
   // 1. Create the PsshData struct
   pssh := &widevine.PsshData{KeyIds: [][]byte{key_id}}
   // 2. Build the License Request directly from the pssh struct
   req_bytes, err := pssh.BuildLicenseRequest(client_id)
   if err != nil {
      t.Fatal(err)
   }
   // 3. Sign the request
   signed_bytes, err := widevine.BuildSignedMessage(req_bytes, private_key)
   if err != nil {
      t.Fatalf("Failed to create signed request: %v", err)
   }
   // 4. Send to License Server
   data, err := os.ReadFile(cache + "/peacock/peacock.txt")
   if err != nil {
      t.Fatal(err)
   }
   id, err := http.ParseSetCookie(string(data))
   if err != nil {
      t.Fatal(err)
   }
   var token AuthToken
   err = token.Fetch(id)
   if err != nil {
      t.Fatal(err)
   }
   play, err := token.Playout(watch.content_id)
   if err != nil {
      t.Fatal(err)
   }
   _, err = play.Widevine(signed_bytes)
   if err != nil {
      t.Fatal(err)
   }
}

func output(name string, arg ...string) (string, error) {
   data, err := exec.Command(name, arg...).Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
}

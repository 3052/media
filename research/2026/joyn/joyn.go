package joyn

import (
   "bytes"
   "crypto/sha1"
   "encoding/hex"
   "encoding/json"
   "errors"
   "net/http"
   "time"
)

func (a Anonymous) Entitlement(content_id string) (*Entitlement, error) {
   body, err := json.Marshal(map[string]string{"content_id": content_id})
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://entitlement.p7s1.io/api/user/entitlement-token",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+a.Access_Token)
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      var b bytes.Buffer
      res.Write(&b)
      return nil, errors.New(b.String())
   }
   title := new(Entitlement)
   err = json.NewDecoder(res.Body).Decode(title)
   if err != nil {
      return nil, err
   }
   return title, nil
}

type Entitlement struct {
   Entitlement_Token string
}

type Anonymous struct {
   Access_Token string
}
func (e Entitlement) Playlist(content_id string) (*Playlist, error) {
   body, err := func() ([]byte, error) {
      var s struct {
         Manufacturer     string `json:"manufacturer"`
         MaxResolution    int    `json:"maxResolution"`
         Model            string `json:"model"`
         Platform         string `json:"platform"`
         ProtectionSystem string `json:"protectionSystem"`
         StreamingFormat  string `json:"streamingFormat"`
      }
      s.Manufacturer = "unknown"
      s.MaxResolution = 1080
      s.Model = "unknown"
      s.Platform = "browser"
      s.ProtectionSystem = "widevine"
      s.StreamingFormat = "dash"
      return json.Marshal(s)
   }()
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://api.vod-prd.s.joyn.de", bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/v1/asset/" + content_id + "/playlist"
   req.URL.RawQuery = "signature=" + e.signature(body)
   req.Header = http.Header{
      "authorization": {"Bearer " + e.Entitlement_Token},
      "content-type":  {"application/json"},
   }
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   play := new(Playlist)
   err = json.NewDecoder(res.Body).Decode(play)
   if err != nil {
      return nil, err
   }
   return play, nil
}

func (Playlist) WrapRequest(b []byte) ([]byte, error) {
   return b, nil
}

func (Playlist) UnwrapResponse(b []byte) ([]byte, error) {
   return b, nil
}

type Playlist struct {
   LicenseUrl  string
   ManifestUrl string
}

func (p Playlist) RequestUrl() (string, bool) {
   return p.LicenseUrl, true
}

func (Playlist) RequestHeader() (http.Header, error) {
   return http.Header{}, nil
}

const signature_key = "5C7838365C7864665C786638265C783064595C783935245C7865395C7838323F5C7866333D3B5C78386635"

func (e Entitlement) signature(text []byte) string {
   text = append(text, ',')
   text = append(text, e.Entitlement_Token...)
   text = hex.AppendEncode(text, []byte(signature_key))
   sum := sha1.Sum(text)
   return hex.EncodeToString(sum[:])
}

func (a *Anonymous) New() error {
   body, err := func() ([]byte, error) {
      m := map[string]string{
         // fuck you:
         // ENT_RVOD_Playback_Restricted
         "client_id": time.Now().String(),
         "client_name":"web",
      }
      return json.Marshal(m)
   }()
   if err != nil {
      return err
   }
   res, err := http.Post(
      "https://auth.joyn.de/auth/anonymous", "application/json",
      bytes.NewReader(body),
   )
   if err != nil {
      return err
   }
   defer res.Body.Close()
   return json.NewDecoder(res.Body).Decode(a)
}

package disney

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "path"
   "strings"
)

func (a *Account) Playback(mediaId string) (*Playback, error) {
   playback_id, err := json.Marshal(map[string]string{
      "mediaId": mediaId,
   })
   if err != nil {
      return nil, err
   }
   data, err := json.Marshal(map[string]any{
      "playbackId": playback_id,
      "playback": map[string]any{
         "attributes": map[string]any{
            "assetInsertionStrategy": "SGAI",
            "codecs": map[string]bool{
               "supportsMultiCodecMaster": true, // 4K
            },
         },
      },
   })
   if err != nil {
      return nil, err
   }
   req, _ := http.NewRequest(
      "POST",
      // ctr-high also works
      "https://disney.playback.edge.bamgrid.com/v7/playback/ctr-regular",
      bytes.NewReader(data),
   )
   req.Header.Set("authorization", "Bearer "+a.Extensions.Sdk.Token.AccessToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-dss-feature-filtering", "true")
   req.Header.Set("x-bamsdk-platform", "")
   req.Header.Set("x-application-version", "")
   req.Header.Set("x-bamsdk-client-id", "")
   req.Header.Set("x-bamsdk-version", "")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Playback
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result, nil
}

type Hls struct {
   Body []byte
   Url  *url.URL
}

type Playback struct {
   Errors []Error
   Stream struct {
      Sources []struct {
         Complete struct {
            Url string
         }
      }
   }
}

func (a *Account) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST",
      "https://disney.playback.edge.bamgrid.com/widevine/v1/obtain-license",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", a.Extensions.Sdk.Token.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var result struct {
         Errors []Error
      }
      err = json.Unmarshal(data, &result)
      if err != nil {
         return nil, err
      }
      return nil, &result.Errors[0]
   }
   return data, nil
}

func (p *Playback) Hls() (*Hls, error) {
   resp, err := http.Get(p.Stream.Sources[0].Complete.Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Hls{data, resp.Request.URL}, nil
}

func (e *Error) Error() string {
   var data strings.Builder
   data.WriteString("code = ")
   data.WriteString(e.Code)
   data.WriteString("\ndescription = ")
   data.WriteString(e.Description)
   return data.String()
}

type Error struct {
   Code        string
   Description string
}

// ZGlzbmV5JmJyb3dzZXImMS4wLjA
// disney&browser&1.0.0
const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"

const mutation_register_device = `
mutation registerDevice($input: RegisterDeviceInput!) {
   registerDevice(registerDevice: $input) {
      token {
         accessToken
         refreshToken
         accessTokenType
      }
   }
}
`

func (d *Device) Register() error {
   data, err := json.Marshal(map[string]any{
      "query": mutation_register_device,
      "variables": map[string]any{
         "input": map[string]any{
            "deviceProfile":      "!",
            "deviceFamily":       "!",
            "applicationRuntime": "!",
            "attributes": map[string]string{
               "operatingSystem":        "",
               "operatingSystemVersion": "",
            },
         },
      },
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://disney.api.edge.bamgrid.com/graph/v1/device/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("authorization", "Bearer "+client_api_key)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(d)
}

func (a *AccountWithoutActiveProfile) SwitchProfile() (*Account, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_switch_profile,
      "variables": map[string]any{
         "input": map[string]string{
            "profileId": a.Data.Login.Account.Profiles[0].Id,
         },
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+a.Extensions.Sdk.Token.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Account{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (d *Device) Login(email, password string) (*AccountWithoutActiveProfile, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_login,
      "variables": map[string]any{
         "input": map[string]string{
            "email":    email,
            "password": password,
         },
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set(
      "authorization", "Bearer "+d.Data.RegisterDevice.Token.AccessToken,
   )
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &AccountWithoutActiveProfile{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

const mutation_login = `
mutation login($input: LoginInput!) {
   login(login: $input) {
      account {
         profiles {
            id
         }
      }
   }
}
`

const mutation_switch_profile = `
mutation switchProfile($input: SwitchProfileInput!) {
   switchProfile(switchProfile: $input) {
      account {
         activeProfile {
            name
         }
      }
   }
}
`

type Account struct {
   Extensions struct {
      Sdk struct {
         Token struct {
            AccessToken     string
            AccessTokenType string // Account
         }
      }
   }
}

type AccountWithoutActiveProfile struct {
   Data struct {
      Login struct {
         Account struct {
            Profiles []struct {
               Id string
            }
         }
      }
   }
   Extensions struct {
      Sdk struct {
         Token struct {
            AccessToken     string
            AccessTokenType string // AccountWithoutActiveProfile
         }
      }
   }
}

type Device struct {
   Data struct {
      RegisterDevice struct {
         Token struct {
            AccessToken     string
            RefreshToken    string
            AccessTokenType string // Device
         }
      }
   }
}

///

type Explore struct {
   Data struct {
      Errors []Error // region
      Page Page
   }
   Errors []Error // explore-not-supported
}

type Action struct {
   ResourceId string
   Visuals    struct {
      DisplayText string
   }
}

type Page struct {
   Actions []Action
   Containers []struct {
      Seasons []struct {
         Items []struct {
            Actions []Action
            Visuals struct {
               EpisodeNumber string
               EpisodeTitle string
               SeasonNumber string
            }
         }
      }
   }
   Visuals struct {
      Title string
   }
}

// https://disneyplus.com/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
// https://disneyplus.com/cs-cz/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
// https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d
func GetEntity(rawLink string) (string, error) {
   // Parse the URL to safely access its components
   link, err := url.Parse(rawLink)
   if err != nil {
      return "", err
   }
   // Get the last part of the URL path
   last_segment := path.Base(link.Path)
   // The entity might be prefixed with "entity-", so we remove it
   return strings.TrimPrefix(last_segment, "entity-"), nil
}

func (a *Account) Explore(entity string) (*Explore, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+a.Extensions.Sdk.Token.AccessToken)
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "disney.api.edge.bamgrid.com",
      Path:   "/explore/v1.12/page/entity-" + entity,
      RawQuery: url.Values{
         "enhancedContainersLimit": {"1"},
         "limit": {"99"},
      }.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Explore
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   if len(result.Data.Errors) >= 1 {
      return nil, &result.Data.Errors[0]
   }
   return &result, nil
}

// to change location you must log in again
package disney

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
   _ "embed"
)

func (a *AccountWithoutActiveProfile) SwitchProfile(profileId string) (*Account, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_switch_profile,
      "variables": map[string]any{
         "input": map[string]string{
            "profileId": profileId,
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

type Profile struct {
   Name string
   Id   string
}

type Token struct {
   AccessToken     string
   AccessTokenType string // Device
}

//go:embed registerDevice.gql
var mutation_register_device string

func (t *Token) RegisterDevice() error {
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
   var result struct {
      Data struct {
         RegisterDevice struct {
            Token Token
         }
      }
      Errors []Error
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return err
   }
   if len(result.Errors) >= 1 {
      return &result.Errors[0]
   }
   *t = result.Data.RegisterDevice.Token
   return nil
}

//go:embed login.gql
var mutation_login string

type Login struct {
   Account struct {
      Profiles []Profile
   }
}

func (t *Token) Login(email, password string) (*Login, error) {
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
   req.Header.Set("authorization", "Bearer " + t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Login Login
      }
      Errors     []Error
      Extensions struct {
         Sdk struct {
            Token Token
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   *t = result.Extensions.Sdk.Token
   return &result.Data.Login, nil
}

//go:embed requestOtp.gql
var mutation_request_otp string

func (r RequestOtp) String() string {
   if r.Accepted {
      return "accepted = true"
   }
   return "accepted = false"
}

type RequestOtp struct {
   Accepted bool
}

func (t *Token) RequestOtp(email string) (*RequestOtp, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_request_otp,
      "variables": map[string]any{
         "input": map[string]string{
            "email": email,
            "reason": "Login",
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
   req.Header.Set("authorization", "Bearer " + t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         RequestOtp RequestOtp
      }
      Errors []Error
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Data.RequestOtp, nil
}

//go:embed authenticateWithOtp.gql
var mutation_authenticate_with_otp string

// passcode can start with 0
func (t *Token) AuthenticateWithOtp(email, passcode string) (*AuthenticateWithOtp, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_authenticate_with_otp,
      "variables": map[string]any{
         "input": map[string]string{
            "email": email,
            "passcode": passcode,
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
   req.Header.Set("authorization", "Bearer " + t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   var result struct {
      Data struct {
         AuthenticateWithOtp AuthenticateWithOtp
      }
      Errors []Error
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Data.AuthenticateWithOtp, nil
}

//go:embed loginWithActionGrant.gql
var mutation_login_with_action_grant string

type LoginWithActionGrant struct {
   Account struct {
      Profiles []Profile
   }
}

type AuthenticateWithOtp struct {
   ActionGrant string
}
func (t *Token) LoginWithActionGrant(actionGrant string) (*LoginWithActionGrant, error) {
   data, err := json.Marshal(map[string]any{
      "query": mutation_login_with_action_grant,
      "variables": map[string]any{
         "input": map[string]string{
            "actionGrant": actionGrant,
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
   req.Header.Set("authorization", "Bearer " + t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         LoginWithActionGrant LoginWithActionGrant 
      }
      Extensions struct {
         Sdk struct {
            Token Token
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   *t = result.Extensions.Sdk.Token
   return &result.Data.LoginWithActionGrant, nil
}
// ZGlzbmV5JmJyb3dzZXImMS4wLjA
// disney&browser&1.0.0
const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"

func (e *Error) Error() string {
   var data strings.Builder
   if e.Code != "" {
      data.WriteString("code = ")
      data.WriteString(e.Code)
   }
   if e.Description != "" {
      if data.Len() >= 1 {
         data.WriteByte('\n')
      }
      data.WriteString("description = ")
      data.WriteString(e.Description)
   }
   if e.Extensions != nil {
      if data.Len() >= 1 {
         data.WriteByte('\n')
      }
      data.WriteString("extensions = ")
      data.WriteString(e.Extensions.Code)
   }
   if e.Message != "" {
      if data.Len() >= 1 {
         data.WriteByte('\n')
      }
      data.WriteString("message = ")
      data.WriteString(e.Message)
   }
   return data.String()
}

type Error struct {
   Code        string
   Description string
   Extensions  *struct {
      Code string
   }
   Message string
}

type Hls struct {
   Body []byte
   Url  *url.URL
}

func (p *Page) String() string {
   var data strings.Builder
   if len(p.Containers[0].Seasons) >= 1 {
      var line bool
      for _, seasonItem := range p.Containers[0].Seasons {
         if line {
            data.WriteString("\n\n")
         } else {
            line = true
         }
         data.WriteString("name = ")
         data.WriteString(seasonItem.Visuals.Name)
         data.WriteString("\nid = ")
         data.WriteString(seasonItem.Id)
      }
   } else {
      data.WriteString(p.Actions[0].InternalTitle)
   }
   return data.String()
}

type Page struct {
   Actions []struct {
      InternalTitle string // movie
   }
   Containers []struct {
      Seasons []struct { // series
         Visuals struct {
            Name string
         }
         Id string
      }
   }
}

func (p *Profile) String() string {
   var data strings.Builder
   data.WriteString("name = ")
   data.WriteString(p.Name)
   data.WriteString("\nid = ")
   data.WriteString(p.Id)
   return data.String()
}

func (s Season) String() string {
   var (
      data strings.Builder
      line bool
   )
   for _, item := range s.Items {
      for _, action := range item.Actions {
         if line {
            data.WriteByte('\n')
         } else {
            line = true
         }
         data.WriteString(action.InternalTitle)
      }
   }
   return data.String()
}

type Season struct {
   Items []struct {
      Actions []struct {
         InternalTitle string
      }
   }
}

func (s *Stream) Hls() (*Hls, error) {
   resp, err := http.Get(s.Sources[0].Complete.Url)
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

type Stream struct {
   Sources []struct {
      Complete struct {
         Url string
      }
   }
}

// https://disneyplus.com/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
// https://disneyplus.com/cs-cz/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
// https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d
func GetEntity(link string) (string, error) {
   // First, explicitly fail if the URL is a "play" link.
   if strings.Contains(link, "/play/") {
      return "", errors.New("URL is a 'play' link and not a 'browse' link")
   }
   // The unique marker for the ID we want is "/browse/entity-".
   const marker = "/browse/entity-"
   // strings.Cut splits the string at the first instance of the marker.
   // It returns the part before, the part after, and a boolean indicating if the marker was found.
   // We don't need the 'before' part, so we discard it with the blank identifier _.
   _, id, found := strings.Cut(link, marker)
   // If the marker was not found, or if the resulting ID string is empty, return an error.
   if !found || id == "" {
      return "", errors.New("failed to find a valid ID in the URL")
   }
   // The 'id' variable now holds the rest of the string after the marker.
   return id, nil
}

// SL2000 720p
// SL3000 2160p
func (a *Account) PlayReady(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST",
      "https://disney.playback.edge.bamgrid.com/playready/v1/obtain-license.asmx",
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
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

// L3 720p
// L1 2160p
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

func (a *Account) Season(id string) (*Season, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+a.Extensions.Sdk.Token.AccessToken)
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "disney.api.edge.bamgrid.com",
      Path:     "/explore/v1.12/season/" + id,
      RawQuery: "limit=99",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Season Season
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Season, nil
}

func (a *Account) Page(entity string) (*Page, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+a.Extensions.Sdk.Token.AccessToken)
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "disney.api.edge.bamgrid.com",
      Path:     "/explore/v1.12/page/entity-" + entity,
      RawQuery: "limit=0",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Errors []Error // region
         Page   Page
      }
      Errors []Error // auth.expired
   }
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
   return &result.Data.Page, nil
}

func (a *Account) Stream(mediaId string) (*Stream, error) {
   playback_id, err := json.Marshal(map[string]string{
      "mediaId": mediaId,
   })
   if err != nil {
      return nil, err
   }
   data, err := json.Marshal(map[string]any{
      "playback": map[string]any{
         "attributes": map[string]any{
            "assetInsertionStrategy": "SGAI",
            "codecs": map[string]any{
               "supportsMultiCodecMaster": true, // 4K
               "video": []string{
                  "h.264",
                  "h.265",
               },
            },
            "videoRanges": []string{"HDR10"},
         },
      },
      "playbackId": playback_id,
   })
   if err != nil {
      return nil, err
   }
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "disney.playback.edge.bamgrid.com",
      // Path: "/v7/playback/ctr-high",
      // Path: "/v7/playback/tv-drm-ctr-h265-atmos",
      Path: "/v7/playback/ctr-regular",
   }
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+a.Extensions.Sdk.Token.AccessToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-application-version", "")
   req.Header.Set("x-bamsdk-client-id", "")
   req.Header.Set("x-bamsdk-platform", "")
   req.Header.Set("x-bamsdk-version", "")
   req.Header.Set("x-dss-feature-filtering", "true")
   req.Body = io.NopCloser(bytes.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Errors []Error
      Stream Stream
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Stream, nil
}

// access token expires in 14400 seconds AKA 240 minutes AKA 4 hours
func (a *Account) RefreshToken() error {
   data, err := json.Marshal(map[string]any{
      "query": mutation_refresh_token,
      "variables": map[string]any{
         "input": map[string]string{
            "refreshToken": a.Extensions.Sdk.Token.RefreshToken,
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
   return json.NewDecoder(resp.Body).Decode(a)
}

type Account struct {
   Extensions struct {
      Sdk struct {
         Token struct {
            AccessToken     string
            AccessTokenType string // Account
            RefreshToken    string
         }
      }
   }
}

//go:embed refreshToken.gql
var mutation_refresh_token string

//go:embed switchProfile.gql
var mutation_switch_profile string

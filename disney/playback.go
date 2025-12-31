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

type Explore struct {
   Data struct {
      Errors []Error // region
      Page Page
   }
   Errors []Error // explore-not-supported
}

// GetFormattedString generates a printable string based on the presence of Seasons
// in the first container. Each value is printed on a new line.
func (p *Page) GetFormattedString() (string, error) {
   // The only length check, as requested.
   if len(p.Containers[0].Seasons) > 0 {
      // PATH 1: Process the contents of the first container.
      var itemStrings []string
      firstContainer := p.Containers[0]
      for _, season := range firstContainer.Seasons {
         for _, item := range season.Items {
            resourceId := item.Actions[0].ResourceId
            mediaId, err := extractMediaId(resourceId)
            if err != nil {
               return "", fmt.Errorf("error in item S%sE%s: %w", item.Visuals.SeasonNumber, item.Visuals.EpisodeNumber, err)
            }
            // Create a multi-line string block for each item.
            itemBlock := fmt.Sprintf(
               "SeasonNumber: %s\nEpisodeNumber: %s\nEpisodeTitle: %s\nMediaId: %s",
               item.Visuals.SeasonNumber,
               item.Visuals.EpisodeNumber,
               item.Visuals.EpisodeTitle,
               mediaId,
            )
            itemStrings = append(itemStrings, itemBlock)
         }
      }
      // Join the blocks for each item with a blank line in between.
      return strings.Join(itemStrings, "\n\n"), nil
   } else {
      // PATH 2: Get a formatted string from top-level values.
      resourceId := p.Actions[0].ResourceId
      mediaId, err := extractMediaId(resourceId)
      if err != nil {
         return "", fmt.Errorf("top-level path failed: %w", err)
      }
      // Create a multi-line string for the top-level info.
      return fmt.Sprintf("Title: %s\nMediaId: %s", p.Visuals.Title, mediaId), nil
   }
}

// extractMediaId is a helper function to decode the ResourceId and extract the mediaId.
func extractMediaId(encodedResourceId string) (string, error) {
   jsonBytes, err := base64.StdEncoding.DecodeString(encodedResourceId)
   if err != nil {
      return "", fmt.Errorf("base64 decoding failed: %w", err)
   }
   // Helper struct to unmarshal the JSON from the decoded ResourceId.
   var payload struct {
      MediaId string `json:"mediaId"`
   }
   if err := json.Unmarshal(jsonBytes, &payload); err != nil {
      return "", fmt.Errorf("JSON unmarshaling failed: %w", err)
   }
   if payload.MediaId == "" {
      return "", errors.New("JSON is valid but is missing the 'mediaId' key")
   }
   return payload.MediaId, nil
}

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

type Action struct {
   ResourceId string
   Visuals    struct {
      DisplayText string
   }
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
         "limit":                   {"1"},
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

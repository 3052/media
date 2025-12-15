package ctv

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (m Manifest) Mpd() (*url.URL, []byte, error) {
   resp, err := http.Get(strings.Replace(string(m), "/best/", "/ultimate/", 1))
   if err != nil {
      return nil, nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, err
   }
   return resp.Request.URL, data, nil
}

type Manifest []byte

func (a *AxisContent) Manifest(play *Playback) (Manifest, error) {
   req, _ := http.NewRequest("", "https://capi.9c9media.com", nil)
   req.URL.Path = func() string {
      data := &strings.Builder{}
      fmt.Fprint(data, "/destinations/")
      fmt.Fprint(data, a.AxisPlaybackLanguages[0].DestinationCode)
      fmt.Fprint(data, "/platforms/desktop/playback/contents/")
      fmt.Fprint(data, a.AxisId)
      fmt.Fprint(data, "/contentPackages/")
      fmt.Fprint(data, play.ContentPackages[0].Id)
      fmt.Fprint(data, "/manifest.mpd")
      return data.String()
   }()
   req.URL.RawQuery = "action=reference"
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var result struct {
         Message string
      }
      err = json.Unmarshal(data, &result)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(result.Message)
   }
   return data, nil
}

// https://ctv.ca/shows/friends/the-one-with-the-bullies-s2e21
func GetPath(rawLink string) (string, error) {
   link, err := url.Parse(rawLink)
   if err != nil {
      return "", err
   }
   if link.Scheme == "" {
      return "", errors.New("invalid URL: scheme is missing")
   }
   return link.Path, nil
}

const query_resolve_path = `
query resolvePath($path: String!) {
   resolvedPath(path: $path) {
      lastSegment {
         content {
            ... on AxisObject {
               id
               ... on AxisMedia {
                  firstPlayableContent {
                     id
                  }
               }
            }
         }
      }
   }
}
`

const query_axis_content = `
query axisContent($id: ID!) {
   axisContent(id: $id) {
      axisId
      axisPlaybackLanguages {
         ... on AxisPlayback {
            destinationCode
         }
      }
   }
}
`

func Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      "https://license.9c9media.ca/widevine", "application/x-protobuf",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type AxisContent struct {
   AxisId                int
   AxisPlaybackLanguages []struct {
      DestinationCode string
   }
}

type Playback struct {
   ContentPackages []struct {
      Id int
   }
}

func (a *AxisContent) Playback() (*Playback, error) {
   req, _ := http.NewRequest("", "https://capi.9c9media.com", nil)
   req.URL.Path = func() string {
      data := &strings.Builder{}
      fmt.Fprint(data, "/destinations/")
      fmt.Fprint(data, a.AxisPlaybackLanguages[0].DestinationCode)
      fmt.Fprint(data, "/platforms/desktop/contents/")
      fmt.Fprint(data, a.AxisId)
      return data.String()
   }()
   req.URL.RawQuery = "$include=[ContentPackages]"
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Playback{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func Resolve(path string) (*ResolvedPath, error) {
   data, err := json.Marshal(map[string]any{
      "query": query_resolve_path,
      "variables": map[string]string{
         "path": path,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://www.ctv.ca/space-graphql/apq/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   // you need this for the first request, then can omit
   req.Header.Set("graphql-client-platform", "entpay_web")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   var result struct {
      Data struct {
         ResolvedPath *ResolvedPath
      }
   }
   err = json.Unmarshal(data, &result)
   if err != nil {
      return nil, err
   }
   if result.Data.ResolvedPath == nil {
      return nil, errors.New(string(data))
   }
   return result.Data.ResolvedPath, nil
}

type ResolvedPath struct {
   LastSegment struct {
      Content struct {
         FirstPlayableContent *struct {
            Id string
         }
         Id string
      }
   }
}

func (r *ResolvedPath) get_id() string {
   if fpc := r.LastSegment.Content.FirstPlayableContent; fpc != nil {
      return fpc.Id
   }
   return r.LastSegment.Content.Id
}

func (r *ResolvedPath) AxisContent() (*AxisContent, error) {
   data, err := json.Marshal(map[string]any{
      "query": query_axis_content,
      "variables": map[string]string{
         "id": r.get_id(),
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://www.ctv.ca/space-graphql/apq/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   // you need this for the first request, then can omit
   req.Header.Set("graphql-client-platform", "entpay_web")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         AxisContent AxisContent
      }
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, errors.New(result.Errors[0].Message)
   }
   return &result.Data.AxisContent, nil
}

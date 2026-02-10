package ctv

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strconv"
   "strings"
)

func (a *AxisContent) Mpd(content1 *Content) (string, error) {
   req, _ := http.NewRequest("", "https://capi.9c9media.com", nil)
   req.URL.Path = func() string {
      b := []byte("/destinations/")
      b = append(b, a.AxisPlaybackLanguages[0].DestinationCode...)
      b = append(b, "/platforms/desktop/playback/contents/"...)
      b = strconv.AppendInt(b, a.AxisId, 10)
      b = append(b, "/contentPackages/"...)
      b = strconv.AppendInt(b, content1.ContentPackages[0].Id, 10)
      b = append(b, "/manifest.mpd"...)
      return string(b)
   }()
   req.URL.RawQuery = "action=reference"
   req.Header.Set("vpn", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }
   data1 := string(data)
   if resp.StatusCode != http.StatusOK {
      return "", errors.New(data1)
   }
   return strings.Replace(data1, "/best/", "/ultimate/", 1), nil
}

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

func (a Address) String() string {
   return a[0]
}

// https://www.ctv.ca/shows/friends/the-one-with-the-bullies-s2e21
func (a *Address) Set(data string) error {
   data = strings.TrimPrefix(data, "https://")
   data = strings.TrimPrefix(data, "www.")
   (*a)[0] = strings.TrimPrefix(data, "ctv.ca")
   return nil
}

const query_axis = `
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
` // do not use `query(`

// this is better than strings.Replace and strings.ReplaceAll
func graphql_compact(data string) string {
   return strings.Join(strings.Fields(data), " ")
}

type Address [1]string

type Content struct {
   ContentPackages []struct {
      Id int64
   }
}

func (a *AxisContent) Content() (*Content, error) {
   req, _ := http.NewRequest("", "https://capi.9c9media.com", nil)
   req.URL.Path = func() string {
      b := []byte("/destinations/")
      b = append(b, a.AxisPlaybackLanguages[0].DestinationCode...)
      b = append(b, "/platforms/desktop/contents/"...)
      b = strconv.AppendInt(b, a.AxisId, 10)
      return string(b)
   }()
   req.URL.RawQuery = "$include=[ContentPackages]"
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   content1 := &Content{}
   err = json.NewDecoder(resp.Body).Decode(content1)
   if err != nil {
      return nil, err
   }
   return content1, nil
}

const query_resolve = `
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

func (a Address) Resolve() (*ResolvedPath, error) {
   data, err := json.Marshal(map[string]any{
      "query": graphql_compact(query_resolve),
      "variables": map[string]string{
         "path": a[0],
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
   var value struct {
      Data struct {
         ResolvedPath *struct {
            LastSegment struct {
               Content ResolvedPath
            }
         }
      }
   }
   err = json.Unmarshal(data, &value)
   if err != nil {
      return nil, err
   }
   if value.Data.ResolvedPath == nil {
      return nil, errors.New(string(data))
   }
   return &value.Data.ResolvedPath.LastSegment.Content, nil
}

type ResolvedPath struct {
   FirstPlayableContent *struct {
      Id string
   }
   Id string
}

func (r *ResolvedPath) get_id() string {
   if r.FirstPlayableContent != nil {
      return r.FirstPlayableContent.Id
   }
   return r.Id
}

func (r *ResolvedPath) Axis() (*AxisContent, error) {
   data, err := json.Marshal(map[string]any{
      "query": graphql_compact(query_axis),
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
   var value struct {
      Data struct {
         AxisContent AxisContent
      }
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if len(value.Errors) >= 1 {
      return nil, errors.New(value.Errors[0].Message)
   }
   return &value.Data.AxisContent, nil
}

type AxisContent struct {
   AxisId                int64
   AxisPlaybackLanguages []struct {
      DestinationCode string
   }
}

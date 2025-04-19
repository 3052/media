package itv

import (
   "encoding/json"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

const query = `
query ProgrammePage( $brandLegacyId: BrandLegacyId ) {
   titles(
      filter: { brandLegacyId: $brandLegacyId }
      sortBy: SEQUENCE_ASC
   ) {
      ... on Episode {
         series {
            seriesNumber
         }
         episodeNumber
      }
      title
      latestAvailableVersion {
         playlistUrl
      }
   }
}
`

// this is better than strings.Replace and strings.ReplaceAll
func graphql_compact(data string) string {
   return strings.Join(strings.Fields(data), " ")
}

func programme_page(brand_legacy_id string) ([]title, error) {
   data, err := json.Marshal(map[string]string{
      "brandLegacyId": brand_legacy_id,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "", "https://content-inventory.prd.oasvc.itv.com/discovery", nil,
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = url.Values{
      "query":     {graphql_compact(query)},
      "variables": {string(data)},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Titles []title
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return value.Data.Titles, nil
}

func (t *title) String() string {
   var b []byte
   if t.Series != nil {
      b = []byte("series = ")
      b = strconv.AppendInt(b, t.Series.SeriesNumber, 10)
      b = append(b, "\nepisode = "...)
      b = strconv.AppendInt(b, t.EpisodeNumber, 10)
   }
   if t.Title != "" {
      if b != nil {
         b = append(b, '\n')
      }
      b = append(b, "title = "...)
      b = append(b, t.Title...)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = append(b, "playlist = "...)
   b = append(b, t.LatestAvailableVersion.PlaylistUrl...)
   return string(b)
}

type title struct {
   Series *struct {
      SeriesNumber int64
   }
   EpisodeNumber          int64
   Title                  string
   LatestAvailableVersion struct {
      PlaylistUrl string
   }
}

package max

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "slices"
)

func (v *video) String() string {
   b := fmt.Appendln(nil, "video type =", v.Attributes.VideoType)
   b = fmt.Appendln(b, "season number =", v.Attributes.SeasonNumber)
   b = fmt.Appendln(b, "episode number =", v.Attributes.EpisodeNumber)
   b = fmt.Appendln(b, "name =", v.Attributes.Name)
   b = fmt.Append(b, "edit id = ", v.Relationships.Edit.Data.Id)
   return string(b)
}

func (s season) episode() []video {
   var videos []video
   for _, video1 := range s.Included {
      if video1.Attributes != nil {
         if video1.Attributes.VideoType == "EPISODE" {
            videos = append(videos, video1)
         }
      }
   }
   slices.SortFunc(videos, func(a, b video) int {
      return a.Attributes.EpisodeNumber - b.Attributes.EpisodeNumber
   })
   return videos
}

func (n Login) season() (*season, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.Path = "/cms/collections/generic-show-page-rail-episodes-tabbed-content"
   req.URL.RawQuery = url.Values{
      "include":          {"default"},
      "pf[seasonNumber]": {"1"},
      "pf[show.id]":      {"14f9834d-bc23-41a8-ab61-5c8abdbea505"},
   }.Encode()
   req.Header.Set("authorization", "Bearer "+n.Data.Attributes.Token)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   season1 := &season{}
   err = json.NewDecoder(resp.Body).Decode(season1)
   if err != nil {
      return nil, err
   }
   return season1, nil
}

type season struct {
   Included []video
}

type video struct {
   Attributes *struct {
      VideoType     string
      SeasonNumber  int
      EpisodeNumber int
      Name          string
   }
   Relationships *struct {
      Edit *struct {
         Data struct {
            Id string
         }
      }
   }
}

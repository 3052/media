package max

import (
   "encoding/json"
   "errors"
   "iter"
   "net/http"
   "net/url"
   "path"
   "strconv"
   "strings"
)

func (n Login) Season(id ShowId, number int) (*Videos, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.Path = "/cms/collections/generic-show-page-rail-episodes-tabbed-content"
   req.Header.Set("authorization", "Bearer "+n.Data.Attributes.Token)
   req.URL.RawQuery = url.Values{
      "include":          {"default"},
      "pf[seasonNumber]": {strconv.Itoa(number)},
      "pf[show.id]":      {string(id)},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   season1 := &Videos{}
   err = json.NewDecoder(resp.Body).Decode(season1)
   if err != nil {
      return nil, err
   }
   return season1, nil
}

func (n Login) Movie(id ShowId) (*Videos, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.Path = "/cms/routes/movie/" + string(id)
   req.URL.RawQuery = url.Values{
      "include":          {"default"},
      "page[items.size]": {"1"},
   }.Encode()
   req.Header.Set("authorization", "Bearer "+n.Data.Attributes.Token)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var movie Videos
   err = json.NewDecoder(resp.Body).Decode(&movie)
   if err != nil {
      return nil, err
   }
   if len(movie.Errors) >= 1 {
      return nil, errors.New(movie.Errors[0].Detail)
   }
   return &movie, nil
}

func (s ShowId) String() string {
   return string(s)
}

// max.com/movies/12199308-9afb-460b-9d79-9d54b5d2514c
// max.com/movies/heretic/12199308-9afb-460b-9d79-9d54b5d2514c
// max.com/shows/14f9834d-bc23-41a8-ab61-5c8abdbea505
// max.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505
func (s *ShowId) Set(data string) error {
   switch {
   case strings.Contains(data, "/movies/"):
   case strings.Contains(data, "/shows/"):
   default:
      return errors.New("/movies/ or /shows/ not found")
   }
   *s = ShowId(path.Base(data))
   return nil
}

type ShowId string

type Video struct {
   Attributes *struct {
      SeasonNumber  int
      EpisodeNumber int
      Name          string
      VideoType     string
   }
   Relationships *struct {
      Edit *struct {
         Data struct {
            Id string
         }
      }
   }
}

func (v *Video) String() string {
   var b []byte
   if v.Attributes.SeasonNumber >= 1 {
      b = append(b, "season number = "...)
      b = strconv.AppendInt(b, int64(v.Attributes.SeasonNumber), 10)
   }
   if v.Attributes.EpisodeNumber >= 1 {
      b = append(b, "\nepisode number = "...)
      b = strconv.AppendInt(b, int64(v.Attributes.EpisodeNumber), 10)
   }
   if b != nil {
      b = append(b, '\n')
   }
   b = append(b, "name = "...)
   b = append(b, v.Attributes.Name...)
   b = append(b, "\nvideo type = "...)
   b = append(b, v.Attributes.VideoType...)
   b = append(b, "\nedit id = "...)
   b = append(b, v.Relationships.Edit.Data.Id...)
   return string(b)
}

type Videos struct {
   Errors []struct {
      Detail string
   }
   Included []Video
}

func (v *Videos) Seq() iter.Seq[Video] {
   return func(yield func(Video) bool) {
      for _, video1 := range v.Included {
         if video1.Attributes != nil {
            switch video1.Attributes.VideoType {
            case "EPISODE", "MOVIE":
               if !yield(video1) {
                  return
               }
            }
         }
      }
   }
}

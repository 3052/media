package max

import (
   "encoding/json"
   "errors"
   "fmt"
   "iter"
   "net/http"
   "net/url"
   "slices"
   "strings"
)

func (n Login) season(number int) (*season, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.Path = "/cms/collections/generic-show-page-rail-episodes-tabbed-content"
   req.URL.RawQuery = url.Values{
      "include":          {"default"},
      "pf[show.id]":      {"14f9834d-bc23-41a8-ab61-5c8abdbea505"},
      "pf[seasonNumber]": {fmt.Sprint(number)},
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

func (v *video) String() string {
   b := fmt.Appendln(nil, "season number =", v.Attributes.SeasonNumber)
   b = fmt.Appendln(b, "episode number =", v.Attributes.EpisodeNumber)
   b = fmt.Appendln(b, "name =", v.Attributes.Name)
   b = fmt.Appendln(b, "video type =", v.Attributes.VideoType)
   b = fmt.Append(b, "edit id = ", v.Relationships.Edit.Data.Id)
   return string(b)
}

type video struct {
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

func (s *season) episodes() iter.Seq[video] {
   return func(yield func(video) bool) {
      for _, video1 := range s.Included {
         if video1.Attributes != nil {
            if video1.Attributes.VideoType == "EPISODE" {
               if !yield(video1) {
                  break
               }
            }
         }
      }
   }
}

type season struct {
   Data struct {
      Attributes struct {
         Component struct {
            Filters []struct {
               Options []struct{}
            }
         }
      }
   }
   Included []video
}

func (n Login) seasons() iter.Seq2[*season, error] {
   return func(yield func(*season, error) bool) {
      number := 1
      for {
         season1, err := n.season(number)
         if !yield(season1, err) {
            break
         }
         for _, filter := range season1.Data.Attributes.Component.Filters {
            if number >= len(filter.Options) {
               return
            }
         }
         number++
      }
   }
}

func (s *season) sorted() []video {
   return slices.SortedFunc(s.episodes(), func(a, b video) int {
      return a.Attributes.EpisodeNumber - b.Attributes.EpisodeNumber
   })
}

func (m *movie_item) String() string {
   var b strings.Builder
   b.WriteString("name = ")
   b.WriteString(m.Attributes.Name)
   b.WriteString("\nvideo type = ")
   b.WriteString(m.Attributes.VideoType)
   b.WriteString("\nedit id = ")
   b.WriteString(m.Relationships.Edit.Data.Id)
   return b.String()
}

type movie_item struct {
   Attributes *struct {
      Name string
      VideoType string
   }
   Relationships *struct {
      Edit *struct {
         Data struct {
            Id string
         }
      }
   }
}

func (m *movie_items) movie() (*movie_item, bool) {
   for _, item := range m.Included {
      if item.Attributes != nil {
         if item.Attributes.VideoType == "MOVIE" {
            return &item, true
         }
      }
   }
   return nil, false
}

type movie_items struct {
   Errors []struct {
      Detail string
   }
   Included []movie_item
}

func (n Login) movie(route string) (*movie_items, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.RawQuery = url.Values{
      "include":          {"default"},
      "page[items.size]": {"1"},
   }.Encode()
   req.Header.Set("authorization", "Bearer "+n.Data.Attributes.Token)
   req.URL.Path = "/cms/routes" + route
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var movie movie_items
   err = json.NewDecoder(resp.Body).Decode(&movie)
   if err != nil {
      return nil, err
   }
   if len(movie.Errors) >= 1 {
      return nil, errors.New(movie.Errors[0].Detail)
   }
   return &movie, nil
}

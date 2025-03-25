package max

import (
   "encoding/json"
   "errors"
   "fmt"
   "iter"
   "net/http"
   "net/url"
)

func (n Login) movie(id string) (*videos, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.Path = "/cms/routes/movie/" + id
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
   var movie videos
   err = json.NewDecoder(resp.Body).Decode(&movie)
   if err != nil {
      return nil, err
   }
   if len(movie.Errors) >= 1 {
      return nil, errors.New(movie.Errors[0].Detail)
   }
   return &movie, nil
}

func (n Login) season(show string, number int) (*videos, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.Path = "/cms/collections/generic-show-page-rail-episodes-tabbed-content"
   req.Header.Set("authorization", "Bearer "+n.Data.Attributes.Token)
   req.URL.RawQuery = url.Values{
      "include":          {"default"},
      "pf[show.id]":      {show},
      "pf[seasonNumber]": {fmt.Sprint(number)},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   season1 := &videos{}
   err = json.NewDecoder(resp.Body).Decode(season1)
   if err != nil {
      return nil, err
   }
   return season1, nil
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

func (v *video) String() string {
   b := fmt.Appendln(nil, "season number =", v.Attributes.SeasonNumber)
   b = fmt.Appendln(b, "episode number =", v.Attributes.EpisodeNumber)
   b = fmt.Appendln(b, "name =", v.Attributes.Name)
   b = fmt.Appendln(b, "video type =", v.Attributes.VideoType)
   b = fmt.Append(b, "edit id = ", v.Relationships.Edit.Data.Id)
   return string(b)
}

type videos struct {
   Errors []struct {
      Detail string
   }
   Included []video
}

func (v *videos) seq(video_type string) iter.Seq[video] {
   return func(yield func(video) bool) {
      for _, video1 := range v.Included {
         if video1.Attributes != nil {
            if video1.Attributes.VideoType == video_type {
               if !yield(video1) {
                  break
               }
            }
         }
      }
   }
}

func (v *videos) movie() iter.Seq[video] {
   return v.seq("MOVIE")
}

func (v *videos) episode() iter.Seq[video] {
   return v.seq("EPISODE")
}

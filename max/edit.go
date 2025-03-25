package max

import (
   "encoding/json"
   "errors"
   "iter"
   "net/http"
   "net/url"
   "strconv"
)

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

func (n Login) Movie(id string) (*Videos, error) {
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

func (n Login) Season(show string, number int) (*Videos, error) {
   req, _ := http.NewRequest("", prd_api, nil)
   req.URL.Path = "/cms/collections/generic-show-page-rail-episodes-tabbed-content"
   req.Header.Set("authorization", "Bearer "+n.Data.Attributes.Token)
   req.URL.RawQuery = url.Values{
      "include":          {"default"},
      "pf[show.id]":      {show},
      "pf[seasonNumber]": {strconv.Itoa(number)},
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

func (v *Videos) Movie() iter.Seq[Video] {
   return v.seq("MOVIE")
}

func (v *Videos) Episode() iter.Seq[Video] {
   return v.seq("EPISODE")
}

func (v *Videos) seq(video_type string) iter.Seq[Video] {
   return func(yield func(Video) bool) {
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

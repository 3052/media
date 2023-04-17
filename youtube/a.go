package youtube

import (
   "2a.pages.dev/rosso/http"
   "2a.pages.dev/rosso/json"
   "2a.pages.dev/rosso/protobuf"
   "io"
)

func (r Request) Player(id string, tok *Token) (*Player, error) {
   // fails randomly without this. you can change the inner varint, any value
   // seems to work, examples:
   // 0 CgIQAA==
   // 1 CgIQAQ==
   // 2 CgIQAg==
   // 3 CgIQAw==
   // 4 CgIQBA==
   // 5 CgIQBQ==
   // 6 CgIQBg== github.com/TeamNewPipe/NewPipe/issues/9038
   // 7 CgIQBw==
   // 8 CgIQCA==
   // 9 CgIQCQ==
   r.Params = protobuf.Message{
      1: protobuf.Message{
         2: protobuf.Varint(0),
      },
   }.Marshal()
   r.Video_ID = id
   body, err := json.MarshalIndent(r, "", " ")
   if err != nil {
      return nil, err
   }
   req := http.Post()
   req.Body_Bytes(body)
   if tok != nil {
      req.Header.Set("Authorization", "Bearer " + tok.Access_Token)
   } else {
      req.Header.Set("X-Goog-API-Key", api_key)
   }
   req.URL.Host = "www.youtube.com"
   req.URL.Path = "/youtubei/v1/player"
   req.URL.Scheme = "https"
   res, err := http.Default_Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   play := new(Player)
   if err := json.NewDecoder(res.Body).Decode(play); err != nil {
      return nil, err
   }
   return play, nil
}

func Android() Request {
   var r Request
   r.Content_Check_OK = true
   r.Context.Client.Name = "ANDROID"
   r.Context.Client.Version = android_version
   return r
}

func Android_Check() Request {
   var r Request
   r.Content_Check_OK = true
   r.Context.Client.Name = "ANDROID"
   r.Context.Client.Version = android_version
   r.Racy_Check_OK = true
   return r
}

func Mobile_Web() Request {
   var r Request
   r.Context.Client.Name = "MWEB"
   r.Context.Client.Version = mweb_version
   return r
}

func Android_Embed() Request {
   var r Request
   r.Context.Client.Name = "ANDROID_EMBEDDED_PLAYER"
   r.Context.Client.Version = android_version
   return r
}

type Request struct {
   Content_Check_OK bool `json:"contentCheckOk,omitempty"`
   Context struct {
      Client struct {
         Name string `json:"clientName"`
         Version string `json:"clientVersion"`
      } `json:"client"`
   } `json:"context"`
   Params []byte `json:"params,omitempty"`
   Query string `json:"query,omitempty"`
   Racy_Check_OK bool `json:"racyCheckOk,omitempty"`
   Video_ID string `json:"videoId,omitempty"`
}
const (
   // com.google.android.youtube
   // all versions should be valid starting with 16
   android_version = "18.14.40"
   api_key = "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8"
   mweb_version = "2.20230405.01.00"
)

type config struct {
   Innertube_API_Key string
   Innertube_Client_Name string
   Innertube_Client_Version string
}

func new_config() (*config, error) {
   req := http.Get()
   req.URL.Scheme = "https"
   req.URL.Host = "m.youtube.com"
   req.Header.Set("User-Agent", "iPad")
   res, err := http.Default_Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   data, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   sep := []byte("\nytcfg.set(")
   con := new(config)
   if err := json.Cut(data, sep, con); err != nil {
      return nil, err
   }
   return con, nil
}

func (r Request) Search(query string) (*Search, error) {
   param := New_Params()
   param.Type(Type["Video"])
   r.Params = param.Marshal()
   r.Query = query
   body, err := json.MarshalIndent(r, "", " ")
   if err != nil {
      return nil, err
   }
   req := http.Post()
   req.Body_Bytes(body)
   req.Header.Set("X-Goog-API-Key", api_key)
   req.URL.Host = "www.youtube.com"
   req.URL.Path = "/youtubei/v1/search"
   req.URL.Scheme = "https"
   res, err := http.Default_Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   search := new(Search)
   if err := json.NewDecoder(res.Body).Decode(search); err != nil {
      return nil, err
   }
   return search, nil
}
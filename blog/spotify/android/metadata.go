package android

import (
   "154.pages.dev/protobuf"
   "bytes"
   "encoding/hex"
   "errors"
   "io"
   "net/http"
)

type file_format struct {
   m protobuf.Message
}

func (m metadata) file() chan file_format {
   vs := make(chan file_format)
   go func() {
      for v := range m.m.Get(2) {
         v = <-v.Get(3)
         v = <-v.Get(3)
         v = <-v.Get(2)
         for v := range v.Get(12) {
            vs <- file_format{v}
         }
      }
      close(vs)
   }()
   return vs
}

// github.com/librespot-org/librespot/blob/dev/protocol/proto/media_format.proto
// github.com/librespot-org/librespot/blob/dev/protocol/proto/metadata.proto
// github.com/librespot-org/librespot/blob/dev/protocol/proto/media_manifest.proto
type metadata struct {
   m protobuf.Message
}

func (o LoginOk) metadata(canonical_uri string) (*metadata, error) {
   token, err := o.AccessToken()
   if err != nil {
      return nil, err
   }
   var m protobuf.Message
   m.Add(2, func(m *protobuf.Message) {
      m.AddBytes(1, []byte(canonical_uri))
      m.Add(2, func(m *protobuf.Message) {
         m.AddVarint(1, 10)
      })
   })
   req, err := http.NewRequest(
      "POST", "https://guc3-spclient.spotify.com", bytes.NewReader(m.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/extended-metadata/v0/extended-metadata"
   req.Header.Set("authorization", "Bearer "+token)
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      var b bytes.Buffer
      res.Write(&b)
      return nil, errors.New(b.String())
   }
   data, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   var meta metadata
   if err := meta.m.Consume(data); err != nil {
      return nil, err
   }
   return &meta, nil
}

func (f file_format) file_id() (string, bool) {
   if v, ok := <-f.m.GetBytes(1); ok {
      return hex.EncodeToString(v), true
   }
   return "", false
}

func (f file_format) OGG_VORBIS_320() bool {
   if v, ok := f.format(); ok {
      if v == 2 {
         return true
      }
   }
   return false
}

func (f file_format) format() (uint64, bool) {
   if v, ok := <-f.m.GetVarint(2); ok {
      return uint64(v), true
   }
   return 0, false
}
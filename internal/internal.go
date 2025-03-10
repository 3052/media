package internal

import (
   "41.neocities.org/dash"
   "41.neocities.org/sofia/file"
   "41.neocities.org/sofia/pssh"
   "41.neocities.org/widevine"
   "41.neocities.org/x/progress"
   "bufio"
   "bytes"
   "encoding/base64"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "os"
   "slices"
   "strings"
)

func (e *License) get_key(message *pssh_data) ([]byte, error) {
   if message.key_id == nil {
      return nil, nil
   }
   private_key, err := os.ReadFile(e.PrivateKey)
   if err != nil {
      return nil, err
   }
   client_id, err := os.ReadFile(e.ClientId)
   if err != nil {
      return nil, err
   }
   if message.pssh == nil {
      var pssh1 widevine.Pssh
      pssh1.KeyIds = [][]byte{message.key_id}
      message.pssh = pssh1.Marshal()
   }
   log.Println("PSSH", base64.StdEncoding.EncodeToString(message.pssh))
   var module widevine.Cdm
   err = module.New(private_key, client_id, message.pssh)
   if err != nil {
      return nil, err
   }
   data, err := module.RequestBody()
   if err != nil {
      return nil, err
   }
   data, err = e.Widevine(data)
   if err != nil {
      return nil, err
   }
   var body widevine.ResponseBody
   err = body.Unmarshal(data)
   if err != nil {
      return nil, err
   }
   block, err := module.Block(body)
   if err != nil {
      return nil, err
   }
   for container := range body.Container() {
      if bytes.Equal(container.Id(), message.key_id) {
         key := container.Key(block)
         log.Println("key", base64.StdEncoding.EncodeToString(key))
         return key, nil
      }
   }
   return nil, errors.New("get_key")
}

func (e *License) segment_template(represent *dash.Representation) error {
   var message pssh_data
   err := message.New(represent)
   if err != nil {
      return err
   }
   os_file, err := dash_create(represent)
   if err != nil {
      return err
   }
   defer os_file.Close()
   if initial := represent.SegmentTemplate.Initialization; initial != "" {
      initial2, err := initial.Url(represent)
      if err != nil {
         return err
      }
      data, err := get(initial2, nil)
      if err != nil {
         return err
      }
      data, err = message.initialization(data)
      if err != nil {
         return err
      }
      _, err = os_file.Write(data)
      if err != nil {
         return err
      }
   }
   key, err := e.get_key(&message)
   if err != nil {
      return err
   }
   var segments []int
   for represent1 := range represent.Representation() {
      segments = slices.AppendSeq(segments, represent1.Segment())
   }
   var parts progress.Parts
   parts.Set(len(segments))
   head := http.Header{}
   head.Set("silent", "true")
   for _, segment := range segments {
      media, err := represent.SegmentTemplate.Media.Url(represent, segment)
      if err != nil {
         return err
      }
      data, err := get(media, head)
      if err != nil {
         return err
      }
      parts.Next()
      data, err = write_segment(data, key)
      if err != nil {
         return err
      }
      _, err = os_file.Write(data)
      if err != nil {
         return err
      }
   }
   return nil
}

func (e *License) segment_base(represent *dash.Representation) error {
   var message pssh_data
   err := message.New(represent)
   if err != nil {
      return err
   }
   os_file, err := dash_create(represent)
   if err != nil {
      return err
   }
   defer os_file.Close()
   base := represent.SegmentBase
   data, err := get(represent.BaseUrl[0], http.Header{
      "range": {"bytes=" + base.Initialization.Range.String()},
   })
   if err != nil {
      return err
   }
   data, err = message.initialization(data)
   if err != nil {
      return err
   }
   _, err = os_file.Write(data)
   if err != nil {
      return err
   }
   key, err := e.get_key(&message)
   if err != nil {
      return err
   }
   data, err = get(represent.BaseUrl[0], http.Header{
      "range": {"bytes=" + base.IndexRange.String()},
   })
   if err != nil {
      return err
   }
   var file_file file.File
   err = file_file.Read(data)
   if err != nil {
      return err
   }
   var parts progress.Parts
   parts.Set(len(file_file.Sidx.Reference))
   head := http.Header{}
   head.Set("silent", "true")
   for _, reference := range file_file.Sidx.Reference {
      base.IndexRange[0] = base.IndexRange[1] + 1
      base.IndexRange[1] += uint64(reference.Size())
      head.Set("range", "bytes=" + base.IndexRange.String())
      data, err = get(represent.BaseUrl[0], head)
      if err != nil {
         return err
      }
      parts.Next()
      data, err = write_segment(data, key)
      if err != nil {
         return err
      }
      _, err = os_file.Write(data)
      if err != nil {
         return err
      }
   }
   return nil
}

func (e *License) segment_list(represent *dash.Representation) error {
   var message pssh_data
   err := message.New(represent)
   if err != nil {
      return err
   }
   os_file, err := dash_create(represent)
   if err != nil {
      return err
   }
   defer os_file.Close()
   initial, err := represent.SegmentList.Initialization.SourceUrl.Url(represent)
   if err != nil {
      return err
   }
   data, err := get(initial, nil)
   if err != nil {
      return err
   }
   data, err = message.initialization(data)
   if err != nil {
      return err
   }
   _, err = os_file.Write(data)
   if err != nil {
      return err
   }
   key, err := e.get_key(&message)
   if err != nil {
      return err
   }
   var parts progress.Parts
   parts.Set(len(represent.SegmentList.SegmentUrl))
   head := http.Header{}
   head.Set("silent", "true")
   for _, segment := range represent.SegmentList.SegmentUrl {
      media, err := segment.Media.Url(represent)
      if err != nil {
         return err
      }
      data, err := get(media, head)
      if err != nil {
         return err
      }
      parts.Next()
      data, err = write_segment(data, key)
      if err != nil {
         return err
      }
      _, err = os_file.Write(data)
      if err != nil {
         return err
      }
   }
   return nil
}

type pssh_data struct {
   key_id []byte
   pssh   []byte
}

func (p *pssh_data) New(represent *dash.Representation) error {
   for _, content := range represent.ContentProtection {
      if content.SchemeIdUri == widevine_urn {
         if content.Pssh != "" {
            data, err := base64.StdEncoding.DecodeString(content.Pssh)
            if err != nil {
               return err
            }
            var box pssh.Box
            n, err := box.BoxHeader.Decode(data)
            if err != nil {
               return err
            }
            err = box.Read(data[n:])
            if err != nil {
               return err
            }
            p.pssh = box.Data
            break
         }
      }
   }
   return nil
}

func (p *pssh_data) initialization(data []byte) ([]byte, error) {
   var file_file file.File
   err := file_file.Read(data)
   if err != nil {
      return nil, err
   }
   if moov, ok := file_file.GetMoov(); ok {
      for _, pssh1 := range moov.Pssh {
         if pssh1.SystemId.String() == widevine_system_id {
            p.pssh = pssh1.Data
         }
         copy(pssh1.BoxHeader.Type[:], "free") // Firefox
      }
      description := moov.Trak.Mdia.Minf.Stbl.Stsd
      if sinf, ok := description.Sinf(); ok {
         p.key_id = sinf.Schi.Tenc.DefaultKid[:]
         // Firefox
         copy(sinf.BoxHeader.Type[:], "free")
         if sample, ok := description.SampleEntry(); ok {
            // Firefox
            copy(sample.BoxHeader.Type[:], sinf.Frma.DataFormat[:])
         }
      }
   }
   return file_file.Append(nil)
}

const (
   widevine_system_id = "edef8ba979d64acea3c827dcd51d21ed"
   widevine_urn       = "urn:uuid:edef8ba9-79d6-4ace-a3c8-27dcd51d21ed"
)

var Forward = []struct {
   Country string
   Ip      string
}{
   {"Argentina", "186.128.0.0"},
   {"Australia", "1.128.0.0"},
   {"Bolivia", "179.58.0.0"},
   {"Brazil", "179.192.0.0"},
   {"Canada", "99.224.0.0"},
   {"Chile", "191.112.0.0"},
   {"Colombia", "181.128.0.0"},
   {"Costa Rica", "201.192.0.0"},
   {"Denmark", "2.104.0.0"},
   {"Ecuador", "186.68.0.0"},
   {"Egypt", "197.32.0.0"},
   {"Germany", "53.0.0.0"},
   {"Guatemala", "190.56.0.0"},
   {"India", "106.192.0.0"},
   {"Indonesia", "39.192.0.0"},
   {"Ireland", "87.32.0.0"},
   {"Italy", "79.0.0.0"},
   {"Latvia", "78.84.0.0"},
   {"Malaysia", "175.136.0.0"},
   {"Mexico", "189.128.0.0"},
   {"Netherlands", "145.160.0.0"},
   {"New Zealand", "49.224.0.0"},
   {"Norway", "88.88.0.0"},
   {"Peru", "190.232.0.0"},
   {"Russia", "95.24.0.0"},
   {"South Africa", "105.0.0.0"},
   {"South Korea", "175.192.0.0"},
   {"Spain", "88.0.0.0"},
   {"Sweden", "78.64.0.0"},
   {"Taiwan", "120.96.0.0"},
   {"United Kingdom", "25.0.0.0"},
   {"Venezuela", "190.72.0.0"},
}

func write_segment(data, key []byte) ([]byte, error) {
   if key == nil {
      return data, nil
   }
   var file_file file.File
   err := file_file.Read(data)
   if err != nil {
      return nil, err
   }
   track := file_file.Moof.Traf
   if senc := track.Senc; senc != nil {
      for i, data := range file_file.Mdat.Data(&track) {
         err = senc.Sample[i].DecryptCenc(data, key)
         if err != nil {
            return nil, err
         }
      }
   }
   return file_file.Append(nil)
}

func create(name string) (*os.File, error) {
   log.Println("Create", name)
   return os.Create(name)
}

func dash_create(represent *dash.Representation) (*os.File, error) {
   switch *represent.MimeType {
   case "audio/mp4":
      return create(".m4a")
   case "text/vtt":
      return create(".vtt")
   case "video/mp4":
      return create(".m4v")
   }
   return nil, errors.New(*represent.MimeType)
}

func init() {
   log.SetFlags(log.Ltime)
   http.DefaultClient.Transport = &transport{
      // github.com/golang/go/issues/18639
      Protocols: &http.Protocols{},
      Proxy: http.ProxyFromEnvironment,
   }
}

type transport http.Transport

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
   if req.Header.Get("silent") == "" {
      log.Println(req.Method, req.URL)
   }
   return (*http.Transport)(t).RoundTrip(req)
}

type License struct {
   ClientId string
   PrivateKey string
   Widevine func([]byte) ([]byte, error)
}

func get(u *url.URL, head http.Header) ([]byte, error) {
   req := http.Request{Method: "GET", URL: u}
   if head != nil {
      req.Header = head
   } else {
      req.Header = http.Header{}
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   switch resp.StatusCode {
   case http.StatusOK, http.StatusPartialContent:
   default:
      var data strings.Builder
      resp.Write(&data)
      return nil, errors.New(data.String())
   }
   return io.ReadAll(resp.Body)
}
func (e *License) Download(name, id string) error {
   data, err := os.ReadFile(name)
   if err != nil {
      return err
   }
   resp, err := unmarshal(data)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   var media dash.Mpd
   err = media.Unmarshal(data)
   if err != nil {
      return err
   }
   media.Set(resp.Request.URL)
   http.DefaultClient.Jar, err = cookiejar.New(nil)
   if err != nil {
      return err
   }
   http.DefaultClient.Jar.SetCookies(resp.Request.URL, resp.Cookies())
   for represent := range media.Representation() {
      if represent.Id == id {
         if represent.SegmentBase != nil {
            return e.segment_base(&represent)
         }
         if represent.SegmentList != nil {
            return e.segment_list(&represent)
         }
         return e.segment_template(&represent)
      }
   }
   return nil
}

func Mpd(name string, resp *http.Response) error {
   data, err := marshal(resp)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   err = os.WriteFile(name, data, os.ModePerm)
   if err != nil {
      return err
   }
   resp, err = unmarshal(data)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   var media dash.Mpd
   err = media.Unmarshal(data)
   if err != nil {
      return err
   }
   media.Set(resp.Request.URL)
   represents := slices.SortedFunc(media.Representation(),
      func(a, b dash.Representation) int {
         return a.Bandwidth - b.Bandwidth
      },
   )
   for i, represent := range represents {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&represent)
   }
   return nil
}

func marshal(resp *http.Response) ([]byte, error) {
   var buf bytes.Buffer
   _, err := fmt.Fprintln(&buf, resp.Request.URL)
   if err != nil {
      return nil, err
   }
   err = resp.Write(&buf)
   if err != nil {
      return nil, err
   }
   return buf.Bytes(), nil
}

func unmarshal(data []byte) (*http.Response, error) {
   data1, data, _ := bytes.Cut(data, []byte{'\n'})
   var base url.URL
   err := base.UnmarshalBinary(data1)
   if err != nil {
      return nil, err
   }
   return http.ReadResponse(
      bufio.NewReader(bytes.NewReader(data)), &http.Request{URL: &base},
   )
}

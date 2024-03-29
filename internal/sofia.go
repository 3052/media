package internal

import (
   "154.pages.dev/encoding/dash"
   "154.pages.dev/sofia"
   "encoding/base64"
   "errors"
   "io"
   "log/slog"
   "net/http"
)

func write_init(w io.Writer, r io.Reader) ([]byte, error) {
   var file sofia.File
   if err := file.Read(r); err != nil {
      return nil, err
   }
   for _, protect := range file.Movie.Protection {
      copy(protect.BoxHeader.Type[:], "free") // Firefox
   }
   description := &file.Movie.Track.Media.
      MediaInformation.SampleTable.SampleDescription
   if v := description.AudioSample; v != nil {
      copy(v.ProtectionScheme.BoxHeader.Type[:], "free") // Firefox
      copy(
         v.SampleEntry.BoxHeader.Type[:],
         v.ProtectionScheme.OriginalFormat.DataFormat[:],
      ) // Firefox
   }
   if v := description.VisualSample; v != nil {
      copy(v.ProtectionScheme.BoxHeader.Type[:], "free") // Firefox
      copy(
         v.SampleEntry.BoxHeader.Type[:],
         v.ProtectionScheme.OriginalFormat.DataFormat[:],
      ) // Firefox
   }
   if err := file.Write(w); err != nil {
      return nil, err
   }
   pssh, ok := file.Movie.Widevine()
   if !ok {
      return nil, errors.New("sofia.Movie.Widevine")
   }
   slog.Debug("Widevine", "PSSH", base64.StdEncoding.EncodeToString(pssh))
   return pssh, nil
}

func write_segment(w io.Writer, r io.Reader, key []byte) error {
   var file sofia.File
   if err := file.Read(r); err != nil {
      return err
   }
   fragment := file.MovieFragment.TrackFragment
   for i, data := range file.MediaData.Data(fragment.TrackRun) {
      err := fragment.SampleEncryption.Samples[i].DecryptCenc(data, key)
      if err != nil {
         return err
      }
   }
   return file.Write(w)
}

func write_sidx(base_url string, raw dash.RawRange) ([]sofia.Range, error) {
   req, err := http.NewRequest("GET", base_url, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("Range", "bytes=" + string(raw))
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   var file sofia.File
   if err := file.Read(res.Body); err != nil {
      return nil, err
   }
   index, err := raw.Scan()
   if err != nil {
      return nil, err
   }
   return file.SegmentIndex.Ranges(index.End+1), nil
}

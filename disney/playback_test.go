package disney

import "testing"

func Test(t *testing.T) {
   t.Log(tests)
}

var tests = []struct {
   location string
   playback_id string
   url      string
}{
   {
      location: "US",
      url:      "https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
   },
   {
      location: "KR", // MUST DO KR LOGIN FIRST
      url:      "https://disneyplus.com/browse/entity-d0d0796c-a144-42fa-a730-4cbd1014ef1f",
   },
}

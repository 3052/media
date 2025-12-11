package pluto

import (
   "strings"
   "testing"
   "time"
)

/*
widevine:L1   androidmobile    height="1080"   false
widevine:L1   androidmobile    height="576"   true
widevine:L1   androidmobile    height="720"   true
widevine:L1   androidtv    height="1080"   true
widevine:L1   androidtv    height="576"   true
widevine:L1   androidtv    height="720"   true
widevine:L1   web    height="1080"   false
widevine:L1   web    height="576"   true
widevine:L1   web    height="720"   true
widevine:L3   androidmobile    height="1080"   false
widevine:L3   androidmobile    height="576"   true
widevine:L3   androidmobile    height="720"   false
widevine:L3   androidtv    height="1080"   false
widevine:L3   androidtv    height="576"   true
widevine:L3   androidtv    height="720"   false
widevine:L3   web    height="1080"   false
widevine:L3   web    height="576"   true
widevine:L3   web    height="720"   false
*/
func TestFetch(t *testing.T) {
   var capabilities = []string{
      "widevine:L1",
      "widevine:L3",
   }
   var names = []string{
      "androidmobile",
      "androidtv",
      "web",
   }
   var heights = []string{
      ` height="1080"`,
      ` height="576"`,
      ` height="720"`,
   }
   for _, capability := range capabilities {
      for _, name := range names {
         app_name = name
         drm_capabilities = capability
         var series_var Series
         err := series_var.Fetch("6495eff09263a40013cf63a5")
         if err != nil {
            t.Fatal(err)
         }
         _, data, err := series_var.Mpd()
         if err != nil {
            t.Fatal(err)
         }
         for _, height := range heights {
            t.Log(
               capability, " ",
               name, " ",
               height, " ",
               strings.Contains(string(data), height),
            )
         }
         time.Sleep(time.Second)
      }
   }
}

var tests = []string{
   "https://pluto.tv/on-demand/movies/6495eff09263a40013cf63a5",
   "https://pluto.tv/on-demand/series/66d0bb64a1c89200137fb0e6",
}

func TestLog(t *testing.T) {
   t.Log(tests)
}

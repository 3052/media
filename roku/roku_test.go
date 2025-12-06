package roku

import "testing"

var tests = []struct {
   id         string
   type_value string
   url        string
}{
   {
      id:         "597a64a4a25c5bf6af4a8c7053049a6f",
      type_value: "movie",
      url:        "https://therokuchannel.roku.com/watch/597a64a4a25c5bf6af4a8c7053049a6f",
   },
   {
      id:         "105c41ea75775968b670fbb26978ed76",
      type_value: "episode",
      url:        "https://therokuchannel.roku.com/watch/105c41ea75775968b670fbb26978ed76",
   },
}

func Test(t *testing.T) {
   t.Log(tests)
}

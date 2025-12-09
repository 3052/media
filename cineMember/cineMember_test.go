package cineMember

import "testing"

var tests = []string{
   "https://cinemember.nl/nl/title/468845/the-worst-person-in-the-world",
   "https://cinemember.nl/nl/title/469991/knives-out", // buffer too small
}

func Test(t *testing.T) {
   t.Log(tests)
}

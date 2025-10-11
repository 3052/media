package cineMember

import (
   "fmt"
   "testing"
)

var tests = []string{
   "https://cinemember.nl/nl/title/468545/american-hustle",
   "https://cinemember.nl/nl/title/469991/knives-out", // buffer too small
}

func Test(t *testing.T) {
   fmt.Println(tests)
}

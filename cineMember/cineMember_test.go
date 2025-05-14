package cineMember

import (
   "fmt"
   "testing"
)

var tests = []string{
   "https://cinemember.nl/films/american-hustle",
   "https://cinemember.nl/films/knives-out", // buffer too small
}

func Test(t *testing.T) {
   fmt.Println(tests)
}

package cineMember

import "testing"

var tests = []string{
   "https://cinemember.nl/films/american-hustle",
   "https://cinemember.nl/films/knives-out",
}

func Test(t *testing.T) {
   for _, test1 := range tests {
      var web Address
      err := web.Set(test1)
      if err != nil {
         t.Fatal(err)
      }
      _, err = web.Article()
      if err != nil {
         t.Fatal(err)
      }
   }
}

package member

import (
   "fmt"
   "os"
   "testing"
)

func TestAuthenticate(t *testing.T) {
   username := os.Getenv("cine_member_username")
   if username == "" {
      t.Fatal("Getenv")
   }
   password := os.Getenv("cine_member_password")
   text, err := NewAuthenticate(username, password)
   if err != nil {
      t.Fatal(err)
   }
   os.WriteFile("authenticate.json", text, 0666)
}

func TestAsset(t *testing.T) {
   article, err := american_hustle.Article()
   if err != nil {
      t.Fatal(err)
   }
   text, err := os.ReadFile("authenticate.json")
   if err != nil {
      t.Fatal(err)
   }
   var auth Authenticate
   err = auth.Unmarshal(text)
   if err != nil {
      t.Fatal(err)
   }
   asset, ok := article.Film()
   if !ok {
      t.Fatal(ArticleAsset{})
   }
   play, err := auth.Play(asset)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(play.Dash())
}

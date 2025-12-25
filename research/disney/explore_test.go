package disney

import "testing"

func TestExplore(t *testing.T) {
   email, err := output("credential", "-h=disneyplus.com", "-k=user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h=disneyplus.com")
   if err != nil {
      t.Fatal(err)
   }
   token, err := register_device()
   if err != nil {
      t.Fatal(err)
   }
   account, err := token.login(email, password)
   if err != nil {
      t.Fatal(err)
   }
   _, err = account.explore(test.entity)
   if err != nil {
      t.Fatal(err)
   }
}

var test = struct {
   entity string
   url    string
}{
   entity: "7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
   url:    "https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
}

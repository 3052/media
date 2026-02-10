package sky

import (
   "encoding/xml"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

type Cookie struct {
   Cookie *http.Cookie
}

func (c *Cookie) Set(data string) error {
   var err error
   c.Cookie, err = http.ParseSetCookie(data)
   if err != nil {
      return err
   }
   return nil
}

func (c Cookie) String() string {
   return c.Cookie.String()
}

func (c cookies) session_id() (*Cookie, bool) {
   for _, cookie1 := range c {
      if cookie1.Name == "_ASP.NET_SessionId_" {
         return &Cookie{cookie1}, true
      }
   }
   return nil, false
}

const (
   out_of_country = "/out-of-country"
   verification_token = "__RequestVerificationToken"
)

// hard geo block
func (n *login) New() error {
   req, _ := http.NewRequest("", "https://show.sky.ch/de/login", nil)
   req.Header.Set("tv", "Emulator")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return errors.New(resp.Status)
   }
   if strings.HasSuffix(resp.Request.URL.Path, out_of_country) {
      return errors.New(out_of_country)
   }
   err = xml.NewDecoder(resp.Body).Decode(&n.section)
   if err != nil {
      return err
   }
   n.cookies = resp.Cookies()
   return nil
}

type login struct {
   cookies []*http.Cookie
   section struct {
      Div     struct {
         Form struct {
            Input []struct {
               Name  string `xml:"name,attr"`
               Value string `xml:"value,attr"`
            } `xml:"input"`
         } `xml:"form"`
      } `xml:"div"`
   }
}

func (n *login) cookie_token() (*http.Cookie, error) {
   for _, cookie1 := range n.cookies {
      if cookie1.Name == verification_token {
         return cookie1, nil
      }
   }
   return nil, http.ErrNoCookie
}

func (n *login) input_token() (string, error) {
   for _, input := range n.section.Div.Form.Input {
      if input.Name == verification_token {
         return input.Value, nil
      }
   }
   return "", errors.New(verification_token)
}

// hard geo block
func (n *login) login(username, password string) (cookies, error) {
   input_token, err := n.input_token()
   if err != nil {
      return nil, err
   }
   data := url.Values{
      "password": {password},
      "username": {username},
      verification_token: {input_token},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://show.sky.ch/de/Authentication/Login",
      strings.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Header.Set("tv", "Emulator")
   cookie_token, err := n.cookie_token()
   if err != nil {
      return nil, err
   }
   req.AddCookie(cookie_token)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   if err != nil {
      return nil, err
   }
   return resp.Cookies(), nil
}

type cookies []*http.Cookie

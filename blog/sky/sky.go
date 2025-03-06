package sky

import (
   "encoding/xml"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// x-forwarded-for fail
// mullvad.net fail
// smartproxy.com fail
// proxy-seller.com pass
func sky_player(cookie *http.Cookie) ([]byte, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "show.sky.ch"
   req.URL.Path = "/de/SkyPlayerAjax/SkyPlayer"
   req.URL.Scheme = "https"
   values := url.Values{}
   values["id"] = []string{"2035"}
   values["contentType"] = []string{"2"}
   req.URL.RawQuery = values.Encode()
   req.Header["X-Requested-With"] = []string{"XMLHttpRequest"}
   req.AddCookie(cookie)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if strings.Contains(string(data), not_available[0]) {
      return nil, not_available
   }
   return data, nil
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

type cookies []*http.Cookie

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

func (c *cookie) Set(data string) error {
   var err error
   (*c)[0], err = http.ParseSetCookie(data)
   if err != nil {
      return err
   }
   return nil
}

type cookie [1]*http.Cookie
func (c cookie) String() string {
   return c[0].String()
}

func (c cookies) session_id() (cookie, bool) {
   for _, cookie1 := range c {
      if cookie1.Name == "_ASP.NET_SessionId_" {
         return cookie{cookie1}, true
      }
   }
   return cookie{}, false
}

var not_available = service{
   "We're sorry our service is not available in your region yet.",
}

type service [1]string

func (s service) Error() string {
   s[0] = strings.TrimSuffix(s[0], ".")
   return strings.ToLower(s[0])
}

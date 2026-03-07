package disney

import (
   "fmt"
   "testing"
)

func TestAuthenticateWithOtp(t *testing.T) {
   var device_item Device
   device_item.Token.AccessToken = otp_test.access_token
   authenticate, err := device_item.AuthenticateWithOtp(
      otp_test.email, otp_test.passcode,
   )
   if err != nil {
      t.Fatal(err)
   }
   inactive, err := device_item.LoginWithActionGrant(authenticate.ActionGrant)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", inactive)
}

func TestRequestOtp(t *testing.T) {
   var device_item Device
   device_item.Token.AccessToken = otp_test.access_token
   otp, err := device_item.RequestOtp(otp_test.email)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(otp)
}

func TestDevice(t *testing.T) {
   device_item, err := RegisterDevice()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(device_item.Token.AccessToken)
}

func TestEntity(t *testing.T) {
   t.Log(entity_tests)
}

var entity_tests = []struct {
   entity string
   format string
   url    string
}{
   {
      entity: "movie",
      format: "4K ULTRA HD",
      url:    "https://disneyplus.com/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
   },
   {
      entity: "movie",
      format: "4K ULTRA HD",
      url:    "https://disneyplus.com/browse/entity-917f1bf3-3db4-4df0-afe2-60b2c5e67618",
   },
   {
      entity: "series",
      format: "HD",
      url:    "https://disneyplus.com/browse/entity-21e70fbf-6a51-41b3-88e9-f111830b046c",
   },
}

var otp_test = struct {
   email string
   passcode string
   access_token string
}{
   email: "27@riseup.net",
   passcode: "177224",
   access_token: "eyJ6aXAiOiJERUYiLCJraWQiOiJ0Vy10M2ZQUTJEN2Q0YlBWTU1rSkd4dkJlZ0ZXQkdXek5KcFFtOGRJMWYwIiwiY3R5IjoiSldUIiwiZW5jIjoiQzIwUCIsImFsZyI6ImRpciJ9..xHCIfoGXxlc5-0V9.PF4O6PDDE-FadTNm2PWIgeFgpvbXSPtc4QfjpuYbGzCXIwXKNoXnt5EaA2rIsTOTYl__Go33gkGX2TSan9sXqwme8pq5Nb5A61XSfwTSa13XWPN0e-bOQf18B0rFA48iVR1hWhhe4mi9d94Ygsv1kK8xRXvGViKXcKX30On8GJI0PHx0zL73BHcnWyYgn7628G3plPW51iqG_i_XhhYenLk8_MlbzKBglmWfcCyxlUVlzpU68EUtHUmx6XrRb0aSELRhJETIJ28jrAKni5pO_deCwFiKq1zfYv-DF3sJKaJJQcsGlnHa4IFK44-U_mLkkX-0Op9MEQw9cwTiwpSQPxMYu7fnugaLDJcXAwaI8CYJ39226yfXlqkINqnNmRCeei1cz9bziUN-uop0XUkWAFpRahqocZvwlrIhFVBAOW703KQSOLEfWh_GElqWaXkgW_Tm8pNRTQFyOC1Z20-MsZNAatx8bu6ZB7RJCO5RgvKRwz-8MaGJZXDCcob6tl36RaxAmAtkl_UQgvdvFhoQX8t9eitF0sjS2d5Aoy-cfuXLO3kF0pAzgFVGqCeQ_n78Oa4YuGI6RjwSt5TRBc2fg8dXUF0phUKyfYy_dU-CQ3Nw-7ci2Na0GKVUa-cLWCNGGrnEJiJ1S2SQ39VL5XM1b7UN09l_PxYBwQ8jUEiIUJ-4fZJJGVJduQ4deEXw6OsVzV0yQcN4dt1LW5u1rO98Ba_sr7ygXV2xYRGNicCLnHRHlKC7-F_mFS8aD7tZuvnYrtxPWS8ODbneiDNKwwZYLdMaZYxY2Wr1z5sotfihOpfej-nUpeF98z0NLamXkVWTmUs0bwdYdWSPR0Cto7tOLvvcQ2FctSwotk9If_Eb_A9qN91JtJQJz4dM0TyfsCxfa0QxsQkz5of-UmfnATLuqetJrOD35dIVu3dW_-qE6i2lnJZ9VlWPhOD3BgHo3z_cp2BeJEn0j4nogJWKypTKMSQqIw3cejyx5O2ioj_LVWhxBVCX6z_qsmK-CyNQweR_ox7VIXEapGYq5g3BPbEvKYql61v8nEVuroN2AEcPwud_YtFfoJlFFfJDJ6nXuS6_Zfd3mtpPTcodW50k4ddSbcKu5VUvtv5iqLSgqNEj2y0TPnq0_bE5zEf8WqB8I9T7MBc_rvfpMQ-aizqMHWzNEUNAvO_HzpaL3OigucZwuOe5LHrUKIBjJcI83yYSKNV3_vZsccQ-WzyziJuQHG_iC7JvhwGzgEXJa42xC1t5dfndjjI4IGpMZgJtJR9IjhoSobYjIvqejCEtUuPOmETuyctuEOGKeTtUDLXLTPYmrvDUcj1dwQCbawQ-mLIJSq5IDAS8RsmIZviYtYLYvEkAUHE9Wf_hIzPxJZXoM8kRgwHVTY1mNyjqTjcEppOd6yNOoOWI9Wvj1DDEQ1RF17giCCGgZlGt59jhQdDB-oM3D6UhCEttCAxOMXH8li0QBlW0vtUfe1-f54d5TvDYLZMELU0WLOyjj4Q-Dde-8p5zDSqYFjcBJTwbBUWh5crPsS3BbJQVOKDL04QhoiLvd1JtpkQYDA9bPDcgkwR2ocEdlgkJEQmie1IkllIYlqQoUrRsKRhNSmyDk_JbHV3ZIrpZPbOpJl3YplMyKHKk-U0xpgrjjjauvE5ScvfK8iBofdUP7yKMqE8KBQ8UMi4GKaVJMZbDM_6RIr1fD-OvkNdDJIHegx8sfI4_44KCAJR_RAsxZlpgt-A44jtXGxtlYJYdmLyGzEPZn6B3BGAUzfCjkY5qYRWfRJtyIAzRtsyBnQ_m0rR7opLln69C35INdxL993mtuPThWwfKmfczXDa_Y7TQyrW0kQnusx4aGXhGkiNJJuU4GvF5vCgAkYJPWwMpPF2CK04fScDnAz5cJkcHjmVlkFcjvus53YWMq4WHGTZo_ARF0mqu5PrJSF-XtSVI37Qocnb_IMez8Y9o6WjuGYIQRDm2ToNE4Z-g7QvED9z1keYcqee9a7mSLUhI-aAdITbnBP5BzCdv9t0_p_KvnrDBviwDigZQK8Ktm6wDvVRFF-Hm6pC7wM91oHYc9ZbLNG2ePgbuK_k8-6wJ1Y06rfDqMhefcETMm4m6EO_p2wnxjodl1iaUvFjNJfdZiEpCxBt8hlyoWAMAy0FDvtlqH9nrbEFrW9UAYy5znX5i8L3MCUq5B_j3riNj1OSR2X5RwMC3JrcqIW7kl_Y3cr-50sk9DoRgC4se-Fcw5vdmP5e73vTyC3CJOb2DJgkZXRD842cvGHlk-zJWPOK_i1U-oFgiDOMM1vH0MwFlsEF_j4_2eVgUUW7DuEmmyClD0lN9ntIcIIpQ9IqYpnbdhVRemUnw6HEaZaddvqCYxtC1FiPWOak-Wjl3CckwcEGOoqQcjaHlRYT0cwgRUTQdSMl9KuItPaWxWl2hT9FPbOrGi8Px37U38IiL7eK2ySHws4sSv_KNMLbfFcv9_Nh5ICz4m_AR3EmgvUwbU0oBRACb4fZXsrqSRY5tO6i3rlVB9yVosQ5WmxSNppukkKpJ--ASDRzyq5aC394qVleqsIso-B6kHAOu.0SPkZdhmm2e3ALewOhZSKw",
}

package disney

import (
   "encoding/json"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (p *playback) fetch() error {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "disney.playback.edge.bamgrid.com"
   req.URL.Path = "/v7/playback/ctr-regular"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(strings.NewReader(data))
   req.Header.Add("Authorization", "Bearer " + bearer)
   req.Header.Add("Content-Type", "application/json")
   req.Header.Add("X-Application-Version", "5d5917f8")
   req.Header.Add("X-Bamsdk-Client-Id", "disney-svod-3d9324fc")
   req.Header.Add("X-Bamsdk-Platform", "javascript/windows/firefox")
   req.Header.Add("X-Bamsdk-Version", "34.3")
   req.Header.Add("X-Dss-Feature-Filtering", "true")
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(p)
   if err != nil {
      return err
   }
   if len(p.Errors) >= 1 {
      return &p.Errors[0]
   }
   return nil
}

func (e *Error) Error() string {
   var data strings.Builder
   data.WriteString("code = ")
   data.WriteString(e.Code)
   data.WriteString("\ndescription = ")
   data.WriteString(e.Description)
   return data.String()
}

type Error struct {
   Code string
   Description string
}

type playback struct {
   Errors []Error
   Stream struct {
      Sources []struct {
         Complete struct {
            Url string
         }
      }
   }
}

const data = `
{
  "playback": {
    "attributes": {
      "assetInsertionStrategies": {
        "point": "SGAI",
        "range": "SGAI"
      }
    }
  },
  "playbackId": "eyJtZWRpYUlkIjoiYWE0MDFhMmItYjdmNC00YzExLWJmNjEtYTNiMDZmOWM5NzRkIiwiYXZhaWxJZCI6ImNkNDkwZmE0LTBkMWYtNDU1ZS04ZGNiLWZmZmQ1MTY2NmMyMSIsImF2YWlsVmVyc2lvbiI6Mywic291cmNlSWQiOiJjZDQ5MGZhNC0wZDFmLTQ1NWUtOGRjYi1mZmZkNTE2NjZjMjEiLCJjb250ZW50VHlwZSI6InZvZCJ9"
}
`

const bearer = "eyJ6aXAiOiJERUYiLCJraWQiOiJ0Vy10M2ZQUTJEN2Q0YlBWTU1rSkd4dkJlZ0ZXQkdXek5KcFFtOGRJMWYwIiwiY3R5IjoiSldUIiwiZW5jIjoiQzIwUCIsImFsZyI6ImRpciJ9..3s8R6t-Jo_KuUS_i.cMaVAsIrZt0H_vBrbxnZCf8YXe8iHp9-uC_5MggFX1OCLXH3NwOKEH756QaeK9kQ6-_ieffrPBf362wTPI6mB4z-YgIxLzUd_0iDmrtlEHUaUgkO1vPLNmu8rdQQYVtAKDef7UM1TLVkAZfZmIx0jN1oUVbVHO-ao_-Izkxe9TwD1lCzj-q7tOmLXfxGvTNn5Qqml00vlApRra_pHDL_tAqYC8L-2E6iRKmbkuwsHS-ZCh3a8Y8Zeo_2qOEUaNpaZHJ3ZW8kgIqc_qKiP1nRf6A9_gB63kdbU-tCKbMisDZaG_D3357vvO1uGGkTNILtJYIT5fgBZE_BuN8fxRol7Dstwyaj1-XMtVE-DoKHYsWnAycC0dgbElMzckDs-vnY_Xez51i1E4O6s1HC4kwYPVp1KUWCWakIA0m0X_qs2bGeYtWXD7bnoesdyjHrdBHXh_3L5Hl54VyCGw8TA4X7KpplzbOq_gZl3oQTCcsQ9UIaCmEdlCx7xMYdJ4x4PYuCeQWXDS1fmyFdeZNOkDLNWDzYNO04jsrooAh87bJCYAdGfKyIC-XLAmgFc4CZiY6Mjae4OCOrFXXlzy0Ra4-EyQ0DprRfhJm0osq-z0MMaWk-NbIH5bVzlpiMfvzcPRJaGDqpznsRrSvloXQ4xh4t7eihAbXIKH0MUTu2KgB7XNXeNqvQXCuMQUcGHDkoQrRJwzQovXOADayJDoRoIZiU04FrNjVM05kvO7Fy4a2MlqlGcmeXD8-UAwB4rW4_1WfPhrUNyeWTJc4AkenNseEy-xKEFtyElecIXM7XchbKs4k4-U6Rsk2cM1ElFTTIqk1OewtQyTx7DLV9EieawOnKZ3wY1JCrFtg-zhCsYutvZl1gc1lihM5KR5A-6Yj2cG3_zgQ8YlcK-gYfS8rl4T_uTPnUCOaUmVnSAjIyuPo2dZYH4lRHICtEjfbOTJC9VCAcZA7IzCPfqg1b7e80F8Kk0-2h8OUl2O29-Gjgrr0dP1adQtZfwUnez9Q8lP-DGWT6BuOY9hVHuPSwptvlsDXdB1mUISzInA5WW5amaQX-zVft_EXpOJBVz6A3LnpxtChzp6FnNpNvJN38EquEBvmYnEOBxxPH6ReMcRKsB5_FSX6QCBxQm8aVBkpSKIGjgu2vI5PbBByLAYFcKrG9x9w3ex89LOPsVMrJugcGZcVdgJD6H-yoB24m5Z3tVpaJbMaFWQHp-7GvuJq2tZPNxQZfldTRc5RNR3Gmlo_5P-YVDVGrGXgA1HE4SwfUwh28aVKKNFuUx-JXdJFjeXb1LaGlqeErTn-D4QGm0PFLXrsfvIXjwcOIcVumgqz7sJV6y_ymQni08wxAx0gWSLhQs1Bet4BjAMC09hCDyubs4l7K0fbV4uibhzfJ1vcC2lr594g9bmA3-5iqQlQiUO73G2DtI0XPQ5A1_7vRjCQ6CxVfpyXovFGkBAFd4hTTwBqtO_mhbZdq5htqB1Cnf3qrxhN2KwB40CieE66809FfXf_BbCW4aFzIsOJBq9PZXKnhw8d6qvwTGTSnyuCEb3SXkT1-CxPe5e3ywAeUw3ur09bwnjPDs0BAkkps9_CjnGKyPpDw8RGXiz9XUhvMULcJkrCZX66UXjK9wW_gi_Tqt7aOqFKbMTdIJnohxySK3sXiZOipHuX_pkgQlf9brrf2LLAF8S4jEbP8OSljqoUTEYc-V5mTraezrX_ZjjlvKiKf51ccsvf7q4mZkS8K2g-F4dkf5cMViBL1H1_euzNAS5uh31tJLXEQJlffMQ-HYU6Re1FIgvvFSAmQrB5GQdEvq4lpV6yYi-ttRohqF7qxywKpCIeh6-Q8Le_nHSUpQFTRRxgtGoh_2XY9JZRAZr9qHU8ZxZ4jBKt21WvokCLFuYRFfFH01TumW-FZ9qSnn70bbFdF7PVS8H1GzDvDnElD3qLJtbyGpdDO-I4PL7DiajdwSDVshfJPe8tXDRiywIogbCQ2lu0HEyDuAi-d7qzX51aGxrvDAg8VRQx8-_ZhrXNzTqBCF3qrrFZe4jfU4ncgNtrwJ7adsIP-yBSTKHESFpstn7KI_1rMya0TqzkPLIEQGL5YKL-dYWItpZSTAawtJo8wjP0zgrqTtMTRBG8ZQPhVUV0wPG5ozYvaLcCiqtrW4ZWyvoUPd0B4K58S1eiu6IOkV3MVuoBZBRkIfFCVyUQljuTlezEzAfX_2DKLM7MkwgGEdW4M8pvJw5YPhIHs4yh9NZCvjAxUB2w5LNn_BPZ8XRHJ9FLqgamRt53dIrN835GxzAo2QeWq3E_Nyg6vAwhVBJVT7Z2pU7_T9oww1Ug4h5kGBf58i-_EJ5D5JY2EC6xU0rFhqwqjI2AN10RLXUSXvHKOFSt34odMeODoS6mzr_rn7N5mkh-7kXm8v4THrWTeM3azoiFeQwVEn1TqzbBMk_iNtCxRmg97oIGTKDoWJLXCULbVSZ_gRu5fnK8lq_DZz27Dwt-v5Jf1IXhEBZg1wCtEEWibQZh0mfFDoKSEhXqI4cC-ewbLO0SrvlmQUJqpQSnwKcsldbhlEF8mNMAsTWIGcrX5LCpfY6z1lvQb9WN41F05X9wTML5dgxAcY0h-0VxaJb2d9UnHgfB2IGphJCXPOGGCLFXpkzB5YoP9IuYjY_3vAntt0LUllc2fpNkCU7T9vGhPjACUS_EWqmswUd-1jwquoiEFCln2XClvj73et0u0MEbl6T3k4Le0NiLL1xbILTlmwNWgiKKCRqvzQL1gh3mRlsWSQziSA14xJOftKCxkcJqMtQgiOmPb1ikIELWFGqpHCpCTcyn6kk0OcAYPdCHfktbezY3gcg_R-GcLB2foanZNwXiW7xnlezLxIwcjvwB-NmqexwfxS0i1cf2hL8Si6aq3CxF7TKRfPtxqummj7fKE9AP7rSffx7DivPbkQLN8owEKnQkPxXqC8lX9mPnXeFZh3sZ_y_mNciGg7rr_Y34T72KmZsaDnz4SpVxktmMM2qRyc2779PJ5fnh-Arrgd7sWGtfV_IaNpJmyHCOB0Af61BCrpblGjA-aTPoRQZWYHeyVfnYr89_AgySXZOw_DWjJaMYiXujbA0J-kKI95-wRTn0IQZV5kcX8dIuC-QT1156nOm7KDG8jRA-tN6UzMJSHAZG9Ung_uUqLulrE5fazNEvaigyIpijp52hmJM3_JDc-aAuG-x2453gxsyG6qa4XgX3qyPhbQV2QlqsRnDLDqeRtYTFwJNL79OMLYt20zqdnad3zH_MNbLv87TDrS26DtBxESVWQZTy27esGZmT_9B40KjfUb98vFLnaO7MeR_NFO1ds8ovVFaL-lEmWKiwySt7-hbUIw1HxGM-S7xFGOoa98J9K28VtM95NS8xcRhrrMWvVeZSDDd1JvMNuNh8toXym4vv0f4_u5wosTSOWCu5KpTFayLw7dqGRp0EVMT_I-0vQwZnbju6G6yPwEYFtiROZJcGGyiH4DIAtV2uLEeU4jhMsMeXP1m0mP-IjJJOsIhPFzTbdsdUz9ZopC4kE1C_LMWq19xS1RvhH3u9wopRCUfqE0qA_w7egl5qHmJyP-NhdvI5_MmvZ_0dtYNs_fGb9NBft38Ay_SdrbXOQM01_lLnzzlRbnmpT6GyC_ePcL2d9YvMlUyqwtmjCHyXDToLSMVXw-3Iia9VbrJrEchYez0PjYAloQjcD3poVC9fUxZ9Vn0bT65kejM__61blqFhp11Ihsz0pZ29-jGPKrYnOmJV-csGA2JTfM2AXmIbS2aYWSnJ5RNQS1wDWlaSXAHVJK4YPhrjS05CJyKDfrvufFl1CEQBfRE-RCJygYHhIVLJ0yAH-JL7E2Nklg-IfZFZ7KjdFqgUnao9OLMb2Taa_CdOeJuoC6BLapAUuNw4gNhxzMySisYuWp7D1HJlyJMvYDwgIgkeJERqS2tzKxCJtMKRRZC8xPwq9BrNtcgOGLQRfxvx9KrhlB6s_gTV2anLp21lqJ6f2p-JCTmbJZeQ2PDYx8qTEUhVYbPutWOZs_R0Xif7DEpxsgGNEzGA6eT2zjsmaQbw9TalufgX8XhNQ-ss85O6IPHGa62HVtB5BViFRD2gfWd60fCe42ENT_sGaGaOieb1D3tRrIpKjqlBb73KeMc2hW7RQ5mI1qL3ARo7OaDvAzQsloYNNhNg9edz9f55EZuFfy6L-Z0AjBD8xP8YO8rzNZ2WLO78cO8rLYU9FwuZuVOwuOvAYp5017j6-UW-xxA20ATuxta-aRoxW7ELHd26NIUbNGKx2NqsIwTQVTHba4g2wMpjSmTYyLTcFODILhcMB0TpwY89dMogD8oAgKJKdUYhldEI93YV97jlgUJ0OMCQWZ4f2UVGHXvTeiXqNd3pb3GMOoSsEtqfqhTJ1V643qQpcDqUfnLpHaqxryV4tGcjCEIDo26dn2RY4Kd5cCakD1mUpsDnUQY942u9ttKYbdj_98UqJdH4FwwzjngI4iQPAiGx0HILT6idifbBJrfOwK3IJgwAEE-cY_Srw6NQTZhUbQc3vjqsHVqUQdVHA_nQcRUYBLQRwyXdx-WoBUx0OLFPoWsm2I51rrgH6ur9af5shHmiRQ5kEVJOPIeMa6SIf74NRcihHQQe5X5zuLzUL4vDAnUwSe_hAb1F-bFOxiFRdNNm8G2bwqBciwDviqWdd49SJwyi4Gxqf8Rs9l-kZrJHrg5iJZOJ0EeT6ndWeU7v6N9smAhw8B3MYWKx2dYSLqQnhUnAcH-zSxZfCM1L2K-8wZaF6YZjbjhHj8cjUBiL30QkRX6Ke8TC5FtNwRuIvLZzfvSLWowLmdWoohFwSHqm6R-lDCMXxT2KZX3h43f40amIHZWypjtMif6bpPxG5G9dsm946rCuolQPn-ETlpPwxibJXBJOTFA33Vv5_YzGpyfCnxrey2j5K9PbO5uLqsgjSHHyI4A1U8sKiwaVMNSIpTh82fRbKPA4lZGHHu0M8XJakrKCLWcJ8gcjGXEDchX6OXtb8DmXAJ3JSRQbkfQlMYrQpRw52aU7TlgHW1v4i95BCAlUzIU9Hs4b3VxQ8Nk6mtgRMeYfOPHlPnJa_ALgH7exdryiEAXu-NYBlxWsrGUwQSF9eUvDo6XDZlbnr4SQlZwYhUieI93GUPz0a1yPtP-DZ8rv5d9xtnDzOThMNYYpJIxOrIBN1woLeRlCFeoafRnphRkUeL2kyTrZajK_GU16Xp-DPDVCemHd4K208I8B1_boZtll5XNGEU9WPzWK3mVo7KWLFX0UoQsfuSUXqRbxf3yHUWf_Uj7_lnfWtse0z5rvbf4KD4UXZIgx33yET7JJw7jivq39E-_Hwp62mNCcEVTKj.ePllIRmDWDyPCJnxQxS2oA"


package disney

import (
   "41.neocities.org/rosso/disney"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func request_otp(device *disney.Device) (*http.Response, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "disney.api.edge.bamgrid.com"
   req.URL.Path = "/v1/public/graphql"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(strings.NewReader(data))
   //req.Header.Add("Authorization", "Bearer " + bearer)
   req.Header.Add("Authorization", "Bearer " + d.Token.AccessToken)
   return http.DefaultClient.Do(&req)
}

const data = `
{
   "query": "\n  mutation requestOtp($input: RequestOtpInput!) {\n    requestOtp(requestOtp: $input) {\n      accepted\n    }\n  }\n",
   "variables": {
      "input": {
         "email": "27@riseup.net",
         "reason": "Login"
      }
   }
}
`

const bearer = "eyJ6aXAiOiJERUYiLCJraWQiOiJ0Vy10M2ZQUTJEN2Q0YlBWTU1rSkd4dkJlZ0ZXQkdXek5KcFFtOGRJMWYwIiwiY3R5IjoiSldUIiwiZW5jIjoiQzIwUCIsImFsZyI6ImRpciJ9..B2Q6vgeP27YelBok.EROPRRC-uS3kAqVUgW3Z7btnWpyk_A3bEu4WWMTxxAF_FAMmth115WpIQrzwqoPx7JfWMdk5DD6O4PA-2XQmEBqa090GrDBVgnQ9PCjCODD49K8lMuTU3QEsuxrZDB7E_68LdBK0yeuWUzh8FdvJhtIQJZM_shZk5ZY20tXoc61CO0139xgjSg4nSTepaPtZ1mz9lUzOHVmwYZS0uy746mUjeiSQTHH2TYtcDaTlKJfgtlrdszRvmliHkr0e6Xomv8tmbAKKoLfEbQXdXfwg2yOxhDLHkAECEgy4osWc794jApeQBRAv9lzZ4e3vipChRQXxqucXKeKelOlhfXhxQsmZq69XAQ2rfHIqfCSK6Q3ZojGcpi1cme8cef_rkzd5RMUBtUqV5RDGAAJNZsfHjPYFnDVVOcDjCkBUD02bALaaqB3us2BEBD1jCwntMuZJNsS4J6yt8e73P5b1Pfl7UsVVPMdw9pF_1z2ezsJ_I4CGEqLxu3f6xZDfD2OD96hu4YaJBuvGb8iSyG27VuSLBGdPzCGI8uYJX7Y-bc2tUyfeTQEkK_gzP8VbiGRywe_dlRHgFzp_gw-OlBCtvMLOkJNClgZk1PMKOEzOhP-McbReEgtupe5KHyEfiFnzWGthmxNC2QZzeVF3oLzfV5ppOcaXzCLGgtZHrqoMJSMJWq2iHOBT6Rc5B_jT7T5cbIUx-qPLeW4NC-ffyHghmdNgyAwQmm4GPFtOIMmq4bHcsFt1MKXwqCPZnipSRXVayxeX2TlC9U8W1n7PyaXE-aAAtVa6Q1-3-ShJf3pK8ACt_NX-s3n-IJa5DL0Nw1weXzfjzh-qp4UkzcRSKLp7Rt5TfNye7Fx-qx9m-cVA7iWNsHItWdVzMuxK6E0ArxZRTVVj6Pek1uiDx1azzOXOP6L9JV0aUnY1klfqfht1KVE3C4inuaRfQ5e3d-eTOB9OLk6elAdp21dwndTjQLXtaJCRoh8ltmFbMOjXk3CoVzlKT_NMj_81iXyOguGzTFeSpMXXLph4mnjbsahVrdXGMGQFRobtXBLiFG1sIVixjLt91odjbgKYHVkgyW-ezZKVtSAHezFUJodoynXaHP5thJKKwebb401wKmDTFRRiZrn8eJjgxmoq_B9z9Qq60vgtatPAxWn4SHrr9978l7xExSBpFii8BscIItq-_QmOTFR1rc58HpRZgsmHqfwiGHDi-Ov8DQ9bxQmziVXGOak-ynQFnLGHkDJ8KC-UsCEGwjnzs_vlVOYNXpbDP1a4mTPOOOoBNvkpMrjR5A8fWa5IDkWDCODnlCxffv2EJynNRt75v8RidB1yQiai80aYLTiDjqe7NqOZqm2_eIqeYDHKB_SBQQ-ACTtyW4STsEATJVzkFSvKlDZy8wsBlKuzf8wZqKOtfBj7DqSiq1VSB26Fsyyst_OKbUOvOK1Jf68JX6tewAmPB6QgUbfHKa_SPcMP-KVK0pw1ikwZEByQm9FLjJF9cnZrtaIgiG1_4Ud_BMoMYOBgSM2aOTGdYKXCGeUIwsQw9sS3hOaOzS2LMufcosfiSEP01hkziPB2oHiQ1MEBE-kDaY8pgxdcL2H_D6COibaUtgO992OLpwKmH_8JI-6cxFPpXc4GIhQDT8-CcjoklFewzE5XZF6XLhKGve_-J4iKkS_SmR04I0241RwMy4kpy_6zE3M9PVb8zCbrlt9TfFMXMHT93seLUjNvhAWY_LSahnOzn2vOcx8UVMIi7hFGFwWqdO-J1ttpmuUty30moeX8LSygHVitDIXCwU6ws85Yg9xQbrMlOGMpuVGx5BLuQgelG0tLsv3zwshw67BH1YtsWHLq1X86K2ay6VeOPGjRnoRCp3PkSFYyrTvMjVBS5Tkk7v6VTU3sMUJ89MMkR6FHupzL4Z3OLp6anNMKj3PhwKJqPvmktdPfk_jAop9CAKmAJxEupLZDuYGpI9b6KHtSeAN_MLrkVVSF48U2KFtI6f5T5yrUBHFKpMOtfEMQXRfHHD9Bh25nWHuYCS1yUb03czkJ5EWS7SM3XhJWaoy1B6HMC0CENQ6pOuSrZDYObsuBtITYijKj95mBWOJcSV7wuMhqZnjB19hn7JFKTjDtYhu9aergO259rpbHhTr_6eD8XfENhpDE8mlYQa0GNlXZ1fHGRoRf-Z94YbzSHOG2mKg4-SljKrATp8V5ZJUGxvvSDk_W9HoC0TPK41N9bf11dP4ZMTN_3tDV9pexzwvA1ngaA61yMpLLtUlJ_mpKcg5fzkvGX1lYS6wlzdWLNlY7zwDCE_MrjalH4Dvq53tllBBkrNgk4rLUdfzK6A2mrRFI4bHSqwTszzdFPUEHzxWYbc8TNKL3pDv3igj89V0TTHcmPnC9xVSKtipQG35DqG3nzlt34VYYnNwxH5ojtKX6afsJ9QFhNyvmUatjLKnrWDcXQSG0B3OBFQGA49z0ZFYoNnJU5jTR7RkHuv0FL57dsbiAceW4u8Ixjx-ZZAfoH02XnMMFBlcPfarBtIx8O_kGVqqlHa1U9IJcDnUrq1X9SXTVlhi9-gbKf21MUYPTsN1L_A4LycrjA96euRHY8ie8o9jCo6WxgIf6q7zYEOdiaMSNYu80MO6k3Dy8dgjn8ejdeZh3aE0D4p5miihO_N2yX9U._xxaA49rHzlSEWbkYHRbCg"


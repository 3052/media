# joyn

- https://joyn.de/filme/barry-seal-only-in-america
- https://justwatch.com/de/Anbieter/joyn-plus
- https://justwatch.com/us/movie/american-made

## android

https://play.google.com/store/apps/details?id=de.prosiebensat1digital.seventv

## web

this is it:

<https://joyn-vod-prd.akamaized.net/v1/CiQwMjYwN2RlZC0xNjU4LTQyMDMtODNjNC1iYzk0NGI0MmJkY2M.Cg1hX3A0c3ZuNGEyOGZxEAEY4AMiJDM5YTU4NTA5LWZjZjEtNGIxZC05ZWZiLTNhNTAyZTIzMWNlYg.Io7JBViGADc4EXQKv-C6v1OUGctJyGOgqZVac71-D7o/a_p4svn4a28fq/.ism/dash/audio_deu_1=157000-480256.dash>

from:

~~~
POST https://api.vod-prd.s.joyn.de/v1/asset/a_p4svn4a28fq/playlist?signature=3a082bf39122c422094360c39d1774897345c821 HTTP/2.0
user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/111.0
accept: */*
accept-language: en-US,en;q=0.5
accept-encoding: gzip, deflate, br
authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXlfc2lnbiI6InByb2QiLCJlbnRpdGxlbWVudF9pZCI6Ijg3MmRkMjNhLTFkNDEtNGIwZi04YzJhLTBmM2E1MDZhNTlmZiIsImNvbnRlbnRfaWQiOiJhX3A0c3ZuNGEyOGZxIiwidXNlcl9pZCI6IkpOQUEtYWUzZjgwY2YtNTc0My01MjFlLTkwOTItMDkyYTg3NzkyYzhkIiwicHJvZmlsZV9pZCI6IkpOQUEtYWUzZjgwY2YtNTc0My01MjFlLTkwOTItMDkyYTg3NzkyYzhkIiwiYW5vbnltb3VzX2lkIjoiNmJkMTZiNDgtYTQ4YS00MDBlLTljZGItMjdmMzM5MWUyZTRlIiwiY2F0YWxvZ19jb3VudHJ5IjoiREUiLCJsb2NhdGlvbl9jb3VudHJ5IjoiREUiLCJkaXN0cmlidXRpb25fdGVuYW50IjoiSk9ZTiIsImNvcHlyaWdodHMiOlsiVW5pdmVyc2FsIFN0dWRpb3MgSW5jLiBBbGwgUmlnaHRzIFJlc2VydmVkLiJdLCJqb3luX3BhY2thZ2VzIjpbIkRFX0ZSRUUiXSwiYnVzaW5lc3NfbW9kZWwiOiJBVk9EIiwicXVhbGl0eSI6IlNEIiwiYWRzX21heF9taWRyb2xsX2Jsb2NrcyI6MTAsImFkc19saW1pdF9wcmVyb2xsIjozLCJhZHNfbGltaXRfbWlkcm9sbCI6NSwiYWRzX3Rlc3QiOiIiLCJhZHNfdmFyaWFudCI6IiIsImFkc19icmVha19zcGFjaW5nIjoxMywiaWF0IjoxNzE1MTIyMTYwLCJleHAiOjE3MTUyMDg1NjB9.Y6KWmtE1Gq5BT4qiJ4W2hLyhkw6mwF8mGreLPJwBxPBV6g1LDm4Pnvn8rTtorLkBC0yZGlVtpJCpEtoXoHZthMjNQvgkCkI5JoDP2ezy-Lh5nIpXtcy9CrKJ_Y6vyvnTDSRz5PQuJpbt-CiHQP5bWxlfoUBYDvkKjMMy8okJiVHqiRoaYQ-ycBG60HhUKthaURy4EgY8v6m2QHH1ygWZJhCj7U45szwTW6Qq2YJKErKVGaecCbPrUNyJm4wC8jdMc_YMd7DhZ0KYu72tJN5FlHhbI0CnawpBfS4ivlSOfIe-yQX_yCh9ONA1-ZbU0U-AAD_h8UvBCbKLtITDmaxLwg
content-type: application/json
content-length: 249
origin: https://www.joyn.de
sec-fetch-dest: empty
sec-fetch-mode: cors
sec-fetch-site: same-site
te: trailers

{"manufacturer":"unknown","platform":"browser","maxSecurityLevel":1,"streamingFormat":"dash","model":"unknown","protectionSystem":"widevine","enableDolbyAudio":false,"enableSubtitles":true,"maxResolution":1080,"variantName":"default","version":"v1"}
~~~

authorization from here:

~~~
POST /api/user/entitlement-token HTTP/1.1
Host: entitlement.p7s1.io
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/111.0
Accept: */*
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br
authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImEwZDQwYjkxZTA2OGEzY2ZhODQ1ZjRkZTViNmY3NjA2NmEzMzc3NTEifQ.eyJkaXN0cmlidXRpb25fdGVuYW50IjoiSk9ZTiIsImNvdW50cnkiOiJERSIsImNhdEN0eSI6IkRFIiwibG9jQ3R5IjoiREUiLCJwcmZMbmciOiJkZSIsImF2b2RDdHkiOiJERSIsImF2b2RFbmwiOnRydWUsInN2b2RDdHkiOiJERSIsInN2b2RFbmwiOnRydWUsImpJZCI6IjcxZGI5YWIyM2ZlMWY1OTA1MTk0ZWNlNWY5NjYyNTBiOjU3Yzk0ODNhMTdlODkxZThhYjQ0MGZkNzcyNTkwODVmZTk0ODRiOTJmMGMwN2RhODkzYTlhYWE0NmQ4MzM0Yjg3YTY1N2EwMGUxODNiMzg3NTJjYmNjMzMyNzBlODZiMyIsInBJZCI6ImExMzcyMjRjNzE0MmM5NjEzMjllZjExMWNkODNmNjlmOmNkOTU5NzViOTQzMGJiNzgxYTA4ZGEyOWRiZDYyNmQ0MDEyODBiMjMzODZkZWNiYzZjOWI2YjUyYmU1ZWJmNTE1YWQ3MDIwMTY1MTU0NDkzMmIyMDk0NzlkOGNjMDJmOCIsImpJZEMiOiJKTkFBLWFlM2Y4MGNmLTU3NDMtNTIxZS05MDkyLTA5MmE4Nzc5MmM4ZCIsInBJZEMiOiJKTkFBLWFlM2Y4MGNmLTU3NDMtNTIxZS05MDkyLTA5MmE4Nzc5MmM4ZCIsImNJZCI6IjZiZDE2YjQ4LWE0OGEtNDAwZS05Y2RiLTI3ZjMzOTFlMmU0ZSIsImNOIjoid2ViIiwiZW50IjoiZWM6MCxmOjAsZmw6MCIsImludGVybmFsIjpmYWxzZSwiZWxpZ2libGVGb3JFbXBsb3llZVN1YnNjcmlwdGlvbiI6ZmFsc2UsImlhdCI6MTcxNTEyMjE0NywiZXhwIjoxNzE1MjA4NTQ3LCJhdWQiOlsid2ViIl0sImlzcyI6Imh0dHBzOi8vam95bi5kZS8iLCJzdWIiOiIyMTI5ZmY4OS1hOGQzLTQxOGMtOTkyNy00ZGY1ZDU1ZDUwNTgifQ.HxWomuc00pn2o4tacXAMdZGrLwlonw_fPfoAlMy2dKuMZ9Lp3v0CsEC1E6HuyLZ38UBbeotwKKSymDG8Rlyn3xiEjyYwqO7pEjLwKedDIgEb1m7AD-KchPjj1xAXZMDOw2IBPBsDqgQ3LdC33lNgWpng9rp-6irNcZFHM_KnfjFW0M_RL3jPlJzHWgbi1PamxsxroEPwBSZjup_xBgRQ9_phatgBezjlbJMJEk1-6eZ6WWn-8yJLyDEp1qD2ObyVE6mKtcHa6agooUlsr8YZ72wzSlV5yRIWKta_qikAxxVy0oI5OqsBU46aQPyoNdpenXY3OllLq3ZYWAWCT3Gxnho5UEUaFW7MtIuvwtwia4eRg2_piZg1dSLPRt0NRtBe_x4j0KX1UDxnGZzU6REIQvE2CxpTNeAXvz7gVqzdN0ZY_nar_BwuY9Gq2PKAU9wyTm9Dyy1YdydFd3pucbPtKM8iLZdE_wDNYFgG-Yi1vpsCfO0_FHokY2l-9L_r198Htk5z-gnE-3_sDxuP7wQlYpIYGdMz6vGsAcVNdbWE8VXTYgiT-zyYtuzUJFFQItI_4AvcizzVaLJf-7G1N6PM2f9FqDV1rQwFeZQ808jQdTJGYoF0qT_FS0ONgI1URSNE7vLvmSFVjoJPsAeMQBf1iOnM4qGH2c4kQcCowY-N3pc
joyn-client-os: UNKNOWN
joyn-client-version: 5.702.5
joyn-platform: web
content-type: application/json
Content-Length: 51
Origin: https://www.joyn.de
Connection: keep-alive
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site

{"content_id":"a_p4svn4a28fq","content_type":"VOD"}
~~~

authorization from here:

~~~
POST https://auth.joyn.de/auth/anonymous HTTP/2.0
user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/111.0
accept: application/json
accept-language: en-US,en;q=0.5
accept-encoding: gzip, deflate, br
content-type: application/json
joyn-client-version: 5.702.5
joyn-platform: web
joyn-distribution-tenant: JOYN
joyn-country: DE
joyn-request-id: 08477561-a8bb-4e38-b298-cd2f680164df
content-length: 128
origin: https://www.joyn.de
sec-fetch-dest: empty
sec-fetch-mode: cors
sec-fetch-site: same-site
te: trailers

{"client_id":"6bd16b48-a48a-400e-9cdb-27f3391e2e4e","client_name":"web","anon_device_id":"6bd16b48-a48a-400e-9cdb-27f3391e2e4e"}
~~~

## /graphql

~~~
GET https://api.joyn.de/graphql?operationName=PageMovieDetailStatic&enable_user_location=true&watch_assistant_variant=true&variables=%7B%22path%22%3A%22%2Ffilme%2Fbarry-seal-only-in-america%22%7D&extensions=%7B%22persistedQuery%22%3A%7B%22version%22%3A1%2C%22sha256Hash%22%3A%225cd6d962be007c782b5049ec7077dd446b334f14461423a72baf34df294d11b2%22%7D%7D HTTP/2.0
joyn-platform: web
x-api-key: 4f0fd9f18abbe3cf0e87fdb556bc39c8
~~~

x-api-key is hard coded in JavaScript
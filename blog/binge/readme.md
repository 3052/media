# binge.com.au

~~~
url = https://binge.com.au/movies/asset-contact-1997!7738
country = Australia
~~~

this is it:

~~~
GET /out/v1/5c62a8170d294a068af0c2107c15e543/b5d724ada4494ebbb8bd55489ea091aa/545bda3186e148d586c7e9d72422b834/index.mpd?aws.manifestfilter=trickplay_height%3A1-2 HTTP/1.1
Host: fxtlgrp-vod-avc-drm7.akamaized.net
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0
Accept: */*
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br, zstd
Referer: https://binge.com.au/
Origin: https://binge.com.au
Sec-GPC: 1
Connection: keep-alive
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
Pragma: no-cache
Cache-Control: no-cache
content-length: 0
~~~

from:

~~~
POST /api/v3/play HTTP/2
Host: play.binge.com.au
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0
Accept: application/json
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br, zstd
Referer: https://binge.com.au/
Authorization: Bearer eyJraWQiOiI3a0UxeCt4bE5xbFJabHNaMm9NeStQNnlBckU9IiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiJhdXRoMHw2N2NlMjU4MzY1Yjc1MjA4NzNlMmU0ZjIiLCJodHRwOi8vZm94c3BvcnRzLmNvbS5hdS9tYXJ0aWFuX2lkIjoiYXV0aDB8NjdjZTI1ODM2NWI3NTIwODczZTJlNGYyIiwiaHR0cHM6Ly9zdHJlYW1vdGlvbi5jb20uYXUvYWNjb3VudC9leHRlcm5hbC1pZGVudGl0aWVzIjp7fSwiaHR0cDovL2lyZGV0by5jb20vY29udHJvbC9qdGkiOiJmMGE0ZmJlMy0xZTM1LTRlNzgtOTg5NS04Yjg1MTE5MGUxZWMiLCJzZWNvbmRhcnlfa2V5IjoiNzMyZjhjYWE0NTMyMjE1NWUyOTgwNDEzODNiNDI5OTc3MGM3MjVkODk5MmMyMjc1ZTdjMjcyZWM4MzY2ZTNjOCIsImlzcyI6Imh0dHBzOi8vdG9rZW5zZXJ2aWNlLnN0cmVhbW90aW9uLmNvbS5hdS8iLCJodHRwOi8vaXJkZXRvLmNvbS9jb250cm9sL2FpZCI6ImZveHRlbG90dCIsImd0eSI6InBhc3N3b3JkIiwiYXVkIjpbInN0cmVhbW90aW9uLmNvbS5hdSIsImh0dHBzOi8vcHJvZC1tYXJ0aWFuLmZveHNwb3J0cy0xYi1wcm9kLmF1dGgwYXBwLmNvbS91c2VyaW5mbyJdLCJodHRwczovL3ZpbW9uZC9lbnRpdGxlbWVudHMiOlt7InN0cmVhbWNvdW50IjoxLCJhZF9zdXBwb3J0ZWQiOnRydWUsInN2b2QiOiIzIiwicXVhbGl0eSI6ImZoZCJ9XSwiYXpwIjoicE04N1RVWEtRdlNTdTkzeWRSakRUcUJnZFllQ2JkaFoiLCJzY29wZSI6Im9wZW5pZCBlbWFpbCBkcm06bG93IG9mZmxpbmVfYWNjZXNzIHVzZXI6cGhvbmVfdmVyaWZpZWQiLCJodHRwOi8vaXJkZXRvLmNvbS9jb250cm9sL2VudCI6W3siZXBpZCI6IkFSRVNfMUZIRF9sb3ciLCJiaWQiOiJMSVRFIn1dLCJleHAiOjE3NDE1NzE4NTEsImlhdCI6MTc0MTU3MTU1MSwianRpIjoiZjAyMTQ1ZTAtNDg4ZS00MTQ1LTk0ZDMtMzgxZjJjOWZkMjM4IiwiaHR0cHM6Ly9hcmVzLmNvbS5hdS9zdGF0dXMiOnsidXBkYXRlZF9hdCI6IjIwMjUtMDMtMDlUMjM6Mzk6MzkuODMzWiIsInBwdl9ldmVudHMiOltdLCJhY2NvdW50X3N0YXR1cyI6IkFDVElWRV9TVUJTQ1JJUFRJT04iLCJzdWJfYWNjb3VudF9zdGF0dXMiOiJQQVlJTkdfU1VCU0NSSVBUSU9OIn19.M9YJaohXpT1Jl8_Ws8VTLGc2YNdLFi4BgsB3NiKYD4U3CBhOh4aOT4fyx0Gel6ZWgkafTWXfTfQKFWrgZQfIDSMEcKuNhuZ2l4SzK9ZeuN_MyGmAEK5Ki-_XdXVEjPCHUzEXGHRQOwp3hseaRjLUzczwOXd-l-BhgPDDF4juyeLr3qAzG0ORAMWLdKi4RvGkUvJ0FJcf8_piA0WUi5AGGODA3DImrzdbgMNFWzWslD5BVEzRTPDbjJXhonBggBUWi3RWvsoTswxOK3xsHsotd7hPInoQZ0G6mULoJR_CCbYNVWYjtJ_QY0OFSRPICVscI3sL9_k26xU9p6NGhM1nHA
X-Vimond-Subprofile: eyJraWQiOiJrZXktaWQtMSIsImFsZyI6IkhTMjU2In0.eyJhdWQiOiJ2aW1vbmQtZXhwZXJpZW5jZS1wb3J0YWwiLCJzdWIiOiJhdXRoMHw2N2NlMjU4MzY1Yjc1MjA4NzNlMmU0ZjIiLCJodHRwczpcL1wvdmltb25kXC9zdWJwcm9maWxlIjoiOWNlMDZlMzNjOGYzZDJjYTgyMmNmZDllOWY2ZWZiN2M4YTE3MWZjMSIsImlzcyI6Imh0dHBzOlwvXC9hcmVzLnN1YnByb2ZpbGUuY2Yuc3RyZWFtb3Rpb24tcHJvZC52bW5kLnR2IiwiZXhwIjoxNzczMTA3NTk5LCJpYXQiOjE3NDE1NzE1OTksImp0aSI6IjkzZjAyMTJiLWI3NDEtNDExMC05NzYzLTExYmViOWFlMDkyMCJ9.w9laN4ewY83czQUp0YJN1qpEE2FZW5IC1qBoQ5PpqSE
Content-Type: application/json
Content-Length: 554
Origin: https://binge.com.au
Sec-GPC: 1
Connection: keep-alive
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: same-site
Pragma: no-cache
Cache-Control: no-cache
TE: trailers

{"preference":{"trackKeys":false},"assetId":"7738","application":{"name":"binge","version":"11.2.0","appId":"binge.com.au"},"device":{"id":"50e785be-4c7f-4781-87e4-a3b4c75a3634","type":"desktop"},"os":{"name":"Browser","version":"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0"},"player":{"name":"VideoFS","version":"38.0.5"},"ads":{"optOut":false},"browser":{"version":"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0"},"capabilities":{"codecs":["avc"]},"session":{"intent":"playback"}}
~~~

from:

~~~
POST /oauth/token HTTP/2
Host: tokenservice.streamotion.com.au
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0
Accept: application/json
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br, zstd
Referer: https://binge.com.au/
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik56aEJPVFJHT0RjNE1FUkRSRFJEUTBVd1FrVkdNRGt4TVVVNVF6RTRRa0UzTkVVMk1rVkRRZyJ9.eyJzZWNvbmRhcnlfa2V5IjoiNzMyZjhjYWE0NTMyMjE1NWUyOTgwNDEzODNiNDI5OTc3MGM3MjVkODk5MmMyMjc1ZTdjMjcyZWM4MzY2ZTNjOCIsImlzcyI6Imh0dHBzOi8vYXV0aC5zdHJlYW1vdGlvbi5jb20uYXUvIiwic3ViIjoiYXV0aDB8NjdjZTI1ODM2NWI3NTIwODczZTJlNGYyIiwiYXVkIjpbInN0cmVhbW90aW9uLmNvbS5hdSIsImh0dHBzOi8vcHJvZC1tYXJ0aWFuLmZveHNwb3J0cy0xYi1wcm9kLmF1dGgwYXBwLmNvbS91c2VyaW5mbyJdLCJpYXQiOjE3NDE1NzE1NDgsImV4cCI6MTc0MTU5MzE0OCwic2NvcGUiOiJvcGVuaWQgZW1haWwgZHJtOmxvdyBvZmZsaW5lX2FjY2VzcyB1c2VyOnBob25lX3ZlcmlmaWVkIiwiYXpwIjoicE04N1RVWEtRdlNTdTkzeWRSakRUcUJnZFllQ2JkaFoifQ.XcEfhjhu5Bwkm-d6Bg-Z3sTDWNr4wQFkt6ns-_lPbaoE6SUHGO8CNmLxK4m-vCdnGus4_bXlKevMnYohZhDGMiQwW-XbCg3FyCWAp-8K-3cMkZ49-AL4YxHZAwZE5HaQduqALSjQbqOE3-PKlpK7hY1Bf1W0qM0InSdV-DzdcxfDUsPBDcFca50uzSyPo-TPSDEvqlOviLTOgHlqByjV7sWArkpeXZGZZRkuAnmDB1MeHp0Z_CgV0OYitcU5zcBMnVyoEDuABlWZGcwqc1G3kiLlH6F-zO6BVRhXZLLrC_4Jz3sUJWhVe4nWPOdhQLer4wmD_Oq8isgbm9BAObPKCA
Content-Type: application/json
Content-Length: 79
Origin: https://binge.com.au
Sec-GPC: 1
Connection: keep-alive
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
Pragma: no-cache
Cache-Control: no-cache
TE: trailers

{"client_id":"pM87TUXKQvSSu93ydRjDTqBgdYeCbdhZ","scope":"openid email drm:low"}
~~~

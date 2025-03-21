# canal

https://justwatch.com/cz/poskytovatel/canal-plus

~~~
url = https://www.canalplus.cz/stream/film/argylle-tajny-agent
monetization = FLATRATE
country = Czech Republic
~~~

this is it:

~~~
GET /bpk-token/1ac@3ecn2ctviurl2iylh02qqe0yfggridxy3ahkisaa/bpk-vod/playout01/default/appletvcz_A007300100102_aa58f539915e9d73863edd75a4f2fe91_HD/appletvcz_A007300100102_aa58f539915e9d73863edd75a4f2fe91_HD/index.mpd?bkm-query HTTP/1.1
Host: cz-bks400-prod08-live.solocoo.tv
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0
Accept: */*
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br, zstd
Referer: https://play.canalplus.cz/
Origin: https://play.canalplus.cz
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
POST /v1/assets/1EBvrU5Q2IFTIWSC2_4cAlD98U0OR0ejZm_dgGJi/play HTTP/2
Host: tvapi-hlm2.solocoo.tv
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0
Accept: application/json, text/plain, */*
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br, zstd
Referer: https://play.canalplus.cz/
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0di5zb2xvY29vLmF1dGgiOnsicyI6IncxZjhhOGZiMC0wNWZiLTExZjAtYjVkYS1mMzJkMWNkNWRkZjciLCJ1IjoiV3ozS0JWRnAzY2xwclEzZWVNUGZZZyIsImwiOiJlbl9VUyIsImQiOiJQQyIsImRtIjoiRmlyZWZveCIsIm9tIjoiTyIsImMiOiIzR01XanAwTldZT2ZhOThVZjhhbU1oUXVSNnJ6dUxvY3FSZ0NKcEZpUjI0Iiwic3QiOiJmdWxsIiwiZyI6ImV5SmljaUk2SW0wM1kzQWlMQ0prWWlJNlptRnNjMlVzSW5CMElqcG1ZV3h6WlN3aVpHVWlPaUppY21GdVpFMWhjSEJwYm1jaUxDSjFjQ0k2SW0wM1kzQWlMQ0p2Y0NJNklqRXdNRFEySW4wIiwiZiI6NiwiYiI6Im03Y3AifSwibmJmIjoxNzQyNTIzNzA5LCJleHAiOjE3NDI1MjU1NzcsImlhdCI6MTc0MjUyMzcwOSwiYXVkIjoibTdjcCJ9.lnjPwQryinqnFccT8ryVqF6joz0c_0vguiKx-w6JoiI
Content-Length: 182
Origin: https://play.canalplus.cz
Connection: keep-alive
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
Pragma: no-cache
Cache-Control: no-cache
TE: trailers

{"player":{"name":"RxPlayer","version":"4.2.0","capabilities":{"mediaTypes":["DASH"],"drmSystems":["Widevine"],"audioQualities":["Stereo"],"smartLib":true,"embeddedSubtitles":true}}}
~~~

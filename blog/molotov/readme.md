# molotov.tv

- https://justwatch.com/fr/plateforme/molotov-tv
- https://play.google.com/store/apps/details?id=tv.molotov.app

1. Regarder maintenant (watch now)
2. e-mail
3. password
4. day
5. month
6. year
7. sex
8. VPN
9. subscribe
10. card number
11. month
12. year
13. security code
14. Essayer gratuitement (try it for free)

480p:

https://vod-molotov.akamaized.net/output/v2/d8/a1/65/32e3c47902de4911dca77b0ad73e9ac34965a1d8/32e3c47902de4911dca77b0ad73e9ac34965a1d8.ism/high.mpd

720p:

https://vod-molotov.akamaized.net/output/v2/d8/a1/65/32e3c47902de4911dca77b0ad73e9ac34965a1d8/32e3c47902de4911dca77b0ad73e9ac34965a1d8.ism/hdready.mpd

1080p, exactly the same:

- https://vod-molotov.akamaized.net/output/v2/d8/a1/65/32e3c47902de4911dca77b0ad73e9ac34965a1d8/32e3c47902de4911dca77b0ad73e9ac34965a1d8.ism/fhdready.mpd
- https://vod-molotov.akamaized.net/output/v2/d8/a1/65/32e3c47902de4911dca77b0ad73e9ac34965a1d8/32e3c47902de4911dca77b0ad73e9ac34965a1d8.ism/fullhd25.mpd

this is it:

~~~
GET /output/v2/d8/a1/65/32e3c47902de4911dca77b0ad73e9ac34965a1d8/32e3c47902de4911dca77b0ad73e9ac34965a1d8.ism/high.mpd#t=4 HTTP/1.1
Host: vod-molotov.akamaized.net
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0
Accept: */*
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br, zstd
Referer: https://app.molotov.tv/
Origin: https://app.molotov.tv
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
url = https://www.molotov.tv/fr_fr/p/15082-531/la-vie-aquatique
monetization = FLATRATE
country = France
~~~

# movistar

https://justwatch.com/es/proveedor/movistar-plus-plus-ficcion-total

~~~
url = http://wl.movistarplus.es/ficha/?id=3427440
monetization = FLATRATE
country = Spain
~~~

smart proxy blocks login - proxy seller works. license header comes from here:

~~~
POST /asvas/ccs/00QSp000009M9gzMAC-L/SMARTTV_OTT/ea3585a776ed444d8677ad8be6ef0db3/Session HTTP/1.1
Host: alkasvaspub.imagenio.telefonica.net
User-Agent: Dalvik/2.1.0 (Linux; U; Android 12; 22126RN91Y Build/SP1A.210812.016)
Accept-Encoding: gzip, deflate
Accept: application/json, text/javascript, */*; q=0.01
Connection: keep-alive
Accept-Language: es-ES,es;q=0.9
Origin: https://ver.movistarplus.es
Referer: https://ver.movistarplus.es/
Content-Type: application/json
X-Hzid: eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI2N2Y1Y2NlN2FkMDg3YjI1YzBmNjRhZGIiLCJpYXQiOjE3NDQ0MTIwNDQsImlzcyI6ImVhMzU4NWE3NzZlZDQ0NGQ4Njc3YWQ4YmU2ZWYwZGIzIiwiZXhwIjoxNzQ0NDU1MjQ0fQ.cYc7fzZFKT1CU5KWxuTZtEhy6CgP0rqFDBFdyjWwyJw
Content-Length: 64

{"contentID":3427440,"drmMediaID":"1176568", "streamType":"AST"}
~~~

## phone

~~~
> play -i es.plus.yomvi
details[8] = 0 USD
details[13][1][4] = 9.8.1
details[13][1][16] = Feb 6, 2025
details[13][1][17] = APK APK APK APK
details[13][1][82][1][1] = 7.0 and up
details[15][18] = https://www.movistar.es/particulares/centro-de-privacidad
downloads = 12.58 million
name = Movistar Plus+
size = 21.03 megabyte
version code = 502
~~~

https://apkmirror.com/apk/movistar-espana/movistar-2

not having luck with older android, so trying with SDK version 34:

mullvad fail

nord fail

smart proxy fail

proxy seller fail

## tv

~~~
> play -i com.movistarplus.androidtv -leanback
details[8] = 0 USD
details[13][1][4] = 2.5.0
details[13][1][16] = Mar 26, 2025
details[13][1][17] = APK APK APK APK
details[13][1][82][1][1] = 7.0 and up
details[15][18] = https://www.movistar.es/particulares/centro-de-privacidad
downloads = 3.27 million
name = Movistar Plus+
size = 20.50 megabyte
version code = 100
~~~

https://apkmirror.com/apk/movistar-espana/movistar-android-tv

## web

this is it:

~~~
GET /_42189/prod/dash/cplus-3427440-md-03_cplus-3427440-mdrm_s4my8zabfhof8ns/manifest.mpd HTTP/1.1
Host: b42189-p14-h51-v0-aggqkswu-tf781b8.1.cdn.telefonica.com
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0
Accept: */*
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br, zstd
Origin: null
Referer: https://ver.movistarplus.es/
Sec-GPC: 1
Connection: keep-alive
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
Priority: u=4
Pragma: no-cache
Cache-Control: no-cache
content-length: 0
~~~

from:

~~~
GET /ficha/longlegs?id=3427440&origin=WEB&id_perfil=SUSCRI2&suscripcion=UT-DO0004,UT-MPARFU,UT-TVRECS&ui=MPLUS_CLINF&demarcation=0 HTTP/1.1
Host: ver.movistarplus.es
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br, zstd
Sec-GPC: 1
Connection: keep-alive
Cookie: OptanonConsent=isGpcEnabled=1&datestamp=Tue+Apr+08+2025+22%3A06%3A42+GMT-0500+(Central+Daylight+Time)&version=202403.1.0&browserGpcFlag=1&isIABGlobal=false&hosts=&consentId=ffbda86e-5844-463f-ad0a-46767ce0433c&interactionCount=1&isAnonUser=1&landingPath=NotLandingPage&groups=C0001%3A1%2CC0003%3A1%2CC0002%3A1%2CC0004%3A1&geolocation=ES%3BMD&AwaitingReconsent=false; OptanonAlertBoxClosed=2025-04-09T03:06:14.269Z; yomvi_permisos=1744167994856; mplus_auth=webplayer
Upgrade-Insecure-Requests: 1
Sec-Fetch-Dest: document
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: none
Sec-Fetch-User: ?1
Priority: u=0, i
Pragma: no-cache
Cache-Control: no-cache
content-length: 0
~~~

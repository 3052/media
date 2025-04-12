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

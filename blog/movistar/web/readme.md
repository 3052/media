# movistar

https://justwatch.com/es/proveedor/movistar-plus-plus-ficcion-total

~~~
url = http://wl.movistarplus.es/ficha/?id=3427440
monetization = FLATRATE
country = Spain
~~~

smart proxy blocks login - proxy seller works. this is it:

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

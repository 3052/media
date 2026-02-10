# sky

~~~
url = https://show.sky.ch/de/filme/2035/a-knights-tale
monetization = FLATRATE
country = Switzerland
~~~

https://justwatch.com/ch/Anbieter/sky

## phone

login is protected:

~~~go
cookie: aws-waf-token=2e86b681-4c6d-40cd-9856-9ec0780664e5:HAoAkAsSO8kGAAAA:wW...
~~~

if you drop the Amazon request or the Cookie, the login fails

https://apkfab.com/it/sky/homedia.sky.sport

~~~
config.armeabi_v7a.apk
~~~

so need Android 9. this request:

~~~
GET https://clientapi.prd.sky.ch/stream/2035/MOVIE HTTP/2.0
devicecode: ANDROID_INAPP
authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiIxOTM4NDQ...
~~~

has no geo block, but the authorization only lasts five minutes, and the
refresh call is geo blocked, so the web client is better. also request above
does not accept these:

~~~
_ASP.NET_SessionId_ fail
SkyCheeseCake fail
sky-auth-token fail
~~~

## tv

if you request TV app, phone app is returned:

~~~
> play -i homedia.sky.sport -leanback
details[8] = 0 USD
details[13][1][4] = 1.18.1.142
details[13][1][16] = Jan 21, 2025
details[13][1][17] = APK APK APK
details[13][1][82][1][1] = 8.0 and up
details[15][18] = https://support.sky.ch/hc/en-us/articles/9520105066140
downloads = 468.36 thousand
name = Sky
size = 35.09 megabyte
version code = 584

> play -i homedia.sky.sport
details[8] = 0 USD
details[13][1][4] = 1.18.1.142
details[13][1][16] = Jan 21, 2025
details[13][1][17] = APK APK APK
details[13][1][82][1][1] = 8.0 and up
details[15][18] = https://support.sky.ch/hc/en-us/articles/9520105066140
downloads = 468.36 thousand
name = Sky
size = 35.09 megabyte
version code = 584
~~~

## web

https://github.com/sunsettrack4/plugin.video.skych/blob/master/addon.py

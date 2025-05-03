# canal

according to this:

https://boutique.canalplus.com/plan-de-site

the URL is here:

https://play.google.com/store/apps/details?id=com.canal.android.canal

according to this:

https://canalplus.cz

the URL is here:

https://play.google.com/store/apps/details?id=com.canalplus.canalplus

## com.canal.android.canal

https://apkmirror.com/apk/groupe-canal/mycanal-vos-programmes-en-live-ou-en-replay

create Android 7 device. install system certificate. need residential proxy

~~~
adb shell am start -a android.intent.action.VIEW `
-d https://www.canalplus.cz/stream/film/argylle-tajny-agent
~~~

then:

~~~xml
action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.scheme = http
data.scheme = https
data.host = www.mycanal.fr

action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.scheme = tvchannels
data.host = com.canal.android.canal

action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.scheme = http
data.scheme = https
data.host = www.canalplus.com

action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.scheme = https
data.host = mycanal.onelink.me
data.pathPrefix = /1424707377

action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.scheme = https
data.scheme = http
data.host = mycan.al

action.name = android.intent.action.VIEW
category.name = android.intent.category.BROWSABLE
category.name = android.intent.category.DEFAULT
data.host = com.canal.android.canal
data.pathPrefix = /content
data.pathPrefix = /startapp
data.scheme = mycanaltvlauncher
~~~

## com.canalplus.canalplus

<https://apk.gold/download?file_id=3155967/canalplus-app>

~~~
The APK failed to install.<br/> Error: INSTALL_FAILED_MISSING_SPLIT: Missing
split for com.canalplus.canalplus
~~~

https://apkpure.com/canal-app/com.canalplus.canalplus/download/11.3

create Android 7 device

1. select your region, cesko
2. submit
3. english
4. sign in
5. email
6. password
   - adb shell input text PASSWORD
7. log in

~~~xml
com.canalplus.canalplus.xml
action.name = android.intent.action.VIEW
category.name = android.intent.category.DEFAULT
category.name = android.intent.category.BROWSABLE
data.scheme = @string/application_scheme
data.scheme = @string/additional_scheme
~~~

then:

~~~
resources\res\values\strings.xml
<string name="additional_scheme">m7cp</string>
<string name="application_scheme">https</string>
~~~

if you do this:

~~~
adb shell am start -a android.intent.action.VIEW `
-d https://www.canalplus.cz/stream/film/argylle-tajny-agent
~~~

it opens the app, but it stays stuck on the splash screen. same for this:

~~~
adb shell am start -a android.intent.action.VIEW `
-d m7cp://www.canalplus.cz/stream/film/argylle-tajny-agent
~~~

also no requests are made. this is it:

~~~
POST https://tvapi-hlm2.solocoo.tv/v1/assets/6kJR8duRW27PS1KJuZKRkzLHw3-Z1OVNKQH5kN56/play HTTP/2.0
authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0di5zb2xvY29vLmF...
content-type: application/json; charset=UTF-8

{"connectionType":"MobileData","player":{"capabilities":{"drmSystems":["Widevine"],"drmSystemsWithInfo":[{"level":"L3","system":"Widevine","systemId":"4464"}],"embeddedSubtitles":true,"hdcpLevel":"Unprotected","maxHdcpLevel":"Unprotected","mediaTypes":["DASH"],"smartLib":true},"name":"AndroidXMedia3","version":"1.2.0"}}
~~~

from:

~~~
GET https://tvapi-hlm2.solocoo.tv/v1/assets?query=episodes%2CWY31FZFpXEpd6-TwSvwG29-SvzpxbpZ-RtzixkdY%2Cseason%2C2&sort=seasonepisode HTTP/2.0
authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0di5zb2xvY29vLmF...
~~~

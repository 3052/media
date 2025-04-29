# canal

no method to get object ID:

- https://github.com/dut-iptv/dut-iptv/blob/master/plugin.video.canaldigitaal/resources/lib/api.py
- https://github.com/add-ons/plugin.video.tvvlaanderen/blob/master/resources/lib/solocoo/asset.py

## web

~~~py
from mitmproxy import http

data = '''
console.log('_0xb40f61', _0xb40f61);
console.log('_0xffbd34', _0xffbd34);
console.log('_0x44b887', _0x44b887);
console.log('_0x5bdf04', _0x5bdf04);
console.log('_0x5430bb', _0x5430bb);
console.log('_0x4ab337', _0x4ab337);
return'Client'''

def response(f: http.HTTPFlow) -> None:
   if f.request.path.startswith('/static/js/main.4c582264.js'):
      f.response.text = f.response.text.replace("return'Client", data)
~~~

## com.canalplus.canalplus

https://play.google.com/store/apps/details?id=com.canalplus.canalplus

Updated on Apr 7, 2025

500K+ Downloads

<https://apk.gold/download?file_id=3155967/canalplus-app>

~~~
The APK failed to install.<br/> Error: INSTALL_FAILED_MISSING_SPLIT: Missing
split for com.canalplus.canalplus
~~~

https://apkpure.com/canal-app/com.canalplus.canalplus/download/11.3

1. select your region, cesko
2. submit
3. english
4. sign in
5. email
6. password
   - adb shell input text PASSWORD
7. log in

~~~
adb shell am start -a android.intent.action.VIEW `
-d https://www.canalplus.cz/stream/film/argylle-tajny-agent
~~~

next:

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

no deep link, so we will need to parse HTML

## com.canal.android.canal

https://play.google.com/store/apps/details?id=com.canal.android.canal

Updated on Apr 24, 2025

10M+ Downloads

https://apkmirror.com/apk/groupe-canal/mycanal-vos-programmes-en-live-ou-en-replay

~~~
SdkVersion='35' fail
SdkVersion='34' fail
SdkVersion='33' ABI ABI
SdkVersion='32'
SdkVersion='31' fail
SdkVersion='30'
SdkVersion='29'
SdkVersion='28'
SdkVersion='27' fail
SdkVersion='26'
SdkVersion='25'
sdkVersion:'24' fail
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

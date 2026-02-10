# NBC

## Android

only requires Android 5:

~~~
> play -a com.nbcuni.nbc
requires: 5.0 and up
version code: 2000004392
version name: 9.4.1
~~~

which means we can install up to Android 6 and install user certificate. If we
do that and start the app, we get:

> Unable to Connect
>
> Please check your internet connection and try again.

If we FORCE STOP the app, set No proxy and try again, it works. next we can try
system certificate. same result. this result also repeats with older versions
back to:

~~~
package: name='com.nbcuni.nbc' versionCode='2000003396' versionName='9.0.0'
compileSdkVersion='31' compileSdkVersionCodename='12'
~~~

with older versions:

~~~
package: name='com.nbcuni.nbc' versionCode='2000003368' versionName='7.35.1'
compileSdkVersion='31' compileSdkVersionCodename='12'
~~~

we get:

> ACTION REQUIRED
>
> To keep watching your favorite shows, movies and more, you'll need to update
> to the latest version of the NBC App.

next we can try:

~~~
pip install frida-tools
~~~

download and extract server:

https://github.com/frida/frida/releases

for example:

~~~
frida-server-16.1.7-android-x86.xz
~~~

https://github.com/httptoolkit/frida-interception-and-unpinning

install app, then push server:

~~~
adb root
adb push frida-server-16.1.7-android-x86 /data/app/frida-server
adb shell chmod +x /data/app/frida-server
adb shell /data/app/frida-server
~~~

then:

~~~
frida -U `
-l config.js `
-l native-connect-hook.js `
-l android/android-system-certificate-injection.js `
-l android/android-certificate-unpinning.js `
-l android/android-certificate-unpinning-fallback.js `
-f com.nbcuni.nbc
~~~

with MITM Proxy, app just gets stuck on loading screen. also, instead of my own
script I tried this instead:

https://github.com/httptoolkit/httptoolkit-server/blob/8577f6a/src/interceptors/android/adb-commands.ts#L271-L294

my device does not have `/apex`, which seems to simplify the process. but even
still, MITM Proxy fails. I guess at this point its a MITM Proxy issue. I did
notice that I was using an older MITM Proxy version 9, but trying version 10
and even 8 didn't resolve the problem either. note I also tried using HTTP
Toolkit certificate with MITM Proxy with no luck. here is something. I had idea
to try older Frida, so I tried again with Frida 15.2.2. still failed, but a
different error:

~~~
javax.net.ssl.SSLHandshakeException:
java.security.cert.CertPathValidatorException: Trust anchor for certification
path not found.
at
com.android.org.conscrypt.OpenSSLSocketImpl.startHandshake(OpenSSLSocketImpl.java:328)
~~~

- https://github.com/mitmproxy/android-unpinner/discussions/20
- https://github.com/mitmproxy/mitmproxy/discussions/6477

using HTTP Toolkit instead of MITM Proxy works.

https://github.com/httptoolkit/frida-interception-and-unpinning/issues/56

## how to get `drmProxySecret` value?

if you visit a page such as this:

https://nbc.com/saturday-night-live/video/october-5-nate-bargatze/9000405633

you should see a request like this:

https://www.nbc.com/generetic/generated/generetic.971bc5df5bddfb45624e.js

in the response body, you should see something like this:

~~~json
{
  "coreVideo": {
    "drmProxyUrl": "https://drmproxy.digitalsvc.apps.nbcuni.com/drm-proxy/license",
    "drmProxySecret": "Whn8QFuLFM7Heiz6fYCYga7cYPM8ARe6"
  }
}
~~~

## how to get `hash` value?

in the code we have this:

~~~js
return ""
   .concat(const174_.drmProxyUrl, "/")
   .concat(const181_, "?time=")
   .concat(const182_, "&hash=")
   .concat(const184_, "&device=web")
   .concat(param178_2 ? "&keyId=".concat(param178_2) : "");
~~~

simplify:

~~~js
var91_ = param74_3(1358),
var92_ = param74_3.n(var91_),
let let179_ = arguments.length > 2 && void 0 !== arguments[2] ? arguments[2] : 0;
const182_ = new Date().getTime() + let179_,
const const181_ = param178_.toLowerCase(),
const183_ = const182_ + const181_,
const184_ = var92_()(const183_, const174_.drmProxySecret);
"&hash=" + const184_
~~~

we should print these:

~~~js
console.log('DRM', var92_().toString());
console.log('DRM', const183_);
console.log('DRM', const174_.drmProxySecret);
~~~

script:

~~~py
from mitmproxy import ctx, http

def response(f: http.HTTPFlow) -> None:
   if f.request.path.startswith('/generetic/generated/chunks/12.ff734ba67f44a707e609.js'):
      f.response.text = open('hello.js', 'r').read()
~~~

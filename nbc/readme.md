# NBC

## Android

only requires Android 5:

~~~
com.nbcuni.nbc
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

## drmProxySecret

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

## GraphQL

~~~py
from mitmproxy import http
import logging

# The target hash you want to intercept
TARGET_HASH = "ac2e08429df2a1b856252b7cc6b38c31975be9b2ad83352e5c74128fb9b8d0ac"

# A dummy hash to replace it with (must be same format/length if the server
# validates format). This is just the target hash with the last character
# changed to '0'
REPLACEMENT_HASH = "ac2e08429df2a1b856252b7cc6b38c31975be9b2ad83352e5c74128fb9b8d0a0"

def request(flow: http.HTTPFlow) -> None:
    # We only care about requests that have a body or query parameters
    # containing the hash
    if flow.request.method in ["POST", "GET"]:
        
        # 1. Check/Replace in URL Query Parameters (common for GET requests)
        if TARGET_HASH in flow.request.url:
            logging.info(f"[+] Found Target Hash in URL. Replacing...")
            flow.request.url = flow.request.url.replace(TARGET_HASH, REPLACEMENT_HASH)

        # 2. Check/Replace in Request Body (common for POST requests)
        if flow.request.content and TARGET_HASH.encode() in flow.request.content:
            logging.info(f"[+] Found Target Hash in Body. Replacing...")
            
            # Use replace explicitly on bytes
            # We assume the encoding is standard (UTF-8/ASCII)
            new_content = flow.request.content.replace(
                TARGET_HASH.encode(), 
                REPLACEMENT_HASH.encode()
            )
            flow.request.content = new_content
            
            # Update content-length header automatically
            flow.request.headers["content-length"] = str(len(flow.request.content))
~~~

## hash

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
from mitmproxy import http

def response(flow: http.HTTPFlow) -> None:
   if flow.request.path.startswith('/generetic/generated/chunks/12.ff734ba67f44a707e609.js'):
      flow.response.text = open('hello.js', 'r').read()
~~~

## platform

web:

~~~json
{
  "playbackUrl": "https://vod-lf-oneapp2-prd.akamaized.net/prod/nbc/gLU/RcQ/9000283422/1698569087378-MEWw4/cmaf/mpeg_cenc_2sec/master_cmaf.mpd",
  "type": "DASH"
}
~~~

android:

~~~json
{
  "playbackUrl": "https://vod-lf-oneapp2-prd.akamaized.net/prod/nbc/gLU/RcQ/9000283422/1698569087378-MEWw4/cmaf/mpeg_cenc/master_cmaf.mpd",
  "type": "DASH"
}
~~~

web seems to be the better option:

~~~diff
--- a/mpeg_cenc
+++ b/mpeg_cenc_2sec
-            <S t="0" d="180180" r="124"/>
-            <S t="22522500" d="60060" r="0"/>
-            <S t="22582560" d="180180" r="164"/>
-            <S t="52312260" d="120120" r="0"/>
-            <S t="52432380" d="180180" r="30"/>
-            <S t="58017960" d="120120" r="0"/>
-            <S t="58138080" d="180180" r="181"/>
-            <S t="90930840" d="120120" r="0"/>
-            <S t="91050960" d="180180" r="99"/>
-            <S t="109068960" d="60060" r="0"/>
-            <S t="109129020" d="180180" r="70"/>
-            <S t="121921800" d="114114" r="0"/>
+            <S t="0" d="60060" r="2030"/>
+            <S t="121981860" d="54054" r="0"/>
~~~

## programmingType

if you provide an invalid value:

https://lemonade.nbc.com/v1/vod/2410887629/9000283422?platform=web&programmingType=Clips

you get:

~~~json
{
  "code": 400,
  "error": "No AssetType/ProtectionScheme/Format Matches",
  "message": "Bad Request",
  "meta": {
    "mpxUrl": "https://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000283422?formats=mpeg-dash&assetTypes=NONE,2SEC,VBW&restriction=108697384&sig=006546b8badb29ec17d8fb9f393733900635f73326b600a1a1736563726574",
    "message": {
      "title": "No AssetType/ProtectionScheme/Format Matches",
      "description": "None of the available releases match the specified AssetType, ProtectionScheme, and/or Format preferences",
      "isException": true,
      "exception": "NoAssetTypeFormatMatches",
      "responseCode": "412"
    }
  }
}
~~~

these all return the same thing:

- http://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000283422?formats=mpeg-dash
- http://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000283422?formats=mpeg-dash&assetTypes=2SEC
- http://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000283422?formats=mpeg-dash&assetTypes=VBW

here is another option:

http://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000283422?formats=m3u

but again its not usable:

~~~
#EXT-X-SESSION-KEY:
   KEYFORMAT="com.apple.streamingkeydelivery",
   KEYFORMATVERSIONS="1",
   METHOD=SAMPLE-AES,IV=0xfcf13caf41cb4ec7bcc918872de873b9,
   URI="skd://fcf13caf41cb4ec7bcc918872de873b9"
~~~

this works though:

http://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000359946?switch=HLSServiceSecure

this option:

http://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000283422

fails:

> Invalid URL

regarding locked content:

https://nbc.com/john-wick/video/john-wick/3448375

this works:

https://lemonade.nbc.com/v1/vod/2304992029/3448375?platform=web&programmingType=Movie

these fail:

~~~
> curl link.theplatform.com/s/NnzsPC/media/guid/2304992029/3448375?formats=mpeg-dash
{
        "title": "Invalid Token",
        "description": "This content requires a valid, unexpired auth token.",
        "isException": true,
        "exception": "InvalidAuthToken",
        "responseCode": "403"
}

> curl link.theplatform.com/s/NnzsPC/media/guid/2304992029/3448375?switch=HLSServiceSecure
{
        "title": "Invalid Token",
        "description": "This content requires a valid, unexpired auth token.",
        "isException": true,
        "exception": "InvalidAuthToken",
        "responseCode": "403"
}
~~~

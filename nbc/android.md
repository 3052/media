# Android

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

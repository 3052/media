# Hulu

## Android

- https://play.google.com/store/apps/details?id=com.hulu.livingroomplus
- https://play.google.com/store/apps/details?id=com.hulu.plus

Create Android 6 device. Install user certificate. after entering password, if
you click LOG IN you get this:

> Hmm. Something’s up. Please check your internet settings and try again. If
> all’s fine on your end, visit our Help Center.

system certificate? same result. if we disable proxy? it works. next:

https://github.com/httptoolkit/frida-interception-and-unpinning

~~~
pip install frida-tools
~~~

download and extract server:

https://github.com/frida/frida/releases

for example:

~~~
frida-server-16.1.4-android-x86.xz
~~~

install app, then push server:

~~~
adb root
adb push frida-server-16.1.5-android-x86 /data/app/frida-server
adb shell chmod +x /data/app/frida-server
adb shell /data/app/frida-server
~~~

then:

~~~
frida -U `
-l config.js `
-l android/android-certificate-unpinning.js `
-f com.hulu.plus
~~~

this worked a couple of times:

~~~diff
+++ b/android/android-certificate-unpinning.js
@@ -223,7 +223,7 @@ const PINNING_FIXES = {

     'okhttp3.CertificatePinner': [
         {
-            methodName: 'check',
+            methodName: 'a',
             overload: ['java.lang.String', 'java.util.List'],
             replacement: () => NO_OP
         },
~~~

but it seems to be a race condition or something, as it only works sometimes.
like it might fail the first time, but then if I restart the app it will work.
not sure.

https://github.com/httptoolkit/frida-interception-and-unpinning/issues/55

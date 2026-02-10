# Hulu

1. signup.hulu.com/plans
2. select (hulu, get first month free)
3. email
   - mailsac.com
4. password
5. name
6. birthdate
7. gender
8. agree & continue
9. name on card
10. card number
11. expiration
12. cvc
13. zip code
14. submit

## Android

~~~
> play -a com.hulu.plus
requires: 5.0 and up
~~~

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

> Hulu requires Recaptcha for authentication so just passing account credentials
> is not possible without captcha solving services. To work around this, this
> tool simply takes a Hulu session cookie.

https://github.com/jkmartindale/hulu

is this true on Android? example request:

~~~
POST https://guide.hulu.com/guide/details?user_token=nk77TZQgj1xc245G... HTTP/2.0
x-hulu-user-agent: androidv4/5.3.0+12541-google/b3d7ca343f99384;OS_23,MODEL_Android SDK built for x86

{"eabs":["EAB::023c49bf-6a99-4c67-851c-4c9e7609cc1d::196861183::262714326"]}
~~~

`user_token` comes from here:

~~~
POST https://auth.hulu.com/v1/mobile/mfa/authenticate HTTP/2.0
x-hulu-user-agent: androidv4/5.3.0+12541-google/b3d7ca343f99384;OS_23,MODEL_Android SDK built for x86
content-type: application/x-www-form-urlencoded

code=941741&
friendly_name=Android%20-%20unknown%20Android%20SDK%20built%20for%20x86%20Android&
serial_number=b3d7ca343f99384&
token=83c42269-296c-47ea-ac62-023d02ef2a47
~~~

code is 2FA. token comes from here:

~~~
POST https://auth.hulu.com/v3/mobile/password/authenticate HTTP/2.0
x-hulu-user-agent: androidv4/5.3.0+12541-google/b3d7ca343f99384;OS_23,MODEL_Android SDK built for x86
content-type: application/x-www-form-urlencoded

additional_properties=%7B%22distro%22%3A%22google%22%2C%22device_platform%22%3A%22Android%22%2C%22device_type%22%3A%22mobile%22%2C%22app_version%22%3A%225.3.0%22%2C%22device_family%22%3A%22Android%22%2C%22build_number%22%3A%225012541%22%2C%22device_os%22%3A%22Android%20REL6.0%22%2C%22device_manufacturer%22%3A%22unknown%22%2C%22device_product%22%3A%22Android%20REL6.0%22%2C%22device_model%22%3A%22Android%20SDK%20built%20for%20x86%22%2C%22device_capabilities%22%3A%7B%22device%22%3A%7B%22hulu%3Aapp%3Aandroid%22%3A%225.3.0%22%2C%22hulu%3Aplatform%3Aandroid%3Agoogleplay%22%3A%2223%22%2C%22hulu%3Adevices%3Aunknown%3Aandroidsdkbuiltforx86%22%3A%22%22%7D%2C%22capabilities%22%3A%5B%22hulu%3Adcs%3Acapabilities%3Acompass%3Asite-map%22%2C%22hulu%3Adcs%3Acapabilities%3Aonboarding-person-collection%22%2C%22hulu%3Adcs%3Acapabilities%3Acompass%3Acompass-mvp%22%5D%7D%7D&
device_id=166&
friendly_name=Android%20-%20unknown%20Android%20SDK%20built%20for%20x86%20Android&
mobile_capabilities=telephony&
password=...&
screen_size=%7B%22width_pixels%22%3A1080%2C%22height_pixels%22%3A1794%2C%22width_pixel_density_in_inches%22%3A420%2C%22height_pixel_density_in_inches%22%3A420%7D&
serial_number=b3d7ca343f99384&
time_zone=America%2FChicago&
user_email=...&
recaptcha_type=android&
recaptcha_token=03AFcWeA6hFI1SkP4tWKM4l23acaqlu6la04aHYSxjAehgrfYIiJJocCXpLnkW...
~~~

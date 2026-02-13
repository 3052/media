# peacock

how to get SkyOTT key?

## account

1. https://privacy.com
2. New Card
3. Create Card
4. $10
5. Single-Use
6. Set $10 Spend Limit
7. https://peacocktv.com/plans/all-monthly
8. Monthly
9. GET PREMIUM
10. Email
11. Password
12. Re-enter Password
13. First Name
14. Last Name
15. Gender
16. Birth Year
17. Zip Code
18. CREATE ACCOUNT
19. first name
20. last name
21. address
22. city
23. state
24. zip
25. card number
26. expiry date
27. CVC
28. SUBSCRIBE
29. PAY NOW

## android

~~~
> play -a com.peacocktv.peacockandroid
downloads: 34.82 million
files: APK APK APK APK
name: Peacock TV: Stream TV & Movies
offered by: Peacock TV LLC
price: 0 USD
requires: 7.0 and up
size: 67.11 megabyte
updated on: Feb 7, 2024
version code: 124050214
version name: 5.2.14
~~~

https://play.google.com/store/apps/details?id=com.peacocktv.peacockandroid

If you start the app and Sign In, this request:

~~~
POST https://rango.id.peacocktv.com/signin/service/international HTTP/2.0
content-type: application/x-www-form-urlencoded
x-skyott-device: MOBILE
x-skyott-proposition: NBCUOTT
x-skyott-provider: NBCU
x-skyott-territory: US

userIdentifier=MY_EMAIL&password=MY_PASSWORD
~~~

will fail:

~~~
HTTP/2.0 429
~~~

You can fix this problem by removing this request header before starting the
app:

~~~
set modify_headers '/~u signin.service.international/x-skyott-device/'
~~~

Header needs to be removed from that request only, since other requests need the
header.

## tv

~~~
$env:path = 'C:\windows\system32'
.\rootAVD.bat system-images\android-29\android-tv\x86\ramdisk.img

adb shell mkdir -p /data/local/tmp/cacerts
adb push C:/Users/Steven/.mitmproxy/mitmproxy-ca-cert.pem /data/local/tmp/cacerts/c8750f0d.0
adb shell cp /system/etc/security/cacerts/* /data/local/tmp/cacerts

adb shell su -c 'mount -t tmpfs tmpfs /system/etc/security/cacerts'
adb shell su -c 'cp /data/local/tmp/cacerts/* /system/etc/security/cacerts'
adb shell su -c 'chcon u:object_r:system_file:s0 /system/etc/security/cacerts/*'

adb push frida-server-16.2.1-android-x86 /data/local/tmp/frida-server
adb shell chmod +x /data/local/tmp/frida-server
adb shell su -c /data/local/tmp/frida-server
~~~

https://issuetracker.google.com/issues/331256113

## web

you can get `x-skyott-usertoken` with web client via `/auth/tokens`, but it
need `idsession` cookie. Looks like Android is the same.

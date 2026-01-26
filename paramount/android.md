# android

intl:
https://apkmirror.com/apk/viacomcbs-streaming/paramount-4

old:
https://apkmirror.com/apk/viacomcbs-streaming/paramount-3

android TV:
https://apkmirror.com/apk/viacomcbs-streaming/paramount-2

us:
https://apkmirror.com/apk/viacomcbs-streaming/paramount

minimum version: Android 7 (24)

~~~
~/.android/avd/Pixel_XL.avd/emulator-user.ini
~~~

to:

~~~
window.x = 0
window.y = 0
~~~

then:

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

install system certificate

## How to get secret\_key?

~~~
com\cbs\app\androiddata\retrofit\util\RetrofitUtil.java
SecretKeySpec secretKeySpec = new SecretKeySpec(b("302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"), "AES");
~~~

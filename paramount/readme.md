# Paramount+

## link.theplatform.com

why do we need link.theplatform.com? because its the only anonymous option.
logged out web client is missing MPD:

https://paramountplus.com/shows/mayor-of-kingstown/video/xhr/episodes/page/0/size/18/xs/0/season/3

logged in the web client embeds MPD in HTML. with the below items, you need UK IP
and Android cookie, else MPD will be missing. web cookie fails. get Android
cookie:

~~~
POST https://www.paramountplus.com/apps-api/v2.0/androidphone/auth/login.json?at=ABDFhCKlU... HTTP/1.1
content-type: application/x-www-form-urlencoded

j_username=EMAIL&j_password=PASSWORD
~~~

https://www.intl.paramountplus.com/apps-api/v2.0/androidtv/video/cid/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ.json?locale=en-us&at=ABA3WXXZwgC0rQPN9WtWEUmpHsGCFJb6NP4tGjIFVLTuScgId9WA3LdC44hdHUJysQ0%3D

https://www.intl.paramountplus.com/apps-api/v3.0/androidtv/movies/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ.json?includeTrailerInfo=true&includeContentInfo=true&locale=en-us&at=ABDSbrWqqlbSWOrrXk8u9NaNdokPC88YiXcPvIFhPobM3a%2FJWNOSwiCMklwJDDJq4c0%3D

<https://paramountplus.com/apps-api/v3.1/androidphone/irdeto-control/anonymous-session-token.json?contentId=Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ&model=AOSP%20on%20IA%20Emulator&firmwareVersion=9&version=15.0.28&platform=PP_AndroidApp&locale=en-us&locale=en-us&at=ABBoPFHuygkRnnCKELRhypuq5uEAJvSiVATsY9xOASH88ibse11WuoLrFnSDf0Bv7EY%3D>

<https://www.intl.paramountplus.com/apps-api/v3.1/androidtv/irdeto-control/session-token.json?contentId=Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ&model=sdk_google_atv_x86&firmwareVersion=9&version=15.0.28&platform=PPINTL_AndroidTV&locale=en-us&at=ABBoPFHuygkRnnCKELRhypuq5uEAJvSiVATsY9xOASH88ibse11WuoLrFnSDf0Bv7EY%3D>

## android

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

## how to get app secret?

us:

~~~
sources\com\cbs\app\config\UsaMobileAppConfigProvider.java
~~~

- https://apkmirror.com/apk/cbs-interactive-inc/paramount
- https://play.google.com/store/apps/details?id=com.cbs.app

international:

~~~
sources/com/cbs/app/config/DefaultAppSecretProvider.java
~~~

- https://apkmirror.com/apk/viacomcbs-streaming/paramount-android-tv
- https://play.google.com/store/apps/details?id=com.cbs.ca

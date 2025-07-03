# Paramount+

## How to get secret\_key?

~~~
com\cbs\app\androiddata\retrofit\util\RetrofitUtil.java
SecretKeySpec secretKeySpec = new SecretKeySpec(b("302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"), "AES");
~~~

## link.theplatform.com

why do we need link.theplatform.com? because its the only anonymous option.
logged out web client is missing MPD:

https://paramountplus.com/shows/mayor-of-kingstown/video/xhr/episodes/page/0/size/18/xs/0/season/3

logged in the web client embeds MPD in HTML. with the below items, you need
`gb-lon-ovpn-001` and Android cookie, else MPD will be missing. web cookie
fails. get Android cookie:

~~~
POST https://www.paramountplus.com/apps-api/v2.0/androidphone/auth/login.json?at=ABDFhCKlU... HTTP/1.1
content-type: application/x-www-form-urlencoded

j_username=EMAIL&j_password=PASSWORD
~~~

<https://paramountplus.com/apps-api/v3.1/androidphone/irdeto-control/anonymous-session-token.json?contentId=Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ&model=AOSP%20on%20IA%20Emulator&firmwareVersion=9&version=15.0.28&platform=PP_AndroidApp&locale=en-us&locale=en-us&at=ABBoPFHuygkRnnCKELRhypuq5uEAJvSiVATsY9xOASH88ibse11WuoLrFnSDf0Bv7EY%3D>

https://www.intl.paramountplus.com/apps-api/v2.0/androidtv/video/cid/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ.json?locale=en-us&at=ABA3WXXZwgC0rQPN9WtWEUmpHsGCFJb6NP4tGjIFVLTuScgId9WA3LdC44hdHUJysQ0%3D

https://www.intl.paramountplus.com/apps-api/v3.0/androidtv/movies/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ.json?includeTrailerInfo=true&includeContentInfo=true&locale=en-us&at=ABDSbrWqqlbSWOrrXk8u9NaNdokPC88YiXcPvIFhPobM3a%2FJWNOSwiCMklwJDDJq4c0%3D

<https://www.intl.paramountplus.com/apps-api/v3.1/androidtv/irdeto-control/session-token.json?contentId=Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ&model=sdk_google_atv_x86&firmwareVersion=9&version=15.0.28&platform=PPINTL_AndroidTV&locale=en-us&at=ABBoPFHuygkRnnCKELRhypuq5uEAJvSiVATsY9xOASH88ibse11WuoLrFnSDf0Bv7EY%3D>

## apps

create Android 6 device. install user certificate. start video. after the
commercial you might get an error, try again.

## paramount phone us

- https://apkmirror.com/apk/cbs-interactive-inc/paramount
- https://play.google.com/store/apps/details?id=com.cbs.app

15.0.52:

~~~
sources\com\cbs\app\dagger\DataLayerModule.java
dataSourceConfiguration.setCbsAppSecret("4fb47ec1f5c17caa");

sources\com\cbs\app\dagger\SharedComponentModule.java
return new ci.a("{\"amazon_tablet\":\"c4abf90e3aa8131f\",\"amazon_mobile\":\"c1353af7ed0252d8\",\"google_mobile\":\"8c4edb1155a410e4\"}");
~~~

15.0.50:

~~~
sources\com\cbs\app\dagger\DataLayerModule.java
dataSourceConfiguration.setCbsAppSecret("cdaf0c8e254c4424");

sources\com\cbs\app\dagger\SharedComponentModule.java
return new di.a("{\"amazon_tablet\":\"c4abf90e3aa8131f\",\"amazon_mobile\":\"c1353af7ed0252d8\",\"google_mobile\":\"8c4edb1155a410e4\"}");
~~~

## paramount tv intl

- https://apkmirror.com/apk/viacomcbs-streaming/paramount-android-tv
- https://play.google.com/store/apps/details?id=com.cbs.ca

15.0.52:

~~~
sources\com\cbs\app\BuildConfig.java
put("swisscom", "6d5824edfa1e56d6");
put("timvision", "893b6cb2e9112879");
put("vodafone", "ace4afb584a31528");

sources\com\cbs\app\config\DefaultAppSecretProvider.java
return "e55edaeb8451f737";

sources\com\cbs\app\config\SetTopBoxAppSecretProvider.java
return "e55edaeb8451f737";
~~~

15.0.50:

~~~
sources\com\cbs\app\BuildConfig.java
put("swisscom", "2751aeb7d5379e3b");
put("timvision", "fc27c8f1cfda2b25");
put("vodafone", "ae7c0cbda94ff4d5");

sources\com\cbs\app\config\DefaultAppSecretProvider.java
return "0f56dbac9fee3a93";

sources\com\cbs\app\config\SetTopBoxAppSecretProvider.java
return "0f56dbac9fee3a93";
~~~

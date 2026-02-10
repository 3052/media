# apps

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

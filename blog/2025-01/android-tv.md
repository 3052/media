# Netflix (Android TV)

## 4K

Android 12+, which means no intercept

## 1080p x86

~~~
> play -i com.netflix.ninja -leanback
details[8] = 0 USD
details[13][1][4] =
details[13][1][16] =
details[13][1][17] =
details[13][1][82][1][1] = Varies with device
details[15][18] = http://www.netflix.com/privacy
downloads = 319.14 million
name = Netflix
size = 0 byte
version code = 0
~~~

2017:

https://apkmirror.com/apk/netflix-inc/netflix-android-tv

## 1080p x86 armeabi-v7a

https://play.google.com/store/apps/details?id=com.netflix.ninja

~~~
> play -i com.netflix.ninja -leanback -abi armeabi-v7a
details[8] = 0 USD
details[13][1][4] = 11.0.9 build 19837
details[13][1][16] = Nov 21, 2024
details[13][1][17] = APK APK
details[13][1][82][1][1] = 7.0 and up
details[15][18] = http://www.netflix.com/privacy
downloads = 319.14 million
name = Netflix
size = 105.06 megabyte
version code = 19837
~~~

create Television 1080p

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

fails with Android 12:

~~~
Failure [INSTALL_FAILED_NO_MATCHING_ABIS: Failed to extract native libraries,
res=-113]
~~~

works with Android 13, but:

https://issuetracker.google.com/issues/331256113

## 1080p armeabi-v7a armeabi-v7a

only Android 6 or 12 are available, which means either the app wont run, or it
will run but we wont be able to intercept

# Itv

your postcode
SW1A 1AA

even the media content is hard geo blocked

## ITVX

https://apkmirror.com/apk/itv-plc/itv-hub

~~~
air.ITVMobilePlayer
~~~

create Android 7 device. install system certificate

~~~
adb shell am start -a android.intent.action.VIEW `
-d https://www.itv.com/watch/goldeneye/18910
~~~

## ITVX (Android TV)

Android 8:

~~~
sdkVersion:'26'
~~~

result:

~~~
filter=((type=="video"&&DisplayHeight<=576)||(type!="video"))
~~~

- https://apkmirror.com/apk/itv-plc/itv-hub-your-tv-player-watch-live-on-demand-android-tv
- https://play.google.com/store/apps/details?id=air.ITVMobilePlayer

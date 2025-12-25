# disney

https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d

## android

- https://apkmirror.com/apk/disney/disney
- https://play.google.com/store/apps/details?id=com.disney.disneyplus

~~~
~/.android/avd/Pixel_XL.avd/emulator-user.ini
~~~

to:

~~~
window.x = 0
window.y = 0
~~~

https://stackoverflow.com/questions/78813238

then:

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

then:

~~~
adb shell input text HELLO
~~~

APK lies, you need at least Android 12 (S)

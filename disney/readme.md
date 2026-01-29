# disney

## devices

~~~
arcelik_smart_tv_optomax_nt_sl3000
~~~

## clientApiKey

~~~
curl -O https://www.disneyplus.com
~~~

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

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

then:

~~~
adb shell input text HELLO
~~~

APK lies, you need at least Android 12 (S)

## subscribe

1. disneyplus.com
2. Disney+ Hulu
   - select
3. email
   - no trial so just use my own
4. continue
5. confirm
6. password
7. yes, I would like to receive updates
   - false
8. agree & continue
9. birthdate
10. save & continue
11. name on card
12. card number
13. expiration date
14. security code
15. zip code
16. agree & subscribe

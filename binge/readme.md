# android

~~~
url = https://binge.com.au/movies/asset-contact-1997!7738
country = Australia
~~~

- https://play.google.com/store/apps/details?id=au.com.streamotion.ares.tv
- https://play.google.com/store/apps/details?id=au.com.streamotion.ares

~~~
> play -i au.com.streamotion.ares
details[8] = 0 USD
details[13][1][4] = 3.3.5
details[13][1][16] = Feb 25, 2025
details[13][1][17] = APK APK APK APK
details[13][1][82][1][1] = 8.0 and up
details[15][18] = https://help.binge.com.au/s/privacy-policy
downloads = 1.79 million
name = Binge
size = 29.42 megabyte
version code = 3030117
~~~

above is wrong - you need to use Android 9 or more. install system certificate

mullvad fails with all cities

smart proxy works

start proxy:

~~~
mitmproxy --upstream-auth USERNAME:PASSWORD `
--mode "upstream:http://au.smartproxy.com:30001"
~~~

set proxy in android studio

install certificate to device

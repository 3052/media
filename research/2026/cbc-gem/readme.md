# CBC

## Android client

https://play.google.com/store/apps/details?id=ca.cbc.android.cbctv

Create Android device API 23. Install user certificate.

~~~
adb shell am start -a android.intent.action.VIEW `
-d https://gem.cbc.ca/the-fall/s02e03
~~~

## How to create account?

Use Android client

## How to get `apiKey`?

~~~
sources\vd\g.java
private static final String loginRadiusProdKey =
"3f4beddd-2061-49b0-ae80-6f1f2ed65b37";
~~~

https://github.com/skylot/jadx

## X-Forwarded-For

Based on this:

<https://github.com/firehol/blocklist-ipsets/blob/master/geolite2_country/country_ca.netset>

The largest Canada block is:

~~~
99.224.0.0/11
~~~

So in MITM Proxy, press `O` to enter options. Move to `modify_headers` and
press Enter. Then press `a` to add a new entry:

~~~
/~q/X-Forwarded-For/99.224.0.0
~~~

Press Esc when finished, then `q`.

# Paramount+

## Android client

create Android 6 device. install user certificate. start video. after the
commercial you might get an error, try again.

US:

https://play.google.com/store/apps/details?id=com.cbs.app

INTL:

https://play.google.com/store/apps/details?id=com.cbs.ca

## try paramount+

1. paramountplus.com
2. TRY IT FREE
3. CONTINUE
4. make sure MONTHLY is selected, then under Essential click SELECT PLAN
5. if you see a bundle screen, click MAYBE LATER
6. CONTINUE
7. uncheck Yes, I would like to receive marketing
8. CONTINUE
9. START PARAMOUNT+
10. paramountplus.com
11. select a profile
12. under profile click Account
13. Cancel Subscription
14. CONTINUE TO CANCEL
15. I understand
16. YES, CANCEL
17. the first option
18. COMPLETE CANCELLATION

## How to get app\_secret?

~~~
sources\com\cbs\app\dagger\DataLayerModule.java
dataSourceConfiguration.setCbsAppSecret("a624d7b175f5626b");
~~~

## How to get secret\_key?

~~~
com\cbs\app\androiddata\retrofit\util\RetrofitUtil.java
SecretKeySpec secretKeySpec = new SecretKeySpec(b("302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"), "AES");
~~~

# Roku

## starz

1. therokuchannel.roku.com
2. create a free account
3. email
   - mailsac.com
4. first name
5. last name
6. password
7. month
8. date
9. year
10. gender
11. I agree to the terms
12. I'm not a robot
13. continue
14. see plans
15. confirm selection
16. add a payment method
17. card number
18. month
19. year
20. security code
21. address
22. city
23. state
24. zip
25. save
26. see plans
27. start free trial

## Android

com.roku.web.trc is The Roku Channel (Android TV):

https://play.google.com/store/apps/details?id=com.roku.web.trc

~~~
> play -i com.roku.web.trc -t
details[6] = Roku, Inc. & its affiliates
details[8] = 0 USD
details[13][1][4] = 1.1.15
details[13][1][16] = Oct 31, 2023
details[13][1][17] = APK APK APK APK
details[13][1][82][1][1] = 5.0 and up
downloads = 394.82 thousand
name = The Roku Channel
size = 9.93 megabyte
version code = 10017
~~~

create Android 9 device. install system certificate.

## client

~~~
/api/v1/account/activation
/api/v1/account/activation/{{.Code}}
/api/v1/account/token
/api/v3/playback
~~~

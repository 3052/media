# RTBF

1. auvio.rtbf.be
2. log in
3. start my registration
4. enable JavaScript
5. email
   - tempmail.best
6. password
7. password confirmation
8. first name
9. last name
10. date of birth
11. gender
12. postal code
13. country
   - United States (etats-unis)
14. I accept it
15. Je m'inscris (I want to register)
16. To validate your RTBF account, please access your email address. Receive an
   email with the latest information.

## android

https://play.google.com/store/apps/details?id=be.rtbf.auvio

create Android 8 device. install system certificate

~~~
adb shell am start -a android.intent.action.VIEW `
-d https://auvio.rtbf.be/emission/i-care-a-lot-27462
~~~

## client

~~~
/v2/customer/RTBF/businessunit/Auvio/entitlement
entitlement

/v2/customer/RTBF/businessunit/Auvio/auth/gigyaLogin
gigya login

/auvio/v1.23/pages
content

/accounts.login
login

/accounts.getJWT
jwt
~~~

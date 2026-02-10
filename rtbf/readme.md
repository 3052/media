# RTBF

1. rtbf.be/auvio
2. Se connecter (login)
3. Démarrer mon inscription (start my registration)
4. email
   - mailsac.com
5. password
6. confirm password
7. first name
8. last name
9. Date de naissance (date of birth)
10. gender
11. postal code
12. country
   - États-Unis
13. J'accepte le Contrat d’utilisation Mon RTBF (I accept the mon RTBF user
   agreement)
14. Je m'inscris (I want to register)
15. To validate your RTBF account, please access your email address. Receive an
   email with the latest information.

## android

https://play.google.com/store/apps/details?id=be.rtbf.auvio

~~~
> play -i be.rtbf.auvio -s
details[6] = RTBF
details[8] = 0 USD
details[13][1][4] = 3.1.35
details[13][1][16] = May 15, 2024
details[13][1][17] = APK
details[13][1][82][1][1] = 8.0 and up
downloads = 1.58 million
name = RTBF Auvio : direct et replay
size = 28.57 megabyte
version code = 1301035
~~~

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

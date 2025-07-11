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

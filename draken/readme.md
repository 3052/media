# draken

1. sweden VPN/proxy
2. drakenfilm.se/offers
3. Månadsprenumeration, Aktivera nu (monthly subscription, activate now)
4. Din e-postadress (your e-mail address)
   - mailsac.com
5. Gå vidare (go further)
6. password
7. Upprepa lösenord (repeat password)
8. Jag godkänner användarvillkoren (I accept the terms of use)
9. Skapa konto (create an account)
10. Gå vidare (go further)
11. card number
12. Utgångsdatum (expiry date)
13. CVC
14. Namn på kort (name on card)
15. Bekräfta förauktorisering (confirm pre-authorization)

## web

Magine-Accesstoken is hard coded in the JavaScript

## android

https://play.google.com/store/apps/details?id=com.draken.android

~~~
> play -i com.draken.android
details[6] = Draken Film
details[8] = 0 USD
details[13][1][4] = 4.5.0
details[13][1][16] = Feb 15, 2024
details[13][1][17] = APK APK APK
details[13][1][82][1][1] = 5.0 and up
downloads = 27.68 thousand
name = Draken Film
size = 14.65 megabyte
version code = 1707910466
~~~

https://apkcombo.com/draken-film/com.draken.android

Create Android 6 device. Install user certificate. Magine-AccessToken in the
APK:

~~~java
@Override // wc.b
public String a() {
  return "22cc71a2-8b77-4819-95b0-8c90f4cf5663";
}

@Override // wc.b
public String b() {
  return "ea6fc0bb-8352-4bd6-9c4d-040a2c478fe8";
}
~~~

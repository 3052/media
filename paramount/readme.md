# Paramount+

<blockquote>
   <img height=400
   src="https://static1.cbrimages.com/wordpress/wp-content/uploads/2020/01/Dr.-Manhattan-Featured-Image.jpg"
   />
</blockquote>

## try paramount+

1. paramountplus.com
2. try it free
3. continue
4. make sure monthly is selected, then under essential click select plan
5. if you see a bundle screen, click maybe later
6. continue
7. uncheck yes, i would like to receive marketing
8. continue
9. start paramount+

## How to get secret\_key?

~~~
com\cbs\app\androiddata\retrofit\util\RetrofitUtil.java
SecretKeySpec secretKeySpec = new SecretKeySpec(b("302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"), "AES");
~~~

## link.theplatform.com

why do we need link.theplatform.com? because its the only anonymous option.
logged out web client is missing MPD:

https://paramountplus.com/shows/mayor-of-kingstown/video/xhr/episodes/page/0/size/18/xs/0/season/3

logged in the web client embeds MPD in HTML. with the below items, you need
`gb-lon-ovpn-001` and Android cookie, else MPD will be missing. web cookie
fails. get Android cookie:

~~~
POST https://www.paramountplus.com/apps-api/v2.0/androidphone/auth/login.json?at=ABDFhCKlU... HTTP/1.1
content-type: application/x-www-form-urlencoded

j_username=EMAIL&j_password=PASSWORD
~~~

<https://paramountplus.com/apps-api/v3.1/androidphone/irdeto-control/anonymous-session-token.json?contentId=Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ&model=AOSP%20on%20IA%20Emulator&firmwareVersion=9&version=15.0.28&platform=PP_AndroidApp&locale=en-us&locale=en-us&at=ABBoPFHuygkRnnCKELRhypuq5uEAJvSiVATsY9xOASH88ibse11WuoLrFnSDf0Bv7EY%3D>

https://www.intl.paramountplus.com/apps-api/v2.0/androidtv/video/cid/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ.json?locale=en-us&at=ABA3WXXZwgC0rQPN9WtWEUmpHsGCFJb6NP4tGjIFVLTuScgId9WA3LdC44hdHUJysQ0%3D

https://www.intl.paramountplus.com/apps-api/v3.0/androidtv/movies/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ.json?includeTrailerInfo=true&includeContentInfo=true&locale=en-us&at=ABDSbrWqqlbSWOrrXk8u9NaNdokPC88YiXcPvIFhPobM3a%2FJWNOSwiCMklwJDDJq4c0%3D

<https://www.intl.paramountplus.com/apps-api/v3.1/androidtv/irdeto-control/session-token.json?contentId=Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ&model=sdk_google_atv_x86&firmwareVersion=9&version=15.0.28&platform=PPINTL_AndroidTV&locale=en-us&at=ABBoPFHuygkRnnCKELRhypuq5uEAJvSiVATsY9xOASH88ibse11WuoLrFnSDf0Bv7EY%3D>

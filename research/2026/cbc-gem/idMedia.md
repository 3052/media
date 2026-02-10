# idMedia

How to get `idMedia`? If we start with this:

https://gem.cbc.ca/the-fall/s02e03

we should return:

~~~
958273
~~~

using `device=web`, we get this request:

~~~
GET /ott/subscription/v2/gem/ads/vod/the-fall/s02e03?consumerConsentId=null&device=web&isMobileBrowser=false&playerWidth=1536&ppid=3T79NngogE4wtb0PdB_Rh34tgBDtvYHjq9KEhMAq1p&tcProvider=None HTTP/1.1
Host: services.radio-canada.ca
~~~

response contains `idMedia`, but is missing `title`:

~~~json
{
  "adParameters": {
    "vid": "958273"
  }
}
~~~

and this request:

https://services.radio-canada.ca/ott/catalog/v2/gem/show/the-fall/s02e03?device=web

but the `idMedia` is buried under 7 keys:

~~~
content	
0	
lineups	
1	
items	
3	
idMedia	958273
~~~

and this request:

~~~
GET https://gem.cbc.ca/_next/data/WDfYZ6y6xkdFmrbKb5BJy/the-fall/s02e03.json?show=the-fall&content=s02e03 HTTP/2.0
~~~

but the `idMedia` is buried under 9 keys:

~~~
pageProps
data	
content	
0	
lineups	
1	
items	
3	
idMedia	958273
~~~

and this request:

~~~
GET https://gem.cbc.ca/the-fall/s02e03 HTTP/2.0
~~~

but the `idMedia` is buried under 10 keys:

~~~
props	
pageProps	
data	
content	
0	
lineups	
1	
items	
3	
idMedia	958273
~~~

using `device=phone_android` we get this same request:

~~~
GET /ott/subscription/v2/gem/ads/vod/the-fall/s02e03?mobileDeviceId=1335a619-0fda-463a-b0f0-5701704b2688&ppid=1335a619-0fda-463a-b0f0-5701704b2688&tcProvider=none&device=phone_android HTTP/1.1
Host: services.radio-canada.ca
~~~

and this same request:

<https://services.radio-canada.ca/ott/catalog/v2/gem/show/the-fall/s02e03?device=phone_android>

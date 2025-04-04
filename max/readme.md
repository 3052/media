# max.com

1. auth.max.com/product
2. monthly
3. basic with ads
4. continue
5. email
   - mailsac.com
6. confirm email
7. password
8. first name
9. last name
10. create account
11. debit card
12. name on card
13. card number
14. expiration date
15. security code
16. zip code
17. state
18. pay now

~~~json
{
  "errors" : [ {
    "status" : "400",
    "code" : "payment.refused",
    "id" : "a04af317b64acbf60c853a3c3d2fc6c1",
    "detail" : ""
  } ]
}
~~~

## android

https://play.google.com/store/apps/details?id=com.wbd.stream

~~~
> play -i com.wbd.stream
details[8] = 0 USD
details[13][1][4] = 4.12.0.64
details[13][1][16] = Oct 11, 2024
details[13][1][17] = APK APK APK APK
details[13][1][82][1][1] = 5.0 and up
details[15][18] = https://www.max.com/privacy
downloads = 80.50 million
name = Max: Stream HBO, TV, & Movies
size = 112.78 megabyte
version code = 35352971
~~~

above is wrong, needs to be at least Android 7. install system certificate

## movie

https://justwatch.com/us/movie/heretic-2024

~~~
url = https://play.max.com/show/12199308-9afb-460b-9d79-9d54b5d2514c
monetization = FLATRATE
country = United States
~~~

correct URLs:

- https://max.com/movies/12199308-9afb-460b-9d79-9d54b5d2514c
- https://max.com/movies/heretic/12199308-9afb-460b-9d79-9d54b5d2514c
- https://play.max.com/video/watch/2a9b19c2-7dad-4f46-97f1-58c282824bd5/ea64405b-c32a-4ece-aeca-61ad47d6bfb0

~~~
https://play.max.com/video/watch/
2a9b19c2-7dad-4f46-97f1-58c282824bd5 VideoId
/
ea64405b-c32a-4ece-aeca-61ad47d6bfb0 EditId
~~~

all we care about is EditId. how do we get EditId from public URL? like this:

~~~
GET /cms/routes/movie/12199308-9afb-460b-9d79-9d54b5d2514c?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&page[items.size]=10 HTTP/2
Host: default.any-emea.prd.api.max.com
Cookie: st=eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJ0b2tlbi1iMWM5ZTc2MS1lMjQwLTRmZWMtOG...
~~~

result:

~~~
json.data.attributes.url = "/movie/12199308-9afb-460b-9d79-9d54b5d2514c";

json.included[119].relationships.show.data.id = "12199308-9afb-460b-9d79-9d54b5d2514c";
json.included[119].relationships.edit.data.id = "ea64405b-c32a-4ece-aeca-61ad47d6bfb0";
~~~

## tv-show

https://justwatch.com/us/tv-show/the-white-lotus

~~~
url = https://play.max.com/video/watch/28ae9450-8192-4277-b661-e76eaad9b2e6
monetization = FLATRATE
count = 3
country = United States
~~~

correct URLs:

- https://max.com/shows/14f9834d-bc23-41a8-ab61-5c8abdbea505
- https://max.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505

~~~
GET /cms/collections/227084608563650952176059252419027445293?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&pf[show.id]=14f9834d-bc23-41a8-ab61-5c8abdbea505&pf[seasonNumber]=2 HTTP/2
Host: default.any-emea.prd.api.max.com
Cookie: st=eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJ0b2tlbi1iMWM5ZTc2MS1lMjQwLTRmZWMtOG...
~~~

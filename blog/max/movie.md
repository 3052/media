# movie

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

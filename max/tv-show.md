# tv-show

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

all we care about is the EditIds. how do we get the EditIds from public URL?

~~~
GET /cms/routes/show/14f9834d-bc23-41a8-ab61-5c8abdbea505?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&page[items.size]=10 HTTP/2
Host: default.any-emea.prd.api.max.com
Cookie: st=eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJ0b2tlbi1iMWM5ZTc2MS1lMjQwLTRmZWMtOG...
~~~

season 1: 6 episodes
season 2: 7 episodes
season 3: 5 episodes

total 18

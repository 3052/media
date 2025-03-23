# tv-show

https://justwatch.com/us/tv-show/the-white-lotus

~~~
url = https://play.max.com/video/watch/28ae9450-8192-4277-b661-e76eaad9b2e6
monetization = FLATRATE
count = 3
country = Argentina
country = Belgium
country = Bolivia
country = Brazil
country = Bulgaria
country = Chile
country = Colombia
country = Costa Rica
country = Croatia
country = Czech Republic
country = Denmark
country = Ecuador
country = Finland
country = France
country = Guatemala
country = Hungary
country = Indonesia
country = Malaysia
country = Mexico
country = Netherlands
country = Norway
country = Peru
country = Philippines
country = Poland
country = Portugal
country = Romania
country = Singapore
country = Slovakia
country = Spain
country = Sweden
country = Taiwan
country = Thailand
country = United States
country = Venezuela
~~~

correct URLs:

- https://max.com/shows/14f9834d-bc23-41a8-ab61-5c8abdbea505
- https://max.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505

all we care about is EditId. how do we get EditId from public URL?

~~~
GET /cms/routes/show/14f9834d-bc23-41a8-ab61-5c8abdbea505?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&page[items.size]=10 HTTP/2
Host: default.any-emea.prd.api.max.com
Cookie: st=eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJ0b2tlbi1iMWM5ZTc2MS1lMjQwLTRmZWMtOG...
~~~

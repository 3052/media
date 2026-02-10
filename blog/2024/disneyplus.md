# Disney+

If we run this command:

~~~
justwatch -a /us/movie/martha-marcy-may-marlene
~~~

we get this result:

~~~
FLATRATE
   https://www.apps.disneyplus.com/ph/movies/martha-marcy-may-marlene/1770015428?filters=content_type%3Dmovie
   - Philippines
   https://disneyplus.bn5x.net/c/1206980/705874/9358?u=https%3A%2F%2Fwww.disneyplus.com%2Fmovies%2Fmartha-marcy-may-marlene%2F599QO3ufvTPZ&subId3=justappsvod
   - Australia
   - Canada
   - New Zealand
   - Singapore
~~~

Mullvad VPN does not offer Philippines. If I try the other URL from United
States, I just get a redirect to:

https://disneyplus.com

If I Forget About This Site, and try the other four locations, I get a redirect
to:

https://disneyplus.com/movies/martha-marcy-may-marlene/599QO3ufvTPZ

and:

~~~
(error): SDK.SdkSession.initialize() - LocationNotAllowedException
vendor.ce5bd9512bb606ba418d.js:2:1296786
~~~

found this:

> Solution for me was to change the timezone to manual and not automatic, in
> accordance of where your vpn is

https://reddit.com/r/mullvadvpn/comments/o1lset/-/ibs4j4m

same result. So we can still consider Disney+, but only for United States videos.

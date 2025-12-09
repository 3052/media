# pluto reductio ad absurdum

we start with this:

https://pluto.tv/on-demand/movies/6495eff09263a40013cf63a5

using the old API:

https://api.pluto.tv/v2/episodes/6495eff09263a40013cf63a5/clips.json

we get this:

<https://siloh.pluto.tv/735_Paramount_Pictures_LF/clip/6495efee9263a40013cf638d_Jack_Reacher/1080pDRM/20241115_113001/dash/0-end/main.mpd>

which is not valid:

~~~
> curl -i https://siloh.pluto.tv/735_Paramount_Pictures_LF/clip/6495efee9263a40013cf638d_Jack_Reacher/1080pDRM/20241115_113001/dash/0-end/main.mpd
HTTP/2 403
~~~

but if we change the scheme and host it works:

~~~
> curl -i http://silo-hybrik.pluto.tv.s3.amazonaws.com/735_Paramount_Pictures_LF/clip/6495efee9263a40013cf638d_Jack_Reacher/1080pDRM/20241115_113001/dash/0-end/main.mpd
HTTP/1.1 200 OK
~~~

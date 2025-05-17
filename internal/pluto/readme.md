# pluto

~~~
Skulduggery[CULT]
 — 
5/13/25, 10:20 PM
all good
ok it printed alot of qualities
~~~

## fixed

~~~
Skulduggery[CULT]
 — 
5/13/25, 10:24 PM
nice
its downloading so slow lmfao
~~~

and:

~~~
Skulduggery[CULT]
 — 
5/13/25, 10:18 PM
ahh ok

13:19:08 WriteFile C:/Users/lachl/media/Mpd
panic: open C:/Users/lachl/media/Mpd: The system cannot find the path specified.

goroutine 1 [running]:
main.main()
        C:/Users/lachl/Downloads/media-main/media-main/internal/pluto/pluto.go:28 +0x294

ig make the dir?
~~~

and:

~~~
Skulduggery[CULT]
 — 
5/13/25, 10:16 PM

 .\pluto.exe -a https://pluto.tv/us/on-demand/series/60be6b55ea905900134a266a/season/7 -c C:\Users\lachl\AppData\Local\vinetrimmer\devices\xiaomi\device_client_id_blob -p C:\Users\lachl\AppData\Local\vinetrimmer\devices\xiaomi\device_private_key
13:16:23 GET https://boot.pluto.tv/v4/start?appName=web&appVersion=9&clientID=9&clientModelNumber=9&drmCapabilities=widevine%3AL3&seriesIDs=60be6b55ea905900134a266a
13:16:24 GET https://api.pluto.tv/v2/episodes/60be6b55ea905900134a266a/clips.json
panic: runtime error: index out of range [0] with length 0

goroutine 1 [running]:
41.neocities.org/media/pluto.(*Vod).Clips(0xc000090240)
        C:/Users/lachl/Downloads/media-main/media-main/pluto/pluto.go:102 +0x578
main.(*flags).do_address(0xc00007c120)
        C:/Users/lachl/Downloads/media-main/media-main/internal/pluto/pluto.go:58 +0x5e
main.main()
        C:/Users/lachl/Downloads/media-main/media-main/internal/pluto/pluto.go:26 +0x27f

anything wrong with my command?
tried with -s aswell
~~~

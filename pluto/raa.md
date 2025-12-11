# pluto reductio ad absurdum

try 2:

~~~
https://boot.pluto.tv/v4/start?
appName=web&
appVersion=9.18.0-32296d47c9882754e360f1b28a33027c54cbad16&
clientID=e0292ffd-7e8b-4607-ab89-fcd441a74b40&
clientModelNumber=1.0.0&
seriesIDs=6495eff09263a40013cf63a5
~~~

then:

~~~
"stitcherDash": "https://cfd-v4-service-stitcher-dash-use1-1.prd.pluto.tv",
/v2
"path": "/stitch/dash/episode/6495eff09263a40013cf63a5/main.mpd"
?jwt=
"sessionToken": "eyJhbGciOiJIUzI1NiIsImtpZCI6ImQzYzBlZDU2LTIwYWItNDNmMC05Mzg0..."
~~~

then:

~~~json
{
  "statusCode": 400,
  "message": "Invalid params",
  "errors": [
    {
      "param": "deviceModel",
      "source": "jwt",
      "reason": "empty",
      "message": "deviceModel is empty"
    },
    {
      "param": "deviceVersion",
      "source": "jwt",
      "reason": "empty",
      "message": "deviceVersion is empty"
    },
    {
      "param": "deviceMake",
      "source": "jwt",
      "reason": "empty",
      "message": "deviceMake is empty"
    }
  ]
}
~~~

try 3:

~~~
https://boot.pluto.tv/v4/start?
appName=web&
appVersion=9.18.0-32296d47c9882754e360f1b28a33027c54cbad16&
clientID=e0292ffd-7e8b-4607-ab89-fcd441a74b40&
clientModelNumber=1.0.0&
deviceMake=firefox&
deviceModel=web&
deviceVersion=128.0.0&
seriesIDs=6495eff09263a40013cf63a5
~~~

576p:

~~~xml
<Representation id="3" width="1024" height="576" sar="1:1" mimeType="video/mp4" codecs="avc1.64001f" bandwidth="1531860">
   <SegmentTemplate timescale="15360" startNumber="1" media="video/576p-1600/$Number%05d$.m4s" initialization="video/576p-1600/init.mp4" presentationTimeOffset="0">
   </SegmentTemplate>
</Representation>
~~~

try 4:

~~~
https://boot.pluto.tv/v4/start?
appName=web&
appVersion=9.18.0-32296d47c9882754e360f1b28a33027c54cbad16&
clientID=e0292ffd-7e8b-4607-ab89-fcd441a74b40&
clientModelNumber=1.0.0&
deviceMake=firefox&
deviceModel=web&
deviceVersion=128.0.0&
drmCapabilities=widevine%3AL3&
seriesIDs=6495eff09263a40013cf63a5
~~~

576p also. try 5:

~~~
https://boot.pluto.tv/v4/start?
appName=web&
appVersion=9.18.0-32296d47c9882754e360f1b28a33027c54cbad16&
clientID=e0292ffd-7e8b-4607-ab89-fcd441a74b40&
clientModelNumber=1.0.0&
deviceMake=firefox&
deviceModel=web&
deviceVersion=128.0.0&
drmCapabilities=widevine%3AL1&
seriesIDs=6495eff09263a40013cf63a5
~~~

720p:

~~~xml
<AdaptationSet id="4" width="1280" height="720" frameRate="15360/512" segmentAlignment="true" par="16:9" contentType="video">
   <Representation id="7" sar="1:1" mimeType="video/mp4" codecs="avc1.64001f" bandwidth="2584954">
      <SegmentTemplate timescale="15360" startNumber="1" media="video/720p-2400/$Number%05d$.m4s" initialization="video/720p-2400/init.mp4" presentationTimeOffset="0">
      </SegmentTemplate>
   </Representation>
</AdaptationSet>
~~~

try 6:

~~~
https://boot.pluto.tv/v4/start?
appName=androidmobile&
appVersion=5.61.0&
clientID=b311dee6-18a0-44f6-9351-7eae6eb562ea_93351976e6a032ce&
clientModelNumber=Android%20SDK%20built%20for%20x86&
deviceMake=unknown&
deviceModel=Android%20SDK%20built%20for%20x86&
deviceVersion=6.0_23&
seriesIDs=6495eff09263a40013cf63a5
~~~

576p. try 7:

~~~
https://boot.pluto.tv/v4/start?
appName=androidmobile&
appVersion=5.61.0&
clientID=b311dee6-18a0-44f6-9351-7eae6eb562ea_93351976e6a032ce&
clientModelNumber=Android%20SDK%20built%20for%20x86&
deviceMake=unknown&
deviceModel=Android%20SDK%20built%20for%20x86&
deviceVersion=6.0_23&
drmCapabilities=widevine%3AL1&
seriesIDs=6495eff09263a40013cf63a5
~~~

720p. try 8:

~~~
https://boot.pluto.tv/v4/start?
appName=androidtv&
appVersion=5.53.0-leanback&
clientID=720234b6-ce56-462a-892a-cf0d80c51469_2a547545129d6564&
clientModelNumber=sdk_google_atv_x86&
deviceMake=unknown&
deviceModel=sdk_google_atv_x86&
deviceVersion=9_28&
seriesIDs=6495eff09263a40013cf63a5
~~~

576p. try 9:

~~~
https://boot.pluto.tv/v4/start?
appName=androidtv&
appVersion=5.53.0-leanback&
clientID=720234b6-ce56-462a-892a-cf0d80c51469_2a547545129d6564&
clientModelNumber=sdk_google_atv_x86&
deviceMake=unknown&
deviceModel=sdk_google_atv_x86&
deviceVersion=9_28&
drmCapabilities=widevine%3AL1&
seriesIDs=6495eff09263a40013cf63a5
~~~

1080p:

~~~xml
<AdaptationSet xmlns="urn:mpeg:dash:schema:mpd:2011" id="5" width="1920" height="1080" frameRate="15360/512" segmentAlignment="true" par="16:9" contentType="video">
   <Representation xmlns="urn:mpeg:dash:schema:mpd:2011" id="8" sar="1:1" mimeType="video/mp4" codecs="avc1.640028" bandwidth="4453909">
      <SegmentTemplate xmlns="urn:mpeg:dash:schema:mpd:2011" timescale="15360" startNumber="1" media="video/1080p-4500/$Number%05d$.m4s" initialization="video/1080p-4500/init.mp4" presentationTimeOffset="0">
      </SegmentTemplate>
   </Representation>
</AdaptationSet>
~~~

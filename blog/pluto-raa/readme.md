# pluto reductio ad absurdum

## old client

we start with this:

https://pluto.tv/on-demand/movies/6495eff09263a40013cf63a5

then:

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

and notably, it has 1080p:

~~~xml
<AdaptationSet id="5" contentType="video" width="1920" height="1080" frameRate="15360/512" segmentAlignment="true" par="16:9">
   <Representation id="8" bandwidth="4586756" codecs="avc1.640028" mimeType="video/mp4" sar="1:1">
      <SegmentTemplate timescale="15360" initialization="video/1080p-4500/init.mp4" media="video/1080p-4500/$Number%05d$.m4s" startNumber="1">
      </SegmentTemplate>
   </Representation>
</AdaptationSet>
~~~

## new web client

begin:

~~~
https://boot.pluto.tv/v4/start?
appName=web&
appVersion=9.18.0-32296d47c9882754e360f1b28a33027c54cbad16&
clientID=e0292ffd-7e8b-4607-ab89-fcd441a74b40&
clientModelNumber=1.0.0&
drmCapabilities=widevine%3AL3&
seriesIDs=6495eff09263a40013cf63a5
~~~

next:

https://cfd-v4-service-stitcher-dash-use1-1.prd.pluto.tv/v2/stitch/dash/episode/6495eff09263a40013cf63a5/main.mpd?jwt=eyJhbGciOiJIUzI1NiIsImtpZCI6ImQzYzBlZDU2LTIwYWItNDNmMC05Mzg0LTRiOTNhMmQyZTQ5MyIsInR5cCI6IkpXVCJ9.eyJzZXNzaW9uSUQiOiIyMWVlNjAyMS1kNTVjLTExZjAtOWZmYy00MjEwNTIyOTYwMTkiLCJjbGllbnRJUCI6IjcyLjE4MS4yMy4zOCIsImNpdHkiOiJJcnZpbmciLCJwb3N0YWxDb2RlIjoiNzUwNjMiLCJjb3VudHJ5IjoiVVMiLCJkbWEiOjYyMywiYWN0aXZlUmVnaW9uIjoiVVMiLCJkZXZpY2VMYXQiOjMyLjkxOTk5ODE2ODk0NTMxLCJkZXZpY2VMb24iOi05Ni45NzAwMDEyMjA3MDMxMiwicHJlZmVycmVkTGFuZ3VhZ2UiOiJlbiIsImRldmljZVR5cGUiOiJ3ZWIiLCJkZXZpY2VWZXJzaW9uIjoiMTI4LjAuMCIsImRldmljZU1ha2UiOiJmaXJlZm94IiwiZGV2aWNlTW9kZWwiOiJ3ZWIiLCJhcHBOYW1lIjoid2ViIiwiYXBwVmVyc2lvbiI6IjkuMTguMC0zMjI5NmQ0N2M5ODgyNzU0ZTM2MGYxYjI4YTMzMDI3YzU0Y2JhZDE2IiwiY2xpZW50SUQiOiJlMDI5MmZmZC03ZThiLTQ2MDctYWI4OS1mY2Q0NDFhNzRiNDAiLCJjbUF1ZGllbmNlSUQiOiIiLCJpc0NsaWVudEROVCI6ZmFsc2UsInVzZXJJRCI6IiIsImxvZ0xldmVsIjoiREVGQVVMVCIsInRpbWVab25lIjoiQW1lcmljYS9DaGljYWdvIiwic2VydmVyU2lkZUFkcyI6ZmFsc2UsImUyZUJlYWNvbnMiOmZhbHNlLCJmZWF0dXJlcyI6eyJhZExvYWQiOnsiY29ob3J0IjoiIn0sIm11bHRpQXVkaW8iOnsiZW5hYmxlZCI6dHJ1ZX0sIm11bHRpUG9kQWRzIjp7ImVuYWJsZWQiOnRydWV9LCJzZWFyY2hBUEkiOnsicXVlcnlWZXJzaW9uIjoiaHlicmlkIn0sInN0aXRjaGVySGxzTmciOnsiZGVtdXhlZEF1ZGlvIjoiaml0In0sInN0aXRjaGVySGxzTmdWbGwiOnsiZW5hYmxlZCI6dHJ1ZX0sInN0aXRjaGVySGxzTmdWb2QiOnsiZW5hYmxlZCI6dHJ1ZX19LCJmbXNQYXJhbXMiOnsiZndWY0lEMiI6ImUwMjkyZmZkLTdlOGItNDYwNy1hYjg5LWZjZDQ0MWE3NGI0MCIsImZ3VmNJRDJDb3BwYSI6ImUwMjkyZmZkLTdlOGItNDYwNy1hYjg5LWZjZDQ0MWE3NGI0MCIsImN1c3RvbVBhcmFtcyI6eyJmbXNfbGl2ZXJhbXBfaWRsIjoiIiwiZm1zX2VtYWlsaGFzaCI6IiIsImZtc19zdWJzY3JpYmVyaWQiOiIiLCJmbXNfaWZhIjoiIiwiZm1zX2lkZnYiOiIiLCJmbXNfdXNlcmlkIjoiZTAyOTJmZmQtN2U4Yi00NjA3LWFiODktZmNkNDQxYTc0YjQwIiwiZm1zX3ZjaWQydHlwZSI6InVzZXJpZCIsImZtc19yYW1wX2lkIjoiIiwiZm1zX2hoX3JhbXBfaWQiOiIiLCJmbXNfYmlkaWR0eXBlIjoiIiwiX2Z3XzNQX1VJRCI6IiIsImZtc19ydWxlaWQiOiIxMDAwMCwxMDAwOSJ9fSwiZHJtIjp7Im5hbWUiOiJ3aWRldmluZSIsImxldmVsIjoiTDMifSwiaXNzIjoiYm9vdC5wbHV0by50diIsInN1YiI6InByaTp2MTpwbHV0bzpkZXZpY2VzOlVTOlpUQXlPVEptWm1RdE4yVTRZaTAwTmpBM0xXRmlPRGt0Wm1Oa05EUXhZVGMwWWpRdyIsImF1ZCI6IioucGx1dG8udHYiLCJleHAiOjE3NjU0MTE2MTMsImlhdCI6MTc2NTMyNTIxMywianRpIjoiN2Y3MTg2M2YtMDkyZS00ZGUwLTkxYTMtOGViZGRmZTQyNjgzIn0.dmBaB4riaecyC-gumLwYfG8_SL4yTg69b4XTy6DTxgo

best quality is 576p:

~~~xml
<Representation id="0" width="1024" height="576" sar="1:1" mimeType="video/mp4"
codecs="avc1.64001f" bandwidth="1537432">
~~~

if you base64 decode the `jwt` you get this:

~~~json
{
  "sessionID": "21ee6021-d55c-11f0-9ffc-421052296019",
  "clientIP": "72.181.23.38",
  "city": "Irving",
  "postalCode": "75063",
  "country": "US",
  "dma": 623,
  "activeRegion": "US",
  "deviceLat": 32.91999816894531,
  "deviceLon": -96.97000122070312,
  "preferredLanguage": "en",
  "deviceType": "web",
  "deviceVersion": "128.0.0",
  "deviceMake": "firefox",
  "deviceModel": "web",
  "appName": "web",
  "appVersion": "9.18.0-32296d47c9882754e360f1b28a33027c54cbad16",
  "clientID": "e0292ffd-7e8b-4607-ab89-fcd441a74b40",
  "cmAudienceID": "",
  "isClientDNT": false,
  "userID": "",
  "logLevel": "DEFAULT",
  "timeZone": "America/Chicago",
  "serverSideAds": false,
  "e2eBeacons": false,
  "features": {
    "adLoad": {
      "cohort": ""
    },
    "multiAudio": {
      "enabled": true
    },
    "multiPodAds": {
      "enabled": true
    },
    "searchAPI": {
      "queryVersion": "hybrid"
    },
    "stitcherHlsNg": {
      "demuxedAudio": "jit"
    },
    "stitcherHlsNgVll": {
      "enabled": true
    },
    "stitcherHlsNgVod": {
      "enabled": true
    }
  },
  "fmsParams": {
    "fwVcID2": "e0292ffd-7e8b-4607-ab89-fcd441a74b40",
    "fwVcID2Coppa": "e0292ffd-7e8b-4607-ab89-fcd441a74b40",
    "customParams": {
      "fms_liveramp_idl": "",
      "fms_emailhash": "",
      "fms_subscriberid": "",
      "fms_ifa": "",
      "fms_idfv": "",
      "fms_userid": "e0292ffd-7e8b-4607-ab89-fcd441a74b40",
      "fms_vcid2type": "userid",
      "fms_ramp_id": "",
      "fms_hh_ramp_id": "",
      "fms_bididtype": "",
      "_fw_3P_UID": "",
      "fms_ruleid": "10000,10009"
    }
  },
  "drm": {
    "name": "widevine",
    "level": "L3"
  },
  "iss": "boot.pluto.tv",
  "sub": "pri:v1:pluto:devices:US:ZTAyOTJmZmQtN2U4Yi00NjA3LWFiODktZmNkNDQxYTc0YjQw",
  "aud": "*.pluto.tv",
  "exp": 1765411613,
  "iat": 1765325213,
  "jti": "7f71863f-092e-4de0-91a3-8ebddfe42683"
}
~~~

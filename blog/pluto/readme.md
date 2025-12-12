

https://cfd-v4-service-stitcher-dash-use1-1.prd.pluto.tv
/v2
/stitch/dash/episode/6495eff09263a40013cf63a5/main.mpd?
advertisingId=&appName=androidtv&appVersion=9&app_name=androidtv&clientDeviceType=0&clientID=9&clientModelNumber=9&country=US&deviceDNT=false&deviceId=9&deviceLat=32.9200&deviceLon=-96.9700&deviceMake=9&deviceModel=9&deviceType=9%2Candroid%2Cctv&deviceVersion=9&marketingRegion=US&serverSideAds=false&sessionID=03c62740-d708-11f0-8ce9-32642cdd1b37&sid=03c62740-d708-11f0-8ce9-32642cdd1b37&userId=&
jwt=eyJhbGciOiJIUzI1NiIsImtpZCI6ImVkNTRmNjYxLTlkNTItNDdkYi04ODgzLTA3ZTI0NjdhNTE2ZiIsInR5cCI6IkpXVCJ9.eyJzZXNzaW9uSUQiOiIwM2M2Mjc0MC1kNzA4LTExZjAtOGNlOS0zMjY0MmNkZDFiMzciLCJjbGllbnRJUCI6IjcyLjE4MS4yMy4zOCIsImNpdHkiOiJJcnZpbmciLCJwb3N0YWxDb2RlIjoiNzUwNjMiLCJjb3VudHJ5IjoiVVMiLCJkbWEiOjYyMywiYWN0aXZlUmVnaW9uIjoiVVMiLCJkZXZpY2VMYXQiOjMyLjkxOTk5ODE2ODk0NTMxLCJkZXZpY2VMb24iOi05Ni45NzAwMDEyMjA3MDMxMiwicHJlZmVycmVkTGFuZ3VhZ2UiOiJlbiIsImRldmljZVR5cGUiOiI5LGFuZHJvaWQsY3R2IiwiZGV2aWNlVmVyc2lvbiI6IjkiLCJkZXZpY2VNYWtlIjoiOSIsImRldmljZU1vZGVsIjoiOSIsImFwcE5hbWUiOiJhbmRyb2lkdHYiLCJhcHBWZXJzaW9uIjoiOSIsImNsaWVudElEIjoiOSIsImNtQXVkaWVuY2VJRCI6IiIsImlzQ2xpZW50RE5UIjpmYWxzZSwidXNlcklEIjoiIiwibG9nTGV2ZWwiOiJERUZBVUxUIiwidGltZVpvbmUiOiJBbWVyaWNhL0NoaWNhZ28iLCJzZXJ2ZXJTaWRlQWRzIjpmYWxzZSwiZTJlQmVhY29ucyI6ZmFsc2UsImZlYXR1cmVzIjp7ImFkTG9hZCI6eyJjb2hvcnQiOiIifSwiZW5hYmxlMTA4MHAiOnRydWUsIm11bHRpQXVkaW8iOnsiZW5hYmxlZCI6dHJ1ZX0sIm11bHRpUG9kQWRzIjp7ImVuYWJsZWQiOnRydWV9LCJzZWFyY2hBUEkiOnsicXVlcnlWZXJzaW9uIjoiaHlicmlkIn0sInN0aXRjaGVySGxzTmciOnsiZGVtdXhlZEF1ZGlvIjoiaml0In0sInN0aXRjaGVySGxzTmdWbGwiOnsiZW5hYmxlZCI6dHJ1ZX0sInN0aXRjaGVySGxzTmdWb2QiOnsiZW5hYmxlZCI6dHJ1ZX0sInZvZFByZXJvbGxBZHMiOnsiY29ob3J0IjoiIiwiZW5hYmxlZCI6dHJ1ZX19LCJmbXNQYXJhbXMiOnsiZndWY0lEMiI6IjkiLCJmd1ZjSUQyQ29wcGEiOiI5IiwiY3VzdG9tUGFyYW1zIjp7ImZtc19saXZlcmFtcF9pZGwiOiIiLCJmbXNfZW1haWxoYXNoIjoiIiwiZm1zX3N1YnNjcmliZXJpZCI6IiIsImZtc19pZmEiOiIiLCJmbXNfaWRmdiI6IiIsImZtc191c2VyaWQiOiI5IiwiZm1zX3ZjaWQydHlwZSI6InVzZXJpZCIsImZtc19yYW1wX2lkIjoiIiwiZm1zX2hoX3JhbXBfaWQiOiIiLCJmbXNfYmlkaWR0eXBlIjoiIiwiX2Z3XzNQX1VJRCI6IiIsImZtc19ydWxlaWQiOiIxMDAwMCwxMDAwOSJ9fSwiZHJtIjp7Im5hbWUiOiJ3aWRldmluZSIsImxldmVsIjoiTDEifSwiaXNzIjoiYm9vdC5wbHV0by50diIsInN1YiI6InByaTp2MTpwbHV0bzpkZXZpY2VzOlVTOk9RPT0iLCJhdWQiOiIqLnBsdXRvLnR2IiwiZXhwIjoxNzY1NTk1Mzg3LCJpYXQiOjE3NjU1MDg5ODcsImp0aSI6IjQ5NDI4NTIxLTg1ZDctNDk4YS1hNmFkLTQyMzYwZmE0OTQzNSJ9.zqng7NZWV8rp6KWoM0VIXc90HzQeKDVExo-S4ZmbAfU


if I make this request:

https://boot.pluto.tv/v4/start?appName=androidtv&appVersion=9&clientID=9&clientModelNumber=9&deviceMake=9&deviceModel=9&deviceVersion=9&drmCapabilities=widevine%3AL1&seriesIDs=6495eff09263a40013cf63a5

I get the attached response. then I build the MPD URL from these parts of the
response:

```
```

however the MPD has an issue, because the AdaptationSet ids do not match up
across Periods. I was thinking we can fix it by using another server from the
response:

```
"analytics": "https://sp.pluto.tv",
"api": "https://api.pluto.tv",
"campaigns": "https://service-campaigns-ga.prd.pluto.tv",
"carousel": "https://service-carousel-builder-ga.prd.pluto.tv",
"catalog": "https://service-media-catalog.clusters.pluto.tv",
"channels": "https://service-channels.clusters.pluto.tv",
"concierge": "https://service-concierge.clusters.pluto.tv",
"features": "https://service-features-ga.prd.pluto.tv",
"hub": "https://service-hub-builder-ga.prd.pluto.tv"
"pause": "https://service-ad-image-ga-weighted.prd.pluto.tv",
"preferences": "",
"recommender": "https://service-recommender.clusters.pluto.tv",
"search": "https://service-media-search.clusters.pluto.tv",
"users": "https://service-users.clusters.pluto.tv",
"vod": "https://service-vod.clusters.pluto.tv",
"watchlist": "https://service-watchlist-ga.prd.pluto.tv",
```

but I dont know if that is possible or how to do that

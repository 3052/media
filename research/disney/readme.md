# disney

https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d

this is it:

https://varnish32-c20-mia1-dss-vod-dssc-shield.tr.na.prod.dssedge.com/dvt1=exp=1766696744~url=%2Fgrn%2Fps01%2Fdisney%2Faa401a2b-b7f4-4c11-bf61-a3b06f9c974d%2F~psid=abcca325-a73e-4f76-a3da-9753d3f3cf66~aid=05b49544-06af-43a8-92cf-625412b17d6f~did=3a0085a4-5df9-4593-b8e0-7b807cfcae99~kid=k01~hmac=1c9c9fa8cad432001d1e44d1666db507a903afb746a12c8e7e5701d153c4f6f4/grn/ps01/disney/aa401a2b-b7f4-4c11-bf61-a3b06f9c974d/ctr-all-fb600154-a5e0-4125-ab89-01d627163485-b123e16f-c381-4335-bf76-dcca65425460.m3u8?r=720&v=1&hash=d16b8103ffbd2266d8b5ccc26a08fdc107938aba

from:

~~~
POST https://disney.playback.edge.bamgrid.com/v7/playback/ctr-regular HTTP/2.0
user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0
accept: application/vnd.media-service+json
accept-language: en-US,en;q=0.5
accept-encoding: identity
referer: https://www.disneyplus.com/
authorization: Bearer eyJ6aXAiOiJERUYiLCJraWQiOiJ0Vy10M2ZQUTJEN2Q0YlBWTU1rSkd4...
content-type: application/json
x-dss-feature-filtering: true
x-application-version: 5d5917f8
x-bamsdk-client-id: disney-svod-3d9324fc
x-bamsdk-platform: javascript/windows/firefox
x-bamsdk-version: 34.3
x-dss-edge-accept: vnd.dss.edge+json; version=2
x-request-yp-id: 63626081279ebe65eb50fb54
x-request-id: 7ff709c0-8c6f-492f-9056-d2a5833eae58
content-length: 887
origin: https://www.disneyplus.com
sec-fetch-dest: empty
sec-fetch-mode: cors
sec-fetch-site: cross-site
te: trailers

{
  "playback": {
    "attributes": {
      "resolution": {
        "max": [
          "1280x720"
        ]
      },
      "protocol": "HTTPS",
      "assetInsertionStrategies": {
        "point": "SGAI",
        "range": "SGAI"
      },
      "playbackInitiationContext": "ONLINE",
      "frameRates": [
        30
      ],
      "videoSegmentTypes": [
        "FMP4"
      ],
      "maxSlideDuration": "15_MIN",
      "promosSupported": true
    },
    "adTracking": {
      "limitAdTrackingEnabled": "NOT_SUPPORTED",
      "deviceAdId": "00000000-0000-0000-0000-000000000000",
      "privacyOptOut": "NO"
    },
    "tracking": {
      "playbackSessionId": "abcca325-a73e-4f76-a3da-9753d3f3cf66"
    }
  },
  "playbackId": "eyJtZWRpYUlkIjoiYWE0MDFhMmItYjdmNC00YzExLWJmNjEtYTNiMDZmOWM5NzRkIiwiYXZhaWxJZCI6ImNkNDkwZmE0LTBkMWYtNDU1ZS04ZGNiLWZmZmQ1MTY2NmMyMSIsImF2YWlsVmVyc2lvbiI6Mywic291cmNlSWQiOiJjZDQ5MGZhNC0wZDFmLTQ1NWUtOGRjYi1mZmZkNTE2NjZjMjEiLCJjb250ZW50VHlwZSI6InZvZCJ9",
  "allowedCreatives": [
    "VIDEO"
  ],
  "allowedInsertionVisuals": [
    "PROMO_TEXT",
    "PROMO_FULL_TEXT",
    "ON_SCREEN_RATING",
    "TITLE_TREATMENT",
    "ON_SCREEN_ADVISORY"
  ]
}
~~~

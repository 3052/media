# android

~~~
> play -a au.com.stan.and -s
downloads = 11.17 million
files = APK
name = Stan.
offered by = Stan Entertainment Pty Ltd
price = 0 USD
requires = 5.0 and up
size = 26.92 megabyte (26915531)
updated on = Feb 5, 2024
version code = 50929
version name = 4.31.1
~~~

- https://apkmirror.com/apk/stan-entertainment-pty-ltd/stan
- https://play.google.com/store/apps/details?id=au.com.stan.and

Create Android 6 device. Install user certificate.

## password 1

~~~
POST https://api.stan.com.au/login/v1/sessions/mobile/app/recapture?manufacturer=Android&os=Android&sdk=23&type=mobile&model=unknown&stanName=Stan-Android&stanVersion=4.31.1.50929 HTTP/2.0
content-type: application/x-www-form-urlencoded

email=EMAIL&password=PASSOWRD&recaptureToken=03AFcWeA7atS-iTh4EzBP5LaNOc7xBOOe90Ms1qHOvbePXb1UfknpdC_BrGcg-6kaWScydKg_R3UbZprWjuvTd4sUvgW_CBDXDgWNNqfY1e3gJ8XG72S0QtVGz8fx_DmDPe5QGK2vpEhrCKTXLdT2jK9LcboGvnbrCtcQzMmp0HOpkxieHsurYtWwhIUzRuzsU1JRnmlwQLuwm_SEH-BUTqfIT22csKO7u95vFjg1U1WiFNHfrr0mLsyD_bZcaVRoNr6D45Q7v_-ZHr7TzI_QDdSXDHPOeJXi0IlJiGwboVeo-RBb5Cj8fBkhatK9JkryeALjIrX4LQJt7D-MGy4RYTY201mJqZrBBPTAKJhc_J8geSOlyj6_sZgeJAf9QmFLQ2zLiDptDwJ7r0WPb1MifrBHWMbyc1_k8tcPO6nIFtYP1I3OSpQJxG320xZzeecEUBHu6NTaSwOXVKPT_zcYxA-ZnyhNyHxMIyVLMR-snAq5AepcViF6pveiPEH6_snvnxzHoFpqFQEG7
~~~

## password 2

~~~
POST https://api.stan.com.au/login/v1/sessions/mobile/app?manufacturer=Android&os=Android&sdk=23&type=mobile&model=unknown&stanName=Stan-Android&stanVersion=4.31.1.50929 HTTP/2.0
content-type: application/x-www-form-urlencoded

jwToken=eyJhbGciOiJIUzI1NiIsImtpZCI6InBpa2FjaHUiLCJ0eXAiOiJKV1QifQ.eyJleHAiOjE3MjE3NzQzODQsImp0aSI6IjM2NDk1YjBjOWJiOTRjZGI4Y2RmOWVkNTk3OWI1MTIwIiwiaWF0IjoxNzExNDA2Mzg0LCJyb2xlIjoidXNlciIsInVpZCI6ImUwNzUyOGZkM2I0NDRiMTQ4YTI0NmZmYjM5M2JlNjUyIiwic3RyZWFtcyI6ImhkIiwiY29uY3VycmVuY3kiOjMsInByb2ZpbGVJZCI6ImUwNzUyOGZkM2I0NDRiMTQ4YTI0NmZmYjM5M2JlNjUyIiwicHJvZmlsZU5hbWUiOiJzdGV2ZW4iLCJ0eiI6IkFtZXJpY2EvQ2hpY2FnbyIsImFwcCI6IlN0YW4tQW5kcm9pZCIsInZlciI6IjQuMzEuMS41MDkyOSIsImZlYXQiOjEyMDg5MTExMDR9.-otn_9ERL537Y3XM6pKjxxWfXx83x03MAxlsfGQdfYI&profileId=e07528fd3b444b148a246ffb393be652
~~~

## password 3

https://api.stan.com.au/concurrency/v1/streams?programId=4212387&drm=widevine&quality=sd&clientId=3f9f91a359959c29&format=dash&captions=ttml&jwToken=eyJhbGciOiJIUzI1NiIsImtpZCI6InBpa2FjaHUiLCJ0eXAiOiJKV1QifQ.eyJleHAiOjE3MjE3NzQzODgsImp0aSI6IjY1N2ExYTFkNWZkYzRhNzk5YWFlNmZlMjQwMjhiNDhmIiwiaWF0IjoxNzExNDA2Mzg4LCJyb2xlIjoidXNlciIsInVpZCI6ImUwNzUyOGZkM2I0NDRiMTQ4YTI0NmZmYjM5M2JlNjUyIiwic3RyZWFtcyI6ImhkIiwiY29uY3VycmVuY3kiOjMsInByb2ZpbGVJZCI6ImUwNzUyOGZkM2I0NDRiMTQ4YTI0NmZmYjM5M2JlNjUyIiwicHJvZmlsZU5hbWUiOiJzdGV2ZW4iLCJ0eiI6IkFtZXJpY2EvQ2hpY2FnbyIsImFwcCI6IlN0YW4tQW5kcm9pZCIsInZlciI6IjQuMzEuMS41MDkyOSIsImZlYXQiOjEyMDg5MTExMDR9.gTNITbe2UIoG5Y58IcB6LB2MuujZbWHrdmqGTbWqfN8&manufacturer=Android&os=Android&sdk=23&type=mobile&model=unknown&stanName=Stan-Android&stanVersion=4.31.1.50929&videoCodec=h264

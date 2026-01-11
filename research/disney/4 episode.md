# episode

https://disneyplus.com/play/60da223c-d2a0-411a-95c9-665a839371f9

request:

~~~
GET https://disney.api.edge.bamgrid.com/explore/v1.12/deeplink?action=playback&refId=60da223c-d2a0-411a-95c9-665a839371f9&refIdType=deeplinkId HTTP/2
authorization: Bearer eyJ6aXAiOiJERUYiLCJraWQiOiJ0Vy10M2ZQUTJEN2Q0YlBWTU1rSkd4...
~~~

response:

~~~
data.deeplink.actions[0].resourceId = "eyJtZWRpYUlkIjoiYWM2YzU5YmItMjZmNy00NGQwLWEwMjUtZTJmMGMzNjRkNThmIiwiYXZhaWxJZCI6IjgwZjVmY2RmLWM1OTgtNGRmZS1hNmVlLWRhOTIwMTUwYWU2NCIsImF2YWlsVmVyc2lvbiI6NTAsInNvdXJjZUlkIjoiODBmNWZjZGYtYzU5OC00ZGZlLWE2ZWUtZGE5MjAxNTBhZTY0IiwiY29udGVudFR5cGUiOiJ2b2QifQ==";
~~~

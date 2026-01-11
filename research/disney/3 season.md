# season

https://disneyplus.com/browse/entity-05eb6a8e-90ed-4947-8c0b-e6536cbddd5f

request:

~~~
GET https://disney.api.edge.bamgrid.com/explore/v1.12/season/f1b38ee5-310d-4be7-b8b3-6998d29d8e85?limit=15&offset=0 HTTP/2
authorization: Bearer eyJ6aXAiOiJERUYiLCJraWQiOiJ0Vy10M2ZQUTJEN2Q0YlBWTU1rSkd4...
~~~

response:

~~~
data.season.items[0].visuals.title = "The Bear";
data.season.items[0].visuals.seasonNumber = "1";
data.season.items[0].visuals.episodeNumber = "1";
data.season.items[0].visuals.episodeTitle = "System";
data.season.items[0].actions[0].resourceId = "eyJtZWRpYUlkIjoiYWM2YzU5YmItMjZmNy00NGQwLWEwMjUtZTJmMGMzNjRkNThmIiwiYXZhaWxJZCI6IjgwZjVmY2RmLWM1OTgtNGRmZS1hNmVlLWRhOTIwMTUwYWU2NCIsImF2YWlsVmVyc2lvbiI6NTAsInNvdXJjZUlkIjoiODBmNWZjZGYtYzU5OC00ZGZlLWE2ZWUtZGE5MjAxNTBhZTY0IiwiY29udGVudFR5cGUiOiJ2b2QifQ==";
~~~

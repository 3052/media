# disney

for a show here is what we want:

~~~js
data.page.containers[0].seasons[2].items[0].visuals.episodeTitle = "Electric Sheep";
data.page.containers[0].seasons[2].items[0].actions[0].resourceId = "eyJtZWRpYUlkIjoiZjM5OTk3MWEtYTc4NS00MWM4LTk2NzEtMjQyNThhYmNlY2UxIiwiYXZhaWxJZCI6ImYxMGU0ODFlLTliMzItNGRmYS04NTIyLTFhZmM0MjAzMWI5MSIsImF2YWlsVmVyc2lvbiI6MzMsInNvdXJjZUlkIjoiZjEwZTQ4MWUtOWIzMi00ZGZhLTg1MjItMWFmYzQyMDMxYjkxIiwiY29udGVudFR5cGUiOiJ2b2QifQ==";
~~~

for a movie:

https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d

we want this:

~~~
eyJtZWRpYUlkIjoiYWE0MDFhMmItYjdmNC00YzExLWJmNjEtYTNiMDZmOWM5NzRkIiwiYXZhaWxJZCI6ImNkNDkwZmE0LTBkMWYtNDU1ZS04ZGNiLWZmZmQ1MTY2NmMyMSIsImF2YWlsVmVyc2lvbiI6Mywic291cmNlSWQiOiJjZDQ5MGZhNC0wZDFmLTQ1NWUtOGRjYi1mZmZkNTE2NjZjMjEiLCJjb250ZW50VHlwZSI6InZvZCJ9
~~~

which is this:

~~~json
{
  "mediaId": "aa401a2b-b7f4-4c11-bf61-a3b06f9c974d",
  "availId": "cd490fa4-0d1f-455e-8dcb-fffd51666c21",
  "availVersion": 3,
  "sourceId": "cd490fa4-0d1f-455e-8dcb-fffd51666c21",
  "contentType": "vod"
}
~~~

this also works:

~~~json
{
  "mediaId": "aa401a2b-b7f4-4c11-bf61-a3b06f9c974d"
}
~~~

which is here:

~~~
data.page.actions[0].internalTitle = "The Roses - movie - mediaId:aa401a2b-b7f4-4c11-bf61-a3b06f9c974d";
data.page.actions[1].internalTitle = "The Roses - movie - mediaId:aa401a2b-b7f4-4c11-bf61-a3b06f9c974d";
~~~

better to just:

1. decode base64
2. decode JSON
3. show mediaId to user
4. user request mediaId
5. encode JSON
6. encode base64

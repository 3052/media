# movie

- https://disneyplus.com/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
- https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d

request:

~~~
GET https://disney.api.edge.bamgrid.com/explore/v1.12/page/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d?enhancedContainersLimit=1&limit=99 HTTP/2.0
authorization: Bearer eyJ6aXAiOiJERUYiLCJraWQiOiJ0Vy10M2ZQUTJEN2Q0YlBWTU1rSkd4...
~~~

response:

~~~
data.page.actions[1].visuals.displayText = "RESTART";
data.page.actions[1].resourceId = "eyJtZWRpYUlkIjoiYWE0MDFhMmItYjdmNC00YzExLWJmNjEtYTNiMDZmOWM5NzRkIiwiYXZhaWxJZCI6ImNkNDkwZmE0LTBkMWYtNDU1ZS04ZGNiLWZmZmQ1MTY2NmMyMSIsImF2YWlsVmVyc2lvbiI6Mywic291cmNlSWQiOiJjZDQ5MGZhNC0wZDFmLTQ1NWUtOGRjYi1mZmZkNTE2NjZjMjEiLCJjb250ZW50VHlwZSI6InZvZCJ9";
~~~

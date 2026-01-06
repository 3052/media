# disney

## play

how do I go from this:

https://disneyplus.com/play/d32df5dd-4487-4e4c-9649-20f4bb472923

to something sane? like this:

~~~
GET https://disney.api.edge.bamgrid.com/explore/v1.12/deeplink?action=playback&refId=d32df5dd-4487-4e4c-9649-20f4bb472923&refIdType=deeplinkId HTTP/2
authorization: Bearer eyJ6aXAiOiJERUYiLCJraWQiOiJ0Vy10M2ZQUTJEN2Q0YlBWTU1rSkd4...
~~~

## movie, show

request:

~~~
GET https://disney.api.edge.bamgrid.com/explore/v1.12/page/entity-21e70fbf-6a51-41b3-88e9-f111830b046c?disableSmartFocus=true&enhancedContainersLimit=15&limit=15 HTTP/2
authorization: Bearer eyJ6aXAiOiJERUYiLCJraWQiOiJ0Vy10M2ZQUTJEN2Q0YlBWTU1rSkd4...
~~~

response:

~~~
data.page.containers[0].seasons[0].visuals.name = "Season 1"
data.page.containers[0].seasons[0].id = "f2c23e4f-87ea-4cf1-82c5-bb9767783dae"
~~~

## season

request:

~~~
GET https://disney.api.edge.bamgrid.com/explore/v1.12/season/28226c6e-7c7b-4184-8a86-5dabc4b2832f?limit=15&offset=0 HTTP/2
authorization: Bearer eyJ6aXAiOiJERUYiLCJraWQiOiJ0Vy10M2ZQUTJEN2Q0YlBWTU1rSkd4...
~~~

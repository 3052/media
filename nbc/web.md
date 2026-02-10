# web client

> SORRY! WE'RE HAVING SOME TROUBLE.
>
> If the problem persists, please contact us and send us a note.

if we enable, we get this:

https://lemonade.nbc.com/v1/vod/2410887629/9000283422?platform=web&browser=other&programmingType=Full+Episode

if you change user agent:

~~~
general.useragent.override
~~~

to:

~~~
Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Safari/605.1.15
~~~

you get this instead:

https://lemonade.nbc.com/v1/vod/2410887629/9000283422?platform=web&browser=safari&programmingType=Full+Episode

but it doesnt help:

~~~
#EXT-X-SESSION-KEY:
   IV=0xfcf13caf41cb4ec7bcc918872de873b9,
   KEYFORMAT="com.apple.streamingkeydelivery",
   KEYFORMATVERSIONS="1",
   METHOD=SAMPLE-AES,
   URI="skd://fcf13caf41cb4ec7bcc918872de873b9"
~~~

also value is optional:

https://lemonade.nbc.com/v1/vod/2410887629/9000283422?platform=web&programmingType=Full+Episode

## platform

web:

~~~json
{
  "playbackUrl": "https://vod-lf-oneapp2-prd.akamaized.net/prod/nbc/gLU/RcQ/9000283422/1698569087378-MEWw4/cmaf/mpeg_cenc_2sec/master_cmaf.mpd",
  "type": "DASH"
}
~~~

android:

~~~json
{
  "playbackUrl": "https://vod-lf-oneapp2-prd.akamaized.net/prod/nbc/gLU/RcQ/9000283422/1698569087378-MEWw4/cmaf/mpeg_cenc/master_cmaf.mpd",
  "type": "DASH"
}
~~~

web seems to be the better option:

~~~diff
--- a/mpeg_cenc
+++ b/mpeg_cenc_2sec
-            <S t="0" d="180180" r="124"/>
-            <S t="22522500" d="60060" r="0"/>
-            <S t="22582560" d="180180" r="164"/>
-            <S t="52312260" d="120120" r="0"/>
-            <S t="52432380" d="180180" r="30"/>
-            <S t="58017960" d="120120" r="0"/>
-            <S t="58138080" d="180180" r="181"/>
-            <S t="90930840" d="120120" r="0"/>
-            <S t="91050960" d="180180" r="99"/>
-            <S t="109068960" d="60060" r="0"/>
-            <S t="109129020" d="180180" r="70"/>
-            <S t="121921800" d="114114" r="0"/>
+            <S t="0" d="60060" r="2030"/>
+            <S t="121981860" d="54054" r="0"/>
~~~

## programmingType

if you provide an invalid value:

https://lemonade.nbc.com/v1/vod/2410887629/9000283422?platform=web&programmingType=Clips

you get:

~~~json
{
  "code": 400,
  "error": "No AssetType/ProtectionScheme/Format Matches",
  "message": "Bad Request",
  "meta": {
    "mpxUrl": "https://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000283422?formats=mpeg-dash&assetTypes=NONE,2SEC,VBW&restriction=108697384&sig=006546b8badb29ec17d8fb9f393733900635f73326b600a1a1736563726574",
    "message": {
      "title": "No AssetType/ProtectionScheme/Format Matches",
      "description": "None of the available releases match the specified AssetType, ProtectionScheme, and/or Format preferences",
      "isException": true,
      "exception": "NoAssetTypeFormatMatches",
      "responseCode": "412"
    }
  }
}
~~~

these all return the same thing:

- http://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000283422?formats=mpeg-dash
- http://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000283422?formats=mpeg-dash&assetTypes=2SEC
- http://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000283422?formats=mpeg-dash&assetTypes=VBW

here is another option:

http://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000283422?formats=m3u

but again its not usable:

~~~
#EXT-X-SESSION-KEY:
   KEYFORMAT="com.apple.streamingkeydelivery",
   KEYFORMATVERSIONS="1",
   METHOD=SAMPLE-AES,IV=0xfcf13caf41cb4ec7bcc918872de873b9,
   URI="skd://fcf13caf41cb4ec7bcc918872de873b9"
~~~

this works though:

http://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000359946?switch=HLSServiceSecure

this option:

http://link.theplatform.com/s/NnzsPC/media/guid/2410887629/9000283422

fails:

> Invalid URL

regarding locked content:

https://nbc.com/john-wick/video/john-wick/3448375

this works:

https://lemonade.nbc.com/v1/vod/2304992029/3448375?platform=web&programmingType=Movie

these fail:

~~~
> curl link.theplatform.com/s/NnzsPC/media/guid/2304992029/3448375?formats=mpeg-dash
{
        "title": "Invalid Token",
        "description": "This content requires a valid, unexpired auth token.",
        "isException": true,
        "exception": "InvalidAuthToken",
        "responseCode": "403"
}

> curl link.theplatform.com/s/NnzsPC/media/guid/2304992029/3448375?switch=HLSServiceSecure
{
        "title": "Invalid Token",
        "description": "This content requires a valid, unexpired auth token.",
        "isException": true,
        "exception": "InvalidAuthToken",
        "responseCode": "403"
}
~~~

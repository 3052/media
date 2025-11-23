# hulu

currently I cannot get any video:

~~~
> hulu -a hulu.com/watch/f70dfd4d-dbfb-46b8-abb3-136c841bba11 `
> -f i=H264_1_CMAF_CENC_CTR_400K
12:47:23 MPD PSSH EhD7fAUBMndF07SIpyTgkdCDEhA+9IaL6FRH+JX7CDTdk5cf
12:47:23 POST https://hulu.playback.edge.bamgrid.com/widevine-hulu/v1/hulu/vod/obtain-license/61556664?deejay_device_id=210&nonce=341049817&signature=1763959640_305bdb57e3aabadec68173ca9f8d8d3f523372cc
panic: security-level.insufficient
~~~

but I can get any audio:

~~~
> hulu -a hulu.com/watch/f70dfd4d-dbfb-46b8-abb3-136c841bba11 `
> -f i=EC3_1_CMAF_CENC_CTR_256K
12:47:29 MPD PSSH EhD7fAUBMndF07SIpyTgkdCD
12:47:29 POST https://hulu.playback.edge.bamgrid.com/widevine-hulu/v1/hulu/vod/obtain-license/61556664?deejay_device_id=210&nonce=341049817&signature=1763959646_92cd50d8fd4c014764ff62dc6b2ca8cdb6ff226c
12:47:29 key 12b5853e5a54a79ab84aae29d8079283
~~~

video has this PSSH:

~~~json
{
   "keyIds": [
      "3ef4868be85447f895fb0834dd93971f",
      "fb7c0501327745d3b488a724e091d083"
   ]
}
~~~

audio:

~~~json
{
   "keyIds":[
      "fb7c0501327745d3b488a724e091d083"
   ]
}
~~~

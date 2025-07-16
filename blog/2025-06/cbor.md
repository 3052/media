# cbor

- http://ugorji.net/blog/go-codec-primer#decoding
- https://cs.opensource.google/search?q=cbor&sq=&ss=go

got more info - looks like the body is CBOR encoded:
https://github.com/CastagnaIT/plugin.video.netflix/issues/421
https://github.com/nohajc/netflix-mitm-proxy/issues/7
which makes sense because I saw this:

```json
"mslTransportConfiguration.cborEnabled": true,
```

I also see this which says GZIP/LZW/CBOR:

```
GZIPcLZWsmaxpayloadchunksize %'dCBOR
```

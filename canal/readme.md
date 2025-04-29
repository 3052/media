# canal

no method to get object ID:

- https://github.com/dut-iptv/dut-iptv/blob/master/plugin.video.canaldigitaal/resources/lib/api.py
- https://github.com/add-ons/plugin.video.tvvlaanderen/blob/master/resources/lib/solocoo/asset.py

## web

~~~py
from mitmproxy import http

data = '''
console.log('_0xb40f61', _0xb40f61);
console.log('_0xffbd34', _0xffbd34);
console.log('_0x44b887', _0x44b887);
console.log('_0x5bdf04', _0x5bdf04);
console.log('_0x5430bb', _0x5430bb);
console.log('_0x4ab337', _0x4ab337);
return'Client'''

def response(f: http.HTTPFlow) -> None:
   if f.request.path.startswith('/static/js/main.4c582264.js'):
      f.response.text = f.response.text.replace("return'Client", data)
~~~

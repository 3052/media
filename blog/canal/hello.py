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


'''
_0xb40f61 web.NhFyz4KsZ54
key

_0xffbd34 OXh0-pIwu3gEXz1UiJtqLPscZQot3a0q 
secret

_0x44b887 1742610635 
time

_0x5bdf04 iIDsg8tTosmXlcnH-QbvLfC7Knag3JlnQyTmrk91310 
b64 encoded sha256 hash of the post data

_0x5430bb https://m7cplogin.solocoo.tv/loginiIDsg8tTosmXlcnH-QbvLfC7Knag3JlnQyTmrk913101742610635
URL + data + time

_0x4ab337 NWsZ0zamliCYpEGNdmnvyCbnCSAE1DHNqEAFMaM1GDU // sig
'''

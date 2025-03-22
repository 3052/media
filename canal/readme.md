# canal

~~~
url = https://www.canalplus.cz/stream/film/argylle-tajny-agent
monetization = FLATRATE
country = Czech Republic
~~~

- https://github.com/dut-iptv/dut-iptv/blob/master/plugin.video.canaldigitaal/resources/lib/api.py
- https://github.com/add-ons/plugin.video.tvvlaanderen/blob/master/resources/lib/solocoo/asset.py

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

## smart proxy

1. register
2. register now
3. VPN
4. canal+ complete, order
5. e-mail
   - mailsac.com
6. first name
7. last name
8. I agree with general terms and conditions
9. K platbě (to payment)
10. card number
11. expiry date
12. security code
13. holder name
14. Zaplatit (pay)
15. password
   - min 8 znaků, velké písmeno, číslice, speciální znak (min 8 characters,
     uppercase letter, number, special character)
16. repeat the password
17. Dokončit registraci (complete registration)

## mullvad

~~~
> mullvad status
Connected to cz-prg-wg-201 in Prague, Czech Republic
Your connection appears to be from: Czech Republic, Prague. IPv4: 178.249.209.172

> curl -m 30 canalplus.cz
curl: (28) Failed to connect to canalplus.cz port 80 after 21050 ms: Timed out
~~~

## nord

~~~
Secure Connection Failed
An error occurred during a connection to www.canalplus.cz. PR_END_OF_FILE_ERROR
~~~

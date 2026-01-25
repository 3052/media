# criterion

https://apkmirror.com/apk/the-criterion-collection/the-criterion-channel

minimum version: Android 8

~~~
action.name = android.intent.action.VIEW
category.name = android.intent.category.DEFAULT
data.scheme = @string/appId

action.name = android.intent.action.VIEW
category.name = android.intent.category.DEFAULT
category.name = android.intent.category.BROWSABLE
data.scheme = @string/scheme
~~~

then:

~~~
resources\res\values\strings.xml
<string name="appId">com.criterionchannel</string>
<string name="scheme">vhxcriterionchannel</string>
~~~

these crash app:

~~~
vhxcriterionchannel:%2Fwildcat
vhxcriterionchannel:https%3A%2F%2Fwww.criterionchannel.com%2Fwildcat
vhxcriterionchannel:wildcat
~~~

these just open the app:

~~~
vhxcriterionchannel://%2Fwildcat
vhxcriterionchannel://https%3A%2F%2Fwww.criterionchannel.com%2Fwildcat
vhxcriterionchannel://wildcat
~~~

then:

~~~
~/.android/avd/Pixel_XL.avd/emulator-user.ini
~~~

to:

~~~
window.x = 0
window.y = 0
~~~

https://stackoverflow.com/questions/78813238

~~~
adb install-multiple (Get-ChildItem *.apk)
~~~

then:

~~~
adb shell am start -a android.intent.action.VIEW `
-d https://www.criterionchannel.com/121280-ritual
~~~

## how to get client\_id ?

~~~
resources\res\values\strings.xml
<string name="oauthClientId">9a87f110f79cd25250f6c7f3a6ec8b9851063ca156dae493bf362a7faf146c78</string>
<string name="oauthClientSecret">9U+2KvxiqPBckxKR9dlTXPDd/w+SMDmbNaZeVx7C3tnfUUYAXpkMr+WZWWKGTRn3I+3cTg+30O1cYtYJhZfrez1YFHHcr77b3z3AHzxKFQA=</string>

<string name="tvOauthClientId">349c8358994e8aa948a72593e2c9707b6b717969b98e98d8548f370096338933</string>
<string name="tvOauthClientSecret">8bHCsyIErImeOChN75a0+U0HBg01xQSIdfeqVrTfXKn5UQdF+vZ+M4K03/It/avCEfGjoH5DELaId1/1Zp1tAPEZ1Zf4RmynueOiBc42cPY=</string>
~~~

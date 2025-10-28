package paramount

import (
   "fmt"
   "testing"
)

func TestAppSecret(t *testing.T) {
   fmt.Println(app_secrets)
}

var app_secrets = []struct {
   id      string
   version string
   url     string
   file    string
   code    string
}{
   {
      id:      "com.cbs.app",
      version: "16.0.0",
      url:     "apkmirror.com/apk/cbs-interactive-inc/paramount",
      file:    "sources/com/cbs/app/dagger/DataLayerModule.java",
      code: `return new w60.e(apiEnvironmentTypeA, "9fc14cb03691c342", strInvoke,
      "9ab70ef0883049829a6e3c01a62ca547", "1e8ce303a2f647d4b842bce77c3e713b",
      null, zB, true, false, false, zB2, packageName, strB, 800, null);`,
   },
   {
      id:      "com.cbs.ca",
      version: "15.5.0",
      url:     "apkmirror.com/apk/viacomcbs-streaming/paramount-android-tv",
      file:    "sources/com/cbs/app/config/DefaultAppSecretProvider.java",
      code: `public final class DefaultAppSecretProvider implements g {
          @Override // q60.g
          public String invoke() {
              return "4a81a3c936f63cd5";
          }
      }`,
   },
}

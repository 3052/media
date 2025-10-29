package paramount

import (
   "fmt"
   "testing"
)

var apk1 = apk{
   id:   "com.cbs.ca",
   url:  "apkmirror.com/apk/viacomcbs-streaming/paramount-android-tv",
   file: "sources/com/cbs/app/config/DefaultAppSecretProvider.java",
   detail: []detail{
      {
         version: "16.0.0",
         code: `
         public final class DefaultAppSecretProvider implements e {
             @Override // f70.e
             public String invoke() {
                 return "6c68178445de8138";
             }
         }
         `,
      },
      {
         version: "15.5.0",
         code: `
         public final class DefaultAppSecretProvider implements g {
            @Override // q60.g
            public String invoke() {
               return "4a81a3c936f63cd5";
            }
         }
         `,
      },
   },
}

var apk0 = apk{
   id:   "com.cbs.app",
   url:  "apkmirror.com/apk/cbs-interactive-inc/paramount",
   file: "sources/com/cbs/app/dagger/DataLayerModule.java",
   detail: []detail{
      {
         version: "16.0.0",
         code: `
         return new w60.e(apiEnvironmentTypeA, "9fc14cb03691c342", strInvoke,
         "9ab70ef0883049829a6e3c01a62ca547",
         "1e8ce303a2f647d4b842bce77c3e713b", null, zB, true, false, false,
         zB2, packageName, strB, 800, null);
         `,
      },
   },
}

func TestAppSecret(t *testing.T) {
   fmt.Println(apk0, apk1)
}

type detail struct {
   code    string
   version string
}

type apk struct {
   detail []detail
   file   string
   id     string
   url    string
}

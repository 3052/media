package paramount

import (
   "os"
   "testing"
)

func TestPlayReady(t *testing.T) {
   token, err := ComCbsApp.At()
   if err != nil {
      t.Fatal(err)
   }
   session_var, err := token.playReady("wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q")
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache+"/paramount/PlayReady", []byte(session_var.LsSession), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

func TestLog(t *testing.T) {
   t.Log(apk0, apk1, location_tests)
}

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

var location_tests = []struct {
   location   []string
   url        string
}{
   {
      location:   []string{"USA"},
      url:        "https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
   },
   {
      location:   []string{"USA"},
      url:        "https://paramountplus.com/shows/video/esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
   },
   {
      url:        "https://paramountplus.com/movies/video/Y8sKvb2bIoeX4XZbsfjadF4GhNPwcjTQ",
      location:   []string{"Australia", "United Kingdom"},
   },
   {
      url:        "https://paramountplus.com/movies/video/3DcGhIoTusoQFB_YLGCtLvefraLxuZMJ",
      location: []string{
         "Brazil", "Canada", "Chile", "Colombia", "Mexico", "Peru",
      },
   },
}

package main

import (
   "fmt"
   "os"
   "os/exec"
   "path/filepath"
   "strings"
)

var (
   debug                     = false
   patchFstab                = false
   getUSBHPmodZ              = false
   ramdiskImg                = false
   restore                   = false
   installKernelModules      = false
   installPrebuiltKernelModules = false
   listAllAVDs               = false
   installApps               = false
   noParamsAtAll             = false
   copyAsAdmin               = false
   androidHome               = ""
   avdPathWithRdffile        = ""
   avdPath                   = ""
   rdfFile                   = ""
   rootAVD                   = ""
   magiskZip                 = ""
   bzFile                    = ""
   krFile                    = "kernel-ranchu"
   adbWorkdir                = "/data/data/com.android.shell"
   adbBaseDir                = ""
)

func main() {
   processArguments()

   getANDROIDHOME()

   if debug {
      fmt.Println("[!] We are in Debug Mode")
      fmt.Println("params=", os.Args[1:])
      fmt.Println("DEBUG=", debug)
      fmt.Println("PATCHFSTAB=", patchFstab)
      fmt.Println("GetUSBHPmodZ=", getUSBHPmodZ)
      fmt.Println("RAMDISKIMG=", ramdiskImg)
      fmt.Println("restore=", restore)
      fmt.Println("InstallKernelModules=", installKernelModules)
      fmt.Println("InstallPrebuiltKernelModules=", installPrebuiltKernelModules)
      fmt.Println("ListAllAVDs=", listAllAVDs)
      fmt.Println("InstallApps=", installApps)
      fmt.Println("NOPARAMSATALL=", noParamsAtAll)
      fmt.Println("COPYASADMIN=", copyAsAdmin)
   }

   if !installApps {
      if len(os.Args) < 2 {
         showHelpText()
         os.Exit(0)
      }
      if listAllAVDs {
         showHelpText()
         os.Exit(0)
      }
      if _, err := os.Stat(filepath.Join(androidHome, os.Args[1])); os.IsNotExist(err) {
         fmt.Printf("file %s not found\n", os.Args[1])
         os.Exit(0)
      }
   }

   fmt.Println("[*] Set Directorys")
   avdPathWithRdffile = filepath.Join(androidHome, os.Args[1])
   avdPath = filepath.Dir(avdPathWithRdffile)
   rdfFile = filepath.Base(avdPathWithRdffile)

   testWritePerm(rdfFile)

   if fi, err := os.Stat(avdPathWithRdffile); err == nil && fi.IsDir() {
      showHelpText()
      os.Exit(0)
   }

   if restore {
      restoreBackups()
      os.Exit(0)
   }

   testADB()

   rootAVD, _ = os.Getwd()
   magiskZip = filepath.Join(rootAVD, "Magisk.zip")

   bzFile = filepath.Join(rootAVD, "bzImage")

   if installApps {
      installapps()
      os.Exit(0)
   }

   adbBaseDir = adbWorkdir + "/Magisk"
   fmt.Println("[-] In any AVD via ADB, you can execute code without root in /data/data/com.android.shell")

   testADBWORKDIR()

   os.Chdir(rootAVD)

   fmt.Println("[*] Cleaning up the ADB working space")
   exec.Command("adb", "shell", "rm", "-rf", adbBaseDir).Run()

   fmt.Println("[*] Creating the ADB working space")
   exec.Command("adb", "shell", "mkdir", adbBaseDir).Run()

   fmt.Println("[*] looking for Magisk installer Zip")
   if _, err := os.Stat(magiskZip); os.IsNotExist(err) {
      fmt.Println("[-] Please download Magisk.zip file")
   } else {
      pushtoAVD(magiskZip, "")
   }

   initramfs := filepath.Join(rootAVD, "initramfs.img")

   if ramdiskImg {
      if !strings.Contains(strings.ToLower(rdfFile), "ramdisk") || !strings.HasSuffix(strings.ToLower(rdfFile), ".img") {
         fmt.Println("[!] please give a path to a ramdisk file")
         os.Exit(0)
      }

      createBackup(rdfFile)
      pushtoAVD(avdPathWithRdffile, "ramdisk.img")

      if installKernelModules {
         if _, err := os.Stat(initramfs); err == nil {
            pushtoAVD(initramfs, "")
         }
      }
   }

   fmt.Println("[-] Copy rootAVD Script into Magisk DIR")
   exec.Command("adb", "push", "rootAVD.sh", adbBaseDir).Run()

   fmt.Println("[-] run the actually Boot/Ramdisk/Kernel Image Patch Script")
   fmt.Println("[*] from Magisk by topjohnwu and modded by NewBit XDA")
   cmd := exec.Command("adb", "shell", "sh", adbBaseDir+"/rootAVD.sh", strings.Join(os.Args[1:], " "))
   if err := cmd.Run(); err == nil {
      if !debug {
         if ramdiskImg {
            if copyAsAdmin {
               pullfromAVD("ramdiskpatched4AVD.img", "")
               exec.Command("powershell.exe", "-Command", fmt.Sprintf("Start-Process -Wait cmd '/c copy \"%s\" \"%s\"' -Verb RunAs", filepath.Join(rootAVD, "ramdiskpatched4AVD.img"), avdPathWithRdffile)).Run()
               os.Remove("ramdiskpatched4AVD.img")
            } else {
               pullfromAVD("ramdiskpatched4AVD.img", avdPathWithRdffile)
            }

            pullfromAVD("Magisk.apk", filepath.Join(rootAVD, "Apps"))
            pullfromAVD("Magisk.zip", "")

            if installPrebuiltKernelModules {
               pullfromAVD(bzFile, "")
               installKernelModulesFunc()
            }

            if installKernelModules {
               installKernelModulesFunc()
            }

            fmt.Println("[-] Clean up the ADB working space")
            exec.Command("adb", "shell", "rm", "-rf", adbBaseDir).Run()

            installapps()

            fmt.Println("[-] Shut-Down and Reboot [Cold Boot Now] the AVD and see IF it worked")
            fmt.Println("[-] Root and Su with Magisk for Android Studio AVDs")
            fmt.Println("[-] Modded by NewBit XDA - Jan. 2021")
            fmt.Println("[*] Huge Credits and big Thanks to topjohnwu, shakalaca and vvb2060")
            shutDownAVD()
         }
      }
   }
}

func testADBWORKDIR() {
   fmt.Println("[*] Testing the ADB working space")
   out, _ := exec.Command("adb", "shell", "cd", adbWorkdir).CombinedOutput()
   if strings.Contains(string(out), "No such file or directory") {
      fmt.Printf("[^^!] %s is not available\n", adbWorkdir)
      os.Exit(1)
   }
   fmt.Printf("[^^!] %s is available\n", adbWorkdir)
}

func shutDownAVD() {
   out, _ := exec.Command("adb", "shell", "setprop", "sys.powerctl", "shutdown").CombinedOutput()
   if !strings.Contains(strings.ToLower(string(out)), "error") {
      fmt.Println("[-] Trying to shut down the AVD")
   }
   fmt.Println("[^^!] If the AVD doesnt shut down, try it manually^^!")
}

func installKernelModulesFunc() {
   if _, err := os.Stat(bzFile); err == nil {
      createBackup(krFile)
      fmt.Printf("[*] Copy %s (Kernel) into kernel-ranchu\n", bzFile)

      if copyAsAdmin {
         fmt.Println("[^^!] with elevated write permissions")
         exec.Command("powershell.exe", "-Command", fmt.Sprintf("Start-Process -Wait cmd '/c copy \"%s\" \"%s\"' -Verb RunAs", bzFile, filepath.Join(avdPath, krFile))).Run()
      } else {
         data, _ := os.ReadFile(bzFile)
         os.WriteFile(filepath.Join(avdPath, krFile), data, 0644)
      }

      os.Remove(bzFile)
      os.Remove(filepath.Join(rootAVD, "initramfs.img"))
   }
}

func pullfromAVD(src, dst string) {
   srcBase := filepath.Base(src)
   dstBase := filepath.Base(dst)
   var out []byte
   if dst != "" {
      if dstBase != "" {
         out, _ = exec.Command("adb", "pull", adbBaseDir+"/"+srcBase, dst).CombinedOutput()
      } else {
         out, _ = exec.Command("adb", "pull", adbBaseDir+"/"+srcBase).CombinedOutput()
      }
   } else {
      out, _ = exec.Command("adb", "pull", adbBaseDir+"/"+srcBase).CombinedOutput()
   }

   if !strings.Contains(strings.ToLower(string(out)), "error") {
      fmt.Printf("[*] Pull %s into %s\n", srcBase, dstBase)
      fmt.Printf("[-] %s\n", string(out))
   }
}

func pushtoAVD(src, dst string) {
   srcBase := filepath.Base(src)
   dstBase := filepath.Base(dst)
   var out []byte
   if dst == "" {
      fmt.Printf("[*] Push %s into %s\n", srcBase, adbBaseDir)
      out, _ = exec.Command("adb", "push", src, adbBaseDir).CombinedOutput()
   } else {
      fmt.Printf("[*] Push %s into %s/%s\n", srcBase, adbBaseDir, dstBase)
      out, _ = exec.Command("adb", "push", src, adbBaseDir+"/"+dstBase).CombinedOutput()
   }
   fmt.Printf("[-] %s\n", string(out))
}

func testWritePerm(file string) {
   tempFile := file + ".temp"
   fmt.Println("[*] Testing for write permissions")
   fmt.Println("[-] creating TEMPFILE File")
   data, err := os.ReadFile(filepath.Join(avdPath, file))
   if err != nil {
      copyAsAdmin = true
      fmt.Println("[^!] elevated write permissions are needed to access $ANDROID_HOME")
      return
   }
   err = os.WriteFile(filepath.Join(avdPath, tempFile), data, 0644)
   if err != nil {
      copyAsAdmin = true
      fmt.Println("[^!] elevated write permissions are needed to access $ANDROID_HOME")
   }

   if !copyAsAdmin {
      fmt.Println("[-] deleating TEMPFILE File")
      fmt.Println("[^^!] NO elevated write permissions are needed to access $ANDROID_HOME")
      os.Remove(filepath.Join(avdPath, tempFile))
   }
}

func createBackup(file string) {
   backupFile := file + ".backup"
   if _, err := os.Stat(filepath.Join(avdPath, backupFile)); os.IsNotExist(err) {
      fmt.Println("[*] creating Backup File")
      if copyAsAdmin {
         fmt.Println("[^^!] with elevated write permissions")
         exec.Command("powershell.exe", "-Command", fmt.Sprintf("Start-Process -Wait cmd '/c copy \"%s\" \"%s\"' -Verb RunAs", filepath.Join(avdPath, file), filepath.Join(avdPath, backupFile))).Run()
      } else {
         data, _ := os.ReadFile(filepath.Join(avdPath, file))
         os.WriteFile(filepath.Join(avdPath, backupFile), data, 0644)
      }
      if _, err := os.Stat(filepath.Join(avdPath, backupFile)); err == nil {
         fmt.Println("[-] Backup File was created")
      }
   } else {
      fmt.Println("[-] Backup exists already")
   }
}

func testADB() {
   fmt.Println("[-] Test IF ADB SHELL is working")
   out, _ := exec.Command("adb", "shell", "-n", "echo", "true").CombinedOutput()
   if strings.TrimSpace(string(out)) == "true" {
      fmt.Println("[-] ADB connection possible")
   } else {
      if strings.Contains(strings.ToLower(string(out)), "offline") {
         fmt.Println("[^^!] ADB device is offline")
         fmt.Println("[*] no ADB connection possible")
         os.Exit(1)
      }
      if strings.Contains(strings.ToLower(string(out)), "unauthorized") {
         fmt.Printf("[^^!] %s\n", string(out))
         fmt.Println("[*] no ADB connection possible")
         os.Exit(1)
      }
      // ... (Weitere Fehlerbehandlungen für ADB können hier hinzugefügt werden)
      fmt.Printf("[^^!] %s\n", string(out))
      fmt.Println("[*] no ADB connection possible")
      os.Exit(1)
   }
}

func restoreBackups() {
   files, _ := filepath.Glob(filepath.Join(avdPath, "*.backup"))
   for _, f := range files {
      originalFile := strings.TrimSuffix(f, ".backup")
      fmt.Printf("[^!] Restoring %s to %s\n", f, originalFile)
      if copyAsAdmin {
         fmt.Println("[^!] with elevated write permissions")
         exec.Command("powershell.exe", "-Command", fmt.Sprintf("Start-Process -Wait cmd '/c copy \"%s\" \"%s\"' -Verb RunAs", f, originalFile)).Run()
      } else {
         data, _ := os.ReadFile(f)
         os.WriteFile(originalFile, data, 0644)
      }
   }
   fmt.Println("[*] Backups still remain in place")
}

func processArguments() {
   args := strings.Join(os.Args[1:], " ")
   if strings.Contains(strings.ToUpper(args), "DEBUG") {
      debug = true
   }
   if strings.Contains(strings.ToUpper(args), "PATCHFSTAB") {
      patchFstab = true
   }
   if strings.Contains(strings.ToUpper(args), "GETUSBHPMODZ") {
      getUSBHPmodZ = true
   }
   if strings.Contains(strings.ToUpper(args), "LISTALLAVDS") {
      listAllAVDs = true
   }
   if strings.Contains(strings.ToUpper(args), "INSTALLAPPS") {
      installApps = true
   }
   if len(os.Args) > 2 {
      switch os.Args[2] {
      case "restore":
         restore = true
      case "InstallKernelModules":
         installKernelModules = true
      case "InstallPrebuiltKernelModules":
         installPrebuiltKernelModules = true
      }
   }

   if len(os.Args) > 1 && !listAllAVDs && !installApps {
      ramdiskImg = true
   }

   if len(os.Args) == 1 {
      noParamsAtAll = true
   }
}

func installapps() {
   fmt.Println("[-] Install all APKs placed in the Apps folder")
   apks, _ := filepath.Glob("APPS/*.apk")
   for _, apk := range apks {
   whileloop:
      fmt.Printf("[*] Trying to install %s\n", apk)
      out, _ := exec.Command("adb", "install", "-r", "-d", apk).CombinedOutput()
      fmt.Printf("[-] %s\n", string(out))
      if strings.Contains(string(out), "INSTALL_FAILED_UPDATE_INCOMPATIBLE") {
         parts := strings.Fields(string(out))
         for i, p := range parts {
            if strings.Contains(p, "Package") && i+1 < len(parts) {
               packageName := parts[i+1]
               fmt.Printf("[*] Need to uninstall %s first\n", packageName)
               uninstallOut, _ := exec.Command("adb", "uninstall", packageName).CombinedOutput()
               fmt.Printf("[-] %s\n", string(uninstallOut))
               goto whileloop
            }
         }
      }
   }
}

func showHelpText() {
   fmt.Println("rootAVD A Script to root AVD by NewBit XDA")
   fmt.Println()
   fmt.Println("Usage:\trootAVD [DIR/ramdisk.img] [OPTIONS] | [EXTRA ARGUMENTS]")
   fmt.Println("or:\trootAVD [ARGUMENTS]")
   // ... (Der Rest des Hilfetextes kann hier hinzugefügt werden)
   findSystemImages()
}

func getANDROIDHOME() {
   envVar := ""
   if home, ok := os.LookupEnv("ANDROID_HOME"); ok {
      androidHome = home
      envVar = "%ANDROID_HOME%"
   } else {
      localAppData := os.Getenv("LOCALAPPDATA")
      androidHome = filepath.Join(localAppData, "Android", "Sdk")
      envVar = "%LOCALAPPDATA%\\Android\\Sdk"
   }
   if _, err := os.Stat(filepath.Join(androidHome, "system-images")); os.IsNotExist(err) {
      fmt.Println("Neither system-images nor ramdisk files could be found")
      os.Exit(1)
   }
   fmt.Printf("- use %s to search for AVD system images\n", envVar)
}

func findSystemImages() {
   fmt.Println()
   var sysImgs []string
   filepath.Walk(filepath.Join(androidHome, "system-images"), func(path string, info os.FileInfo, err error) error {
      if !info.IsDir() && strings.HasPrefix(info.Name(), "ramdisk") && strings.HasSuffix(info.Name(), ".img") {
         p, _ := filepath.Rel(androidHome, path)
         sysImgs = append(sysImgs, p)
      }
      return nil
   })

   fmt.Println("Command Examples:")
   fmt.Println("rootAVD.exe")
   fmt.Println("rootAVD.exe ListAllAVDs")
   fmt.Println("rootAVD.exe InstallApps")
   fmt.Println()

   for _, img := range sysImgs {
      fmt.Printf("rootAVD.exe %s\n", img)
      fmt.Printf("rootAVD.exe %s FAKEBOOTIMG\n", img)
      // ... (Weitere Beispiele können hier hinzugefügt werden)
      fmt.Println()
   }
}

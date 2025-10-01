package main

import (
   "fmt"
   "io"
   "log"
   "os"
   "os/exec"
   "path/filepath"
   "runtime"
   "strings"
)

// Global variables hold the application's state, derived from arguments and environment.
var (
   adbBaseDir                   = ""
   adbWorkdir                   = "/data/data/com.android.shell"
   androidHome                  = ""
   avdPath                      = ""
   avdPathWithRdffile           = ""
   bzFile                       = ""
   debug                        = false
   installApps                  = false
   installKernelModules         = false
   installPrebuiltKernelModules = false
   krFile                       = "kernel-ranchu"
   listAllAVDs                  = false
   magiskZip                    = ""
   ramdiskImg                   = false
   rdfFile                      = ""
   restore                      = false
   rootAVD                      = ""
)

func main() {
   // --- Configure the standard logger ---
   // Remove the default timestamp and source file information.
   log.SetFlags(0)
   // Set a custom prefix to match the script's original style for fatal errors.
   log.SetPrefix("[!] FATAL: ")
   // --- End of logger configuration ---

   var err error

   processArguments()

   if err = getANDROIDHOME(); err != nil {
      log.Fatalf("%v", err)
   }

   // Handle standalone commands
   if listAllAVDs {
      showHelpText()
      return
   }
   if installApps {
      if err = testADB(); err != nil {
         log.Fatalf("Cannot install apps without a working ADB connection: %v", err)
      }
      if err = installapps(); err != nil {
         fmt.Printf("[!] ERROR: %v\n", err)
      }
      return
   }

   // From this point, a path to the ramdisk.img is required.
   if len(os.Args) < 2 {
      showHelpText()
      return
   }

   ramdiskArgPath := filepath.Join(androidHome, os.Args[1])
   if _, err := os.Stat(ramdiskArgPath); os.IsNotExist(err) {
      log.Fatalf("File or directory not found: %s", os.Args[1])
   }

   fmt.Println("[*] Setting Directories")
   avdPathWithRdffile = ramdiskArgPath
   avdPath = filepath.Dir(avdPathWithRdffile)
   rdfFile = filepath.Base(avdPathWithRdffile)

   if fi, err := os.Stat(avdPathWithRdffile); err == nil && fi.IsDir() {
      log.Fatalf("The provided path is a directory, not a file: %s", avdPathWithRdffile)
   }

   testWritePerm(rdfFile)

   if restore {
      if err = restoreBackups(); err != nil {
         fmt.Printf("[!] ERROR during restore: %v\n", err)
      }
      return
   }

   if err = testADB(); err != nil {
      log.Fatalf("%v", err)
   }

   rootAVD, err = os.Getwd()
   if err != nil {
      log.Fatalf("Could not get current working directory: %v", err)
   }
   magiskZip = filepath.Join(rootAVD, "Magisk.zip")
   bzFile = filepath.Join(rootAVD, "bzImage")
   adbBaseDir = adbWorkdir + "/Magisk"

   if err = testADBWORKDIR(); err != nil {
      log.Fatalf("%v", err)
   }

   if err = os.Chdir(rootAVD); err != nil {
      log.Fatalf("Could not change directory to %s: %v", rootAVD, err)
   }

   fmt.Println("[*] Preparing ADB working space on AVD...")
   _ = runADBCommand("shell", "rm", "-rf", adbBaseDir)
   if err = runADBCommand("shell", "mkdir", adbBaseDir); err != nil {
      log.Fatalf("Could not create ADB working directory on AVD: %v", err)
   }

   fmt.Println("[*] Looking for Magisk installer Zip...")
   if _, err = os.Stat(magiskZip); os.IsNotExist(err) {
      fmt.Println("[-] Warning: Magisk.zip not found. Please place it in the script directory to patch.")
   } else {
      if err = pushtoAVD(magiskZip, ""); err != nil {
         log.Fatalf("%v", err)
      }
   }

   initramfs := filepath.Join(rootAVD, "initramfs.img")

   if ramdiskImg {
      if !strings.Contains(strings.ToLower(rdfFile), "ramdisk") || !strings.HasSuffix(strings.ToLower(rdfFile), ".img") {
         log.Fatalf("The provided file does not appear to be a ramdisk image.")
      }
      if err = createBackup(rdfFile); err != nil {
         log.Fatalf("%v", err)
      }
      if err = pushtoAVD(avdPathWithRdffile, "ramdisk.img"); err != nil {
         log.Fatalf("%v", err)
      }
      if installKernelModules {
         if _, err = os.Stat(initramfs); err == nil {
            if err = pushtoAVD(initramfs, ""); err != nil {
               log.Fatalf("%v", err)
            }
         }
      }
   }

   if err = pushtoAVD("rootAVD.sh", ""); err != nil {
      log.Fatalf("Could not push the rootAVD.sh script: %v", err)
   }

   fmt.Println("[-] Running the patch script on the AVD...")
   if err = runADBCommand("shell", "sh", adbBaseDir+"/rootAVD.sh", strings.Join(os.Args[1:], " ")); err != nil {
      log.Fatalf("The patch script failed to execute on the AVD.")
   }
   fmt.Println("[+] Patch script executed successfully on the AVD.")

   if !debug && ramdiskImg {
      localPatchedRamdisk := filepath.Join(rootAVD, "ramdiskpatched4AVD.img")
      if err = pullfromAVD("ramdiskpatched4AVD.img", localPatchedRamdisk); err != nil {
         log.Fatalf("%v", err)
      }
      if err = copyFile(localPatchedRamdisk, avdPathWithRdffile); err != nil {
         log.Fatalf("Failed to copy patched ramdisk: %v\n[!] This is likely a permissions issue. Try running with administrator privileges.", err)
      }
      _ = os.Remove(localPatchedRamdisk)

      _ = pullfromAVD("Magisk.apk", filepath.Join(rootAVD, "Apps"))
      _ = pullfromAVD("Magisk.zip", "")

      if installPrebuiltKernelModules {
         if err = pullfromAVD(bzFile, ""); err == nil {
            if err = installKernelModulesFunc(); err != nil {
               fmt.Printf("[!] ERROR: %v\n", err)
            }
         }
      }
      if installKernelModules {
         if err = installKernelModulesFunc(); err != nil {
            fmt.Printf("[!] ERROR: %v\n", err)
         }
      }

      fmt.Println("[-] Cleaning up ADB working space...")
      _ = runADBCommand("shell", "rm", "-rf", adbBaseDir)

      if err = installapps(); err != nil {
         fmt.Printf("[!] ERROR: %v\n", err)
      }

      fmt.Println("[-] Shut-Down and Reboot [Cold Boot Now] the AVD to see if it worked.")
      shutDownAVD()
   }
}

func getANDROIDHOME() error {
   var sdkPath string
   var envVarSource string

   sdkPath, isSet := os.LookupEnv("ANDROID_HOME")
   if isSet && sdkPath != "" {
      envVarSource = "ANDROID_HOME variable"
   } else {
      homeDir, err := os.UserHomeDir()
      if err != nil {
         return fmt.Errorf("could not determine user home directory: %w", err)
      }
      switch runtime.GOOS {
      case "windows":
         sdkPath = filepath.Join(os.Getenv("LOCALAPPDATA"), "Android", "Sdk")
         envVarSource = "default Windows location"
      case "darwin":
         sdkPath = filepath.Join(homeDir, "Library", "Android", "sdk")
         envVarSource = "default macOS location"
      default:
         sdkPath = filepath.Join(homeDir, "Android", "Sdk")
         envVarSource = "default Linux location"
      }
   }

   fmt.Printf("[-] Probing for Android SDK via %s...\n", envVarSource)
   sysImgPath := filepath.Join(sdkPath, "system-images")
   if _, err := os.Stat(sysImgPath); os.IsNotExist(err) {
      return fmt.Errorf("could not find a valid Android SDK at '%s'. Please set ANDROID_HOME", sdkPath)
   }

   androidHome = sdkPath
   fmt.Printf("[+] Android SDK found at: %s\n", androidHome)
   return nil
}

func testADB() error {
   fmt.Println("[-] Testing ADB connection...")
   if err := runADBCommand("shell", "-n", "echo", "true"); err != nil {
      return fmt.Errorf("ADB connection failed. Please ensure an AVD is running and accessible")
   }
   fmt.Println("[+] ADB connection is working.")
   return nil
}

func testADBWORKDIR() error {
   fmt.Println("[*] Testing the ADB working space")
   if err := runADBCommand("shell", "cd", adbWorkdir); err != nil {
      return fmt.Errorf("ADB working directory %s is not available: %w", adbWorkdir, err)
   }
   fmt.Printf("[+] ADB working directory %s is available.\n", adbWorkdir)
   return nil
}

func pushtoAVD(src, dst string) error {
   srcBase := filepath.Base(src)
   var args []string
   var pushDestination string
   if dst == "" {
      args = []string{"push", src, adbBaseDir}
      pushDestination = adbBaseDir
   } else {
      dstBase := filepath.Base(dst)
      args = []string{"push", src, adbBaseDir + "/" + dstBase}
      pushDestination = adbBaseDir + "/" + dstBase
   }
   fmt.Printf("[*] Pushing %s to %s\n", srcBase, pushDestination)
   if err := runADBCommand(args...); err != nil {
      return fmt.Errorf("failed to push %s to AVD: %w", srcBase, err)
   }
   return nil
}

func pullfromAVD(src, dst string) error {
   srcBase := filepath.Base(src)
   adbSrcPath := adbBaseDir + "/" + srcBase
   var args []string
   pullDestination := "(current directory)"
   if dst != "" {
      args = []string{"pull", adbSrcPath, dst}
      pullDestination = dst
   } else {
      args = []string{"pull", adbSrcPath}
   }
   fmt.Printf("[*] Pulling %s to %s\n", srcBase, pullDestination)
   if err := runADBCommand(args...); err != nil {
      return fmt.Errorf("failed to pull %s from AVD: %w", srcBase, err)
   }
   return nil
}

func createBackup(file string) error {
   backupFile := file + ".backup"
   sourcePath := filepath.Join(avdPath, file)
   backupPath := filepath.Join(avdPath, backupFile)

   if _, err := os.Stat(backupPath); os.IsNotExist(err) {
      fmt.Printf("[*] Creating backup of %s...\n", file)
      if err := copyFile(sourcePath, backupPath); err != nil {
         return fmt.Errorf("could not create backup for %s: %w", file, err)
      }
      fmt.Printf("[+] Backup created: %s\n", backupPath)
   } else {
      fmt.Println("[-] Backup file already exists, skipping.")
   }
   return nil
}

func installKernelModulesFunc() error {
   if _, err := os.Stat(bzFile); err == nil {
      if err := createBackup(krFile); err != nil {
         return err // Propagate error from backup creation
      }
      fmt.Printf("[*] Copying %s (Kernel) into %s\n", bzFile, krFile)
      destination := filepath.Join(avdPath, krFile)
      if err := copyFile(bzFile, destination); err != nil {
         return fmt.Errorf("failed to copy kernel file: %w", err)
      }
      _ = os.Remove(bzFile)
      _ = os.Remove(filepath.Join(rootAVD, "initramfs.img"))
   } else {
      fmt.Printf("[-] Kernel file %s not found, skipping installation.\n", bzFile)
   }
   return nil
}

func restoreBackups() error {
   backupPattern := filepath.Join(avdPath, "*.backup")
   files, err := filepath.Glob(backupPattern)
   if err != nil {
      return fmt.Errorf("could not search for backup files: %w", err)
   }
   if len(files) == 0 {
      fmt.Println("[-] No backup files found to restore.")
      return nil
   }
   for _, f := range files {
      originalFile := strings.TrimSuffix(f, ".backup")
      fmt.Printf("[*] Restoring %s -> %s\n", filepath.Base(f), filepath.Base(originalFile))
      if err := copyFile(f, originalFile); err != nil {
         fmt.Printf("[!] Warning: Failed to restore %s: %v\n", filepath.Base(originalFile), err)
         // Continue to next file
      }
   }
   fmt.Println("[+] Backup restoration process finished.")
   return nil
}

func installapps() error {
   fmt.Println("[-] Installing all APKs from the 'Apps' folder...")
   appsDir := "Apps"
   if _, err := os.Stat(appsDir); os.IsNotExist(err) {
      fmt.Println("[-] 'Apps' directory not found, skipping APK installation.")
      return nil // Not a fatal error
   }
   apks, err := filepath.Glob(filepath.Join(appsDir, "*.apk"))
   if err != nil {
      return fmt.Errorf("could not search for APKs: %w", err)
   }
   if len(apks) == 0 {
      fmt.Println("[-] No APKs found in the 'Apps' directory.")
      return nil
   }
   for _, apk := range apks {
   installAttempt:
      fmt.Printf("[*] Trying to install %s\n", filepath.Base(apk))
      out, err := exec.Command("adb", "install", "-r", "-d", apk).CombinedOutput()
      outputStr := string(out)
      fmt.Printf("[-] %s\n", outputStr)
      if err != nil {
         if strings.Contains(outputStr, "INSTALL_FAILED_UPDATE_INCOMPATIBLE") {
            parts := strings.Fields(outputStr)
            for i, p := range parts {
               if strings.Contains(p, "Package") && i+1 < len(parts) {
                  packageName := parts[i+1]
                  fmt.Printf("[*] Incompatible update. Uninstalling %s first...\n", packageName)
                  if err := runADBCommand("uninstall", packageName); err == nil {
                     fmt.Println("[+] Uninstalled successfully, retrying installation...")
                     goto installAttempt
                  } else {
                     fmt.Printf("[!] Failed to uninstall %s.\n", packageName)
                  }
                  break
               }
            }
         } else {
            fmt.Printf("[!] An error occurred during installation of %s.\n", filepath.Base(apk))
         }
      }
   }
   return nil
}

// --- Core Utility Functions ---

func runADBCommand(args ...string) error {
   cmd := exec.Command("adb", args...)
   output, err := cmd.CombinedOutput()
   if err != nil {
      return fmt.Errorf("adb command failed: adb %s\n---Output---\n%s", strings.Join(args, " "), string(output))
   }
   return nil
}

func copyFile(src, dst string) error {
   sourceFile, err := os.Open(src)
   if err != nil {
      return fmt.Errorf("could not open source file %s: %w", src, err)
   }
   defer sourceFile.Close()

   destinationFile, err := os.Create(dst)
   if err != nil {
      return fmt.Errorf("could not create destination file %s: %w", dst, err)
   }
   defer destinationFile.Close()

   _, err = io.Copy(destinationFile, sourceFile)
   if err != nil {
      return fmt.Errorf("failed to copy data from %s to %s: %w", src, err)
   }
   return nil
}

// --- Non-Error-Returning Functions ---

func processArguments() {
   args := strings.Join(os.Args[1:], " ")
   if strings.Contains(strings.ToUpper(args), "DEBUG") {
      debug = true
   }
   if strings.Contains(strings.ToUpper(args), "LISTALLAVDS") {
      listAllAVDs = true
   }
   if strings.Contains(strings.ToUpper(args), "INSTALLAPPS") {
      installApps = true
   }
   if len(os.Args) > 1 {
      for _, arg := range os.Args[1:] {
         switch strings.ToLower(arg) {
         case "restore":
            restore = true
         case "installkernelmodules":
            installKernelModules = true
         case "installprebuiltkernelmodules":
            installPrebuiltKernelModules = true
         }
      }
   }
   if len(os.Args) > 1 && !listAllAVDs && !installApps {
      ramdiskImg = true
   }
}

func shutDownAVD() {
   fmt.Println("[-] Attempting to shut down the AVD...")
   if err := runADBCommand("shell", "setprop", "sys.powerctl", "shutdown"); err != nil {
      fmt.Println("[!] Warning: Failed to send shutdown command. Please shut down the AVD manually.")
   }
   fmt.Println("[+] If the AVD doesn't shut down, please do it manually from Android Studio.")
}

func testWritePerm(file string) {
   fmt.Println("[*] Testing for write permissions in AVD directory...")
   tempFile := filepath.Join(avdPath, file+".temp")
   if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
      fmt.Println("[!] Elevated write permissions appear to be needed to access the Android SDK system images.")
   } else {
      fmt.Println("[+] Write permissions are sufficient.")
      _ = os.Remove(tempFile)
   }
}

// showHelpText displays the main help message for the application.
// It also calls findSystemImages to provide dynamic command examples.
func showHelpText() {
   fmt.Println("rootAVD: A script to root Android Virtual Devices (AVD)")
   fmt.Println("Originally by NewBit XDA, converted to Go.")
   fmt.Println()
   fmt.Println("Usage:   rootAVD [path/to/ramdisk.img] [OPTIONS...]")
   fmt.Println("   or:   rootAVD [COMMAND]")
   fmt.Println()
   fmt.Println("Commands:")
   fmt.Println("  ListAllAVDs        Lists command examples for all found AVD system images.")
   fmt.Println("  InstallApps        Installs all APKs located in the 'Apps' folder.")
   fmt.Println()
   fmt.Println("Options (used with a ramdisk path):")
   fmt.Println("  restore            Restores backups for the given ramdisk's directory.")
   fmt.Println("  InstallKernelModules")
   fmt.Println("                     Installs a custom kernel from local bzImage and initramfs.img files.")
   fmt.Println()
   fmt.Println("Extra Arguments (can be combined):")
   fmt.Println("  DEBUG              Enables debugging mode; prevents writing files back to the AVD.")
   fmt.Println("  PATCHFSTAB         (Functionality to be implemented in the accompanying rootAVD.sh)")

   // Call the function to show examples and gracefully handle any errors.
   if err := findSystemImages(); err != nil {
      // If finding images fails, we print a warning but don't terminate.
      // The user still gets the main help text.
      fmt.Printf("\n[!] Warning: Could not generate dynamic command examples: %v\n", err)
   }
}

// findSystemImages walks the Android SDK's system-images directory to find all
// available ramdisk.img files and prints helpful command examples for each.
// It returns an error if the directory cannot be accessed.
func findSystemImages() error {
   fmt.Println("\n--- Command Examples ---")
   if androidHome == "" {
      // This is a pre-condition failure, return an error.
      return fmt.Errorf("ANDROID_HOME not set; cannot search for system images")
   }

   sysImgDir := filepath.Join(androidHome, "system-images")
   var foundImages []string

   // Walk the system-images directory. filepath.Walk is the correct tool for this.
   walkErr := filepath.Walk(sysImgDir, func(path string, info os.FileInfo, err error) error {
      // If the walker function is passed an error (e.g., a permissions issue),
      // we must handle it. Returning the error will stop the walk.
      if err != nil {
         return err
      }
      // We are only interested in files named "ramdisk.img".
      if !info.IsDir() && info.Name() == "ramdisk.img" {
         // Get the path relative to the androidHome directory for a clean example.
         relativePath, err := filepath.Rel(androidHome, path)
         if err != nil {
            // This is unlikely but possible; continue without this entry.
            fmt.Printf("[!] Warning: Could not calculate relative path for %s\n", path)
            return nil
         }
         foundImages = append(foundImages, relativePath)
      }
      return nil // Continue the walk
   })

   // After the walk is complete, check if it encountered an error.
   if walkErr != nil {
      return fmt.Errorf("error searching for system images: %w", walkErr)
   }

   if len(foundImages) == 0 {
      fmt.Println("[-] No AVD system images with a ramdisk.img were found.")
      return nil
   }

   // Print examples for each found image.
   for _, img := range foundImages {
      // Use filepath.ToSlash for consistent forward slashes in output,
      // which works nicely in Windows cmd/PowerShell as well as Unix shells.
      fmt.Printf("rootAVD.exe %s\n", filepath.ToSlash(img))
      fmt.Printf("rootAVD.exe %s restore\n", filepath.ToSlash(img))
      fmt.Println("--------------------------------------------------")
   }

   return nil // Success
}

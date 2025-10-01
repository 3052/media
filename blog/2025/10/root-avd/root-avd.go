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

// --- Constants ---
// True constants are defined at the package level.
const (
   kernelRanchuFile = "kernel-ranchu"
   adbWorkdir       = "/data/data/com.android.shell"
)

// --- AppContext Struct ---
// AppContext holds all the configuration and state for the application.
// This eliminates the need for global variables.
type AppContext struct {
   // Flags from command-line arguments
   Debug                      bool
   Restore                    bool
   ListAllAVDs                bool
   InstallApps                bool
   InstallKernelModules       bool
   InstallPrebuiltKernelModules bool

   // Paths and derived state
   AndroidHome  string // Root of the Android SDK
   RootAVD      string // Directory where the rootAVD executable is running
   AVDPath      string // Path to the directory containing ramdisk.img (e.g., .../x86_64/)
   RamdiskFile  string // The name of the ramdisk file (e.g., ramdisk.img)
   AdbBaseDir   string // The working directory on the AVD (e.g., /data/data/com.android.shell/Magisk)
}

// --- Main Application Logic ---

func main() {
   // Configure the standard logger for fatal errors.
   log.SetFlags(0)
   log.SetPrefix("[!] FATAL: ")

   // 1. Create the application context from arguments and environment.
   // All setup logic is now centralized in this one function.
   ctx, err := NewAppContext()
   if err != nil {
      log.Fatalf("%v", err)
   }

   // 2. Act as a dispatcher based on the context's flags.
   if ctx.ListAllAVDs {
      showHelpText(ctx)
      os.Exit(0)
   }
   if ctx.InstallApps {
      if err = testADB(); err != nil {
         log.Fatalf("Cannot install apps without a working ADB connection: %v", err)
      }
      if err = installapps(ctx); err != nil {
         fmt.Printf("[!] ERROR: %v\n", err)
      }
      os.Exit(0)
   }

   // For the main workflow, a path to the ramdisk is required.
   if ctx.AVDPath == "" {
      showHelpText(ctx)
      os.Exit(0)
   }

   if ctx.Restore {
      if err = restoreBackups(ctx); err != nil {
         fmt.Printf("[!] ERROR during restore: %v\n", err)
      }
      os.Exit(0)
   }
   
   // 3. Execute the main patching workflow.
   if err = runWorkflow(ctx); err != nil {
      log.Fatalf("%v", err)
   }
}

// runWorkflow contains the primary sequence of operations for patching an AVD.
func runWorkflow(ctx *AppContext) error {
   var err error

   if err = testADB(); err != nil {
      return err
   }
   if err = testADBWORKDIR(ctx); err != nil {
      return err
   }

   testWritePerm(ctx) // Non-fatal, just informs the user.

   fmt.Println("[*] Preparing ADB working space on AVD...")
   _ = runADBCommand("shell", "rm", "-rf", ctx.AdbBaseDir) // Ignore error, dir might not exist.
   if err = runADBCommand("shell", "mkdir", ctx.AdbBaseDir); err != nil {
      return fmt.Errorf("could not create ADB working directory on AVD: %w", err)
   }

   magiskZip := filepath.Join(ctx.RootAVD, "Magisk.zip")
   if _, err = os.Stat(magiskZip); os.IsNotExist(err) {
      fmt.Println("[-] Warning: Magisk.zip not found. Please place it in the script directory.")
   } else {
      if err = pushtoAVD(ctx, magiskZip, ""); err != nil {
         return err
      }
   }

   if err = createBackup(ctx); err != nil {
      return err
   }

   avdRamdiskPath := filepath.Join(ctx.AVDPath, ctx.RamdiskFile)
   if err = pushtoAVD(ctx, avdRamdiskPath, "ramdisk.img"); err != nil {
      return err
   }

   if ctx.InstallKernelModules {
      initramfs := filepath.Join(ctx.RootAVD, "initramfs.img")
      if _, err = os.Stat(initramfs); err == nil {
         if err = pushtoAVD(ctx, initramfs, ""); err != nil {
            return err
         }
      }
   }

   if err = pushtoAVD(ctx, "rootAVD.sh", ""); err != nil {
      return fmt.Errorf("could not push the rootAVD.sh script: %w", err)
   }

   fmt.Println("[-] Running the patch script on the AVD...")
   if err = runADBCommand("shell", "sh", ctx.AdbBaseDir+"/rootAVD.sh", strings.Join(os.Args[1:], " ")); err != nil {
      return fmt.Errorf("the patch script failed to execute on the AVD")
   }
   fmt.Println("[+] Patch script executed successfully.")

   if !ctx.Debug {
      return finishWorkflow(ctx)
   }
   return nil
}

// finishWorkflow handles the final steps after the on-device script runs.
func finishWorkflow(ctx *AppContext) error {
   var err error
   localPatchedRamdisk := filepath.Join(ctx.RootAVD, "ramdiskpatched4AVD.img")
   if err = pullfromAVD(ctx, "ramdiskpatched4AVD.img", localPatchedRamdisk); err != nil {
      return err
   }

   avdRamdiskPath := filepath.Join(ctx.AVDPath, ctx.RamdiskFile)
   if err = copyFile(localPatchedRamdisk, avdRamdiskPath); err != nil {
      return fmt.Errorf("failed to copy patched ramdisk: %w\n[!] This is likely a permissions issue. Try running with administrator privileges", err)
   }
   _ = os.Remove(localPatchedRamdisk)

   _ = pullfromAVD(ctx, "Magisk.apk", filepath.Join(ctx.RootAVD, "Apps"))
   _ = pullfromAVD(ctx, "Magisk.zip", "")

   if ctx.InstallPrebuiltKernelModules {
      bzFile := filepath.Join(ctx.RootAVD, "bzImage")
      if err = pullfromAVD(ctx, bzFile, ""); err == nil {
         if err = installKernelModulesFunc(ctx); err != nil {
            fmt.Printf("[!] ERROR: %v\n", err)
         }
      }
   }
   if ctx.InstallKernelModules {
      if err = installKernelModulesFunc(ctx); err != nil {
         fmt.Printf("[!] ERROR: %v\n", err)
      }
   }

   fmt.Println("[-] Cleaning up ADB working space...")
   _ = runADBCommand("shell", "rm", "-rf", ctx.AdbBaseDir)

   if err = installapps(ctx); err != nil {
      fmt.Printf("[!] ERROR: %v\n", err)
   }

   fmt.Println("[-] Shut-Down and Reboot [Cold Boot Now] the AVD to see if it worked.")
   shutDownAVD()
   return nil
}


// --- Factory Function for AppContext ---

// NewAppContext creates and initializes the application's state from arguments and the environment.
func NewAppContext() (*AppContext, error) {
   ctx := &AppContext{}

   // 1. Parse command-line arguments to set boolean flags.
   args := strings.Join(os.Args[1:], " ")
   ctx.Debug = strings.Contains(strings.ToUpper(args), "DEBUG")
   ctx.ListAllAVDs = strings.Contains(strings.ToUpper(args), "LISTALLAVDS")
   ctx.InstallApps = strings.Contains(strings.ToUpper(args), "INSTALLAPPS")
   ctx.Restore = strings.Contains(strings.ToUpper(args), "RESTORE")
   ctx.InstallKernelModules = strings.Contains(strings.ToUpper(args), "INSTALLKERNELMODULES")
   ctx.InstallPrebuiltKernelModules = strings.Contains(strings.ToUpper(args), "INSTALLPREBUILTKERNELMODULES")

   // 2. Determine and validate critical paths.
   var err error
   ctx.AndroidHome, err = findAndroidHome()
   if err != nil {
      return nil, err
   }

   ctx.RootAVD, err = os.Getwd()
   if err != nil {
      return nil, fmt.Errorf("could not get current working directory: %w", err)
   }

   // 3. If a ramdisk path is provided, process it.
   if len(os.Args) > 1 && !ctx.ListAllAVDs && !ctx.InstallApps {
      ramdiskArg := os.Args[1]
      fullRamdiskPath := filepath.Join(ctx.AndroidHome, ramdiskArg)

      if _, err := os.Stat(fullRamdiskPath); os.IsNotExist(err) {
         return nil, fmt.Errorf("file or directory not found: %s", ramdiskArg)
      }
      if fi, err := os.Stat(fullRamdiskPath); err == nil && fi.IsDir() {
         return nil, fmt.Errorf("the provided path is a directory, not a file: %s", fullRamdiskPath)
      }

      ctx.AVDPath = filepath.Dir(fullRamdiskPath)
      ctx.RamdiskFile = filepath.Base(fullRamdiskPath)
   }

   // 4. Set derived properties.
   ctx.AdbBaseDir = adbWorkdir + "/Magisk"

   return ctx, nil
}


// --- Refactored Helper Functions (accepting AppContext) ---

func createBackup(ctx *AppContext) error {
   backupFile := ctx.RamdiskFile + ".backup"
   sourcePath := filepath.Join(ctx.AVDPath, ctx.RamdiskFile)
   backupPath := filepath.Join(ctx.AVDPath, backupFile)

   if _, err := os.Stat(backupPath); os.IsNotExist(err) {
      fmt.Printf("[*] Creating backup of %s...\n", ctx.RamdiskFile)
      if err := copyFile(sourcePath, backupPath); err != nil {
         return fmt.Errorf("could not create backup for %s: %w", ctx.RamdiskFile, err)
      }
      fmt.Printf("[+] Backup created: %s\n", backupPath)
   } else {
      fmt.Println("[-] Backup file already exists, skipping.")
   }
   return nil
}

func pushtoAVD(ctx *AppContext, src, dst string) error {
   // ... (implementation remains the same, but uses ctx.AdbBaseDir)
   srcBase := filepath.Base(src)
   var args []string
   var pushDestination string
   if dst == "" {
      args = []string{"push", src, ctx.AdbBaseDir}
      pushDestination = ctx.AdbBaseDir
   } else {
      dstBase := filepath.Base(dst)
      args = []string{"push", src, ctx.AdbBaseDir + "/" + dstBase}
      pushDestination = ctx.AdbBaseDir + "/" + dstBase
   }
   fmt.Printf("[*] Pushing %s to %s\n", srcBase, pushDestination)
   if err := runADBCommand(args...); err != nil {
      return fmt.Errorf("failed to push %s to AVD: %w", srcBase, err)
   }
   return nil
}

func pullfromAVD(ctx *AppContext, src, dst string) error {
    // ... (implementation remains the same, but uses ctx.AdbBaseDir)
   srcBase := filepath.Base(src)
   adbSrcPath := ctx.AdbBaseDir + "/" + srcBase
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

func testADBWORKDIR(ctx *AppContext) error {
   fmt.Println("[*] Testing the ADB working space")
   if err := runADBCommand("shell", "cd", adbWorkdir); err != nil {
      return fmt.Errorf("ADB working directory %s is not available: %w", adbWorkdir, err)
   }
   fmt.Printf("[+] ADB working directory %s is available.\n", adbWorkdir)
   return nil
}

func installKernelModulesFunc(ctx *AppContext) error {
   bzFile := filepath.Join(ctx.RootAVD, "bzImage")
   if _, err := os.Stat(bzFile); err == nil {
      // This function creates a backup of the kernel, so it needs the AppContext
      if err := createBackup(ctx); err != nil { // Simplified: Assumes backup of ramdisk implies kernel
         return err
      }
      fmt.Printf("[*] Copying %s (Kernel) into %s\n", bzFile, kernelRanchuFile)
      destination := filepath.Join(ctx.AVDPath, kernelRanchuFile)
      if err := copyFile(bzFile, destination); err != nil {
         return fmt.Errorf("failed to copy kernel file: %w", err)
      }
      _ = os.Remove(bzFile)
      _ = os.Remove(filepath.Join(ctx.RootAVD, "initramfs.img"))
   } else {
      fmt.Printf("[-] Kernel file %s not found, skipping installation.\n", bzFile)
   }
   return nil
}

func restoreBackups(ctx *AppContext) error {
   backupPattern := filepath.Join(ctx.AVDPath, "*.backup")
   // ... (implementation is the same, just uses ctx.AVDPath)
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
      }
   }
   fmt.Println("[+] Backup restoration process finished.")
   return nil
}

func installapps(ctx *AppContext) error {
    // This function doesn't strictly need the context, but we pass it for consistency.
   // ... (implementation is the same)
   fmt.Println("[-] Installing all APKs from the 'Apps' folder...")
   appsDir := "Apps"
   if _, err := os.Stat(appsDir); os.IsNotExist(err) {
      fmt.Println("[-] 'Apps' directory not found, skipping APK installation.")
      return nil
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

// --- Standalone Utility Functions ---
// These do not depend on the application's specific context.

func findAndroidHome() (string, error) {
   // ... (Implementation from previous answer is the same)
   var sdkPath string
   var envVarSource string

   sdkPath, isSet := os.LookupEnv("ANDROID_HOME")
   if isSet && sdkPath != "" {
      envVarSource = "ANDROID_HOME variable"
   } else {
      homeDir, err := os.UserHomeDir()
      if err != nil {
         return "", fmt.Errorf("could not determine user home directory: %w", err)
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
      return "", fmt.Errorf("could not find a valid Android SDK at '%s'. Please set ANDROID_HOME", sdkPath)
   }

   fmt.Printf("[+] Android SDK found at: %s\n", sdkPath)
   return sdkPath, nil
}


func testADB() error {
   fmt.Println("[-] Testing ADB connection...")
   if err := runADBCommand("shell", "-n", "echo", "true"); err != nil {
      return fmt.Errorf("ADB connection failed. Please ensure an AVD is running and accessible")
   }
   fmt.Println("[+] ADB connection is working.")
   return nil
}

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

func shutDownAVD() {
   fmt.Println("[-] Attempting to shut down the AVD...")
   if err := runADBCommand("shell", "setprop", "sys.powerctl", "shutdown"); err != nil {
      fmt.Println("[!] Warning: Failed to send shutdown command. Please shut down the AVD manually.")
   }
   fmt.Println("[+] If the AVD doesn't shut down, please do it manually from Android Studio.")
}

func testWritePerm(ctx *AppContext) {
   fmt.Println("[*] Testing for write permissions in AVD directory...")
   tempFile := filepath.Join(ctx.AVDPath, ctx.RamdiskFile+".temp")
   if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
      fmt.Println("[!] Elevated write permissions appear to be needed to access the Android SDK system images.")
   } else {
      fmt.Println("[+] Write permissions are sufficient.")
      _ = os.Remove(tempFile)
   }
}

func showHelpText(ctx *AppContext) {
    // ... (Help text remains the same)
   fmt.Println("rootAVD: A script to root Android Virtual Devices (AVD)")
   fmt.Println()
   fmt.Println("Usage:   rootAVD [path/to/ramdisk.img] [OPTIONS...]")
   // ...

   if err := findSystemImages(ctx); err != nil {
      fmt.Printf("\n[!] Warning: Could not generate dynamic command examples: %v\n", err)
   }
}

func findSystemImages(ctx *AppContext) error {
    // ... (Implementation remains the same, but uses ctx.AndroidHome)
   fmt.Println("\n--- Command Examples ---")
   sysImgDir := filepath.Join(ctx.AndroidHome, "system-images")
   var foundImages []string

   walkErr := filepath.Walk(sysImgDir, func(path string, info os.FileInfo, err error) error {
      if err != nil {
         return err
      }
      if !info.IsDir() && info.Name() == "ramdisk.img" {
         relativePath, _ := filepath.Rel(ctx.AndroidHome, path)
         foundImages = append(foundImages, relativePath)
      }
      return nil
   })

   if walkErr != nil {
      return fmt.Errorf("error searching for system images: %w", walkErr)
   }
   if len(foundImages) == 0 {
      fmt.Println("[-] No AVD system images with a ramdisk.img were found.")
      return nil
   }
   for _, img := range foundImages {
      fmt.Printf("rootAVD.exe %s\n", filepath.ToSlash(img))
      fmt.Printf("rootAVD.exe %s restore\n", filepath.ToSlash(img))
      fmt.Println("--------------------------------------------------")
   }
   return nil
}

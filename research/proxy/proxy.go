package main

import (
   "encoding/json"
   "flag"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "os"
)

const cacheFileName = "cache.json"

func main() {
   // Define flags
   // flag -x: accept proxy URL
   proxyFlag := flag.String("x", "", "Set the proxy URL (e.g., http://localhost:8080)")
   // flag -a: action flag
   actionFlag := flag.Bool("a", false, "Execute request to example.com")

   flag.Parse()

   // ## no flags: print usage
   if flag.NFlag() == 0 {
      flag.Usage()
      return
   }

   // ## flag -x
   // Note: flag.NFlag() > 0, checks if x was explicitly set and has a value
   if *proxyFlag != "" {
      // 1. accept proxy URL (handled by *proxyFlag)
      // 2. declare map
      config := make(map[string]string)

      // 3. if `cache.json` exist read it into map
      if _, err := os.Stat(cacheFileName); err == nil {
         fileData, err := os.ReadFile(cacheFileName)
         if err != nil {
            log.Fatalf("Failed to read cache file: %v", err)
         }
         if err := json.Unmarshal(fileData, &config); err != nil {
            log.Fatalf("Failed to parse cache JSON: %v", err)
         }
      }

      // 4. update map with proxy URL
      config["proxy"] = *proxyFlag
      fmt.Printf("Updating proxy configuration to: %s\n", *proxyFlag)

      // 5. write to `cache.json`
      jsonData, err := json.MarshalIndent(config, "", "  ")
      if err != nil {
         log.Fatalf("Failed to marshal JSON: %v", err)
      }
      if err := os.WriteFile(cacheFileName, jsonData, 0644); err != nil {
         log.Fatalf("Failed to write to cache file: %v", err)
      }
      fmt.Println("Configuration saved.")
      return
   }

   // ## flag -a
   if *actionFlag {
      // 1. declare map
      config := make(map[string]string)

      // 2. if `cache.json` exist read it into map
      if _, err := os.Stat(cacheFileName); err == nil {
         fileData, err := os.ReadFile(cacheFileName)
         if err == nil {
            // We ignore unmarshal errors here to ensure map stays empty/valid if file is corrupt
            _ = json.Unmarshal(fileData, &config)
         }
      }

      // Prepare HTTP client
      client := &http.Client{}
      targetURL := "http://example.com"

      // 3. if map has proxy then request with proxy
      if proxyURL, ok := config["proxy"]; ok && proxyURL != "" {
         parsedProxy, err := url.Parse(proxyURL)
         if err != nil {
            log.Fatalf("Invalid proxy URL in cache: %v", err)
         }

         fmt.Printf("Using Proxy: %s\n", proxyURL)
         client.Transport = &http.Transport{
            Proxy: http.ProxyURL(parsedProxy),
         }
      } else {
         // 4. else request without proxy
         fmt.Println("No proxy found, connecting directly...")
      }

      // Perform the request
      resp, err := client.Get(targetURL)
      if err != nil {
         log.Fatalf("Request failed: %v", err)
      }
      defer resp.Body.Close()

      fmt.Printf("Response Status: %s\n", resp.Status)
      
      // Optional: Print a bit of the body to prove it worked
      body, _ := io.ReadAll(resp.Body)
      fmt.Printf("Body length: %d bytes\n", len(body))
   }
}

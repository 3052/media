package main

import (
   "encoding/json"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
)

// handleSetProxy reads config, updates it (if doWrite is true), and sets http.DefaultClient
func handleSetProxy(proxyURL string, doWrite bool) error {
   // 1. Declare map
   config := make(map[string]string)
   // 2. Read cache.json into map
   if fileData, err := os.ReadFile(cacheFileName); err == nil {
      if err := json.Unmarshal(fileData, &config); err != nil {
         return fmt.Errorf("failed to parse cache JSON: %w", err)
      }
   }
   // 3. Update map and write to file ONLY if the write flag is true
   if doWrite {
      // Just assign the value. If it's empty, the key becomes empty.
      config["proxy"] = proxyURL
      // Save to file
      jsonData, err := json.MarshalIndent(config, "", "  ")
      if err != nil {
         return fmt.Errorf("failed to marshal config: %w", err)
      }
      if err := os.WriteFile(cacheFileName, jsonData, 0644); err != nil {
         return fmt.Errorf("failed to write cache file: %w", err)
      }
   }
   // 4. Configure http.DefaultClient based on the map state
   if val := config["proxy"]; val != "" {
      parsedURL, err := url.Parse(val)
      if err != nil {
         return fmt.Errorf("invalid proxy URL stored in cache: %w", err)
      }
      http.DefaultClient.Transport = &http.Transport{
         Proxy: http.ProxyURL(parsedURL),
      }
   }
   return nil
}

// handleAction executes the request using http.DefaultClient
func handleAction() error {
   targetURL := "http://example.com"

   // Perform the HEAD request
   resp, err := http.Head(targetURL)
   if err != nil {
      return fmt.Errorf("http request failed: %w", err)
   }
   defer resp.Body.Close()

   fmt.Printf("Response Status: %s\n", resp.Status)

   return nil
}

const cacheFileName = "cache.json"

func main() {
   // Define flags
   proxyFlag := flag.String("x", "", "Set the proxy URL")
   writeFlag := flag.Bool("w", false, "Write/Update the configuration")
   actionFlag := flag.Bool("a", false, "Execute request to example.com")

   flag.Parse()

   // 0. Check if any flags were provided
   if flag.NFlag() == 0 {
      flag.Usage()
      return
   }

   // 1. Always configure the environment
   if err := handleSetProxy(*proxyFlag, *writeFlag); err != nil {
      log.Fatalf("Error configuring proxy: %v", err)
   }

   // 2. If -a is set, perform the action
   if *actionFlag {
      if err := handleAction(); err != nil {
         log.Fatalf("Error executing action: %v", err)
      }
   }
}

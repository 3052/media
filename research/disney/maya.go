package main

import (
   "flag"
   "fmt"
)

func Usage(groups [][]*flag.Flag) error {
   seen := map[string]bool{}
   // 1. Print usage and mark flags as seen
   for i, group := range groups {
      if i >= 1 {
         fmt.Println()
      }
      for _, f := range group {
         fmt.Printf("-%v %v\n", f.Name, f.Usage)
         if f.DefValue != "" {
            fmt.Printf("\tdefault %v\n", f.DefValue)
         }
         seen[f.Name] = true
      }
   }
   // 2. Check for missing flags
   var missing string
   flag.VisitAll(func(f *flag.Flag) {
      if !seen[f.Name] {
         missing = f.Name
      }
   })
   if missing != "" {
      return fmt.Errorf("defined flag missing: -%s", missing)
   }
   return nil
}

func BoolVar(value *bool, name, usage string) *flag.Flag {
   flag.BoolVar(value, name, *value, usage)
   return flag.Lookup(name)
}

func IntVar(value *int, name, usage string) *flag.Flag {
   flag.IntVar(value, name, *value, usage)
   return flag.Lookup(name)
}

func StringVar(value *string, name, usage string) *flag.Flag {
   flag.StringVar(value, name, *value, usage)
   return flag.Lookup(name)
}

func IsSet(f *flag.Flag) bool {
   var set bool
   flag.Visit(func(g *flag.Flag) {
      if g.Name == f.Name {
         set = true
      }
   })
   return set
}

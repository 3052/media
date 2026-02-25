Go language, I need a script. user will provide zero or one flag

## no flags

print usage

## flag -x

1. accept proxy URL
2. declare map
3. if `cache.json` exist read it into map
4. update map with proxy URL
5. write to `cache.json`

## flag -a  

1. declare map
2. if `cache.json` exist read it into map
3. if map has proxy then request http://example.com with proxy
4. else request without proxy

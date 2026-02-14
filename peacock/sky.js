function jhexdump(array) {
    if(!array) return;
    console.log("---------jhexdump start---------");
    var ptr = Memory.alloc(array.length);
    for(var i = 0; i < array.length; ++i)
        Memory.writeS8(ptr.add(i), array[i]);
    console.log(hexdump(ptr, {offset: 0, length: array.length, header: false, ansi: false}));
    console.log("---------jhexdump end---------");
}

function java_hook(){
    Java.perform(function(){
        let HMACCls = Java.use("com.sky.sps.security.HMAC");
        let SecurityUtilsCls = Java.use("com.sky.sps.utils.SecurityUtils");
        let SecretKeySpecCls = Java.use("javax.crypto.spec.SecretKeySpec");
        let MacCls = Java.use("javax.crypto.Mac");
        HMACCls.calculate.overload('java.lang.String', 'boolean').implementation = function(text, flag){
            console.log("---------enter calculate---------");
            let ret = this.calculate(text, flag);
            console.log(text, flag, ret);
            jhexdump(ret);
            return ret;
        }
        SecurityUtilsCls.createMD5Digest.overload('java.lang.String').implementation = function(text){
            console.log("---------enter createMD5Digest---------");
            let ret = this.createMD5Digest(text);
            console.log(text, ret);
            return ret;
        }
        SecretKeySpecCls.$init.overload('[B', 'java.lang.String').implementation = function(key, method){
            console.log("---------enter SecretKeySpec init---------");
            jhexdump(key);
            let ret = this.$init(key, method);
            console.log(key, method, ret);
            return ret;
        }
        MacCls.doFinal.overload('[B').implementation = function(data){
            console.log("---------enter SecretKeySpec init---------");
            jhexdump(data);
            let ret = this.doFinal(data);
            console.log(data, ret);
            return ret;
        }
    })
}

setImmediate(java_hook)

// frida 14.2.18
// frida -U -n com.peacocktv.peacockandroid -l peacock.js -o peacock.log

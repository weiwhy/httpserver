package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)
var (
	flagHelp    = flag.Bool("h", false, "Shows usage options.")
	flagDir 	= flag.String("d", ".", "指定文件目录，默认当前目录")
	flagLis     = flag.String("l", "0.0.0.0:80", "指定监听地址和端口")
	flagSSL		= flag.Bool("s",false,"指定使用https")
	flagCert	= flag.String("c","cert.pem","指定https证书文件")
	flagKey		= flag.String("k","key.pem","指定证书私钥文件")
)

func banner() {
	t := `
  _      __        _            __        
 | | /| / / ___   (_) _    __  / /   __ __
 | |/ |/ / / -_) / / | |/|/ / / _ \ / // /
 |__/|__/  \__/ /_/  |__,__/ /_//_/ \_, / 
                                   /___/  
`
	fmt.Println(t)
}

func main() {
	var  cert,key string
	cert="-----BEGIN CERTIFICATE-----\nMIIDtzCCAp+gAwIBAgIIB5iK9ibdiJwwDQYJKoZIhvcNAQELBQAwgYoxCzAJBgNV\nBAYTAkNOMRcwFQYDVQQJEw5Qb3J0U3dpZ2dlciBDQTEXMBUGA1UEERMOUG9ydFN3\naWdnZXIgQ0ExFzAVBgNVBAoTDlBvcnRTd2lnZ2VyIENBMRcwFQYDVQQLEw5Qb3J0\nU3dpZ2dlciBDQTEXMBUGA1UEAxMOUG9ydFN3aWdnZXIgQ0EwHhcNMjIwMjIzMDIz\nNzU3WhcNMjMwMjIzMDIzNzU3WjAhMQswCQYDVQQGEwJDTjESMBAGA1UEAxMJbG9j\nYWxob3N0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmAdHzwJW4m3/\nK0iZhbYnkqwc0FwQpvMdhIxVMxHquFVm6DkkCIh5yf5p2sxoLNuYjmaFD6ujNrOo\njocaxbx1WTVHYCPM9uqGw5eeG+shYJ0TLyZ0ukLpauGfnBcqN0601U8Nm3Opdhbg\n1J9EfdQxsemWFpkqyLT7sqW1U5DR3WeSxqERjs+0BVZ14q9xsSIv0CfSSTv+Xggr\nhBOXSxA5dOgSn7xqH9YBlTIezzTbK9tWmlzcCQEsgtDUvqSz/8R/QiEKYVB9C4xD\nJcX3dm9758Ca2lEJxLUVg/t9sIAEKahYYJwYoOtR1q/q376xY/AmWsFXZcrhlMrm\nmwMfJj2UUQIDAQABo4GIMIGFMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggr\nBgEFBQcDAQYIKwYBBQUHAwIwHQYDVR0OBBYEFP9kui/48VE138QPSQSjTez/i8we\nMB8GA1UdIwQYMBaAFHG1AJAmc2VRMaKvllq48qSvYj3PMBQGA1UdEQQNMAuCCWxv\nY2FsaG9zdDANBgkqhkiG9w0BAQsFAAOCAQEAaYLPKgx6PfUp9eAqKWIQgYLWVW6K\nvBD/bAvsgAYf0FfplB0O8O0lr7ReheXhkaXma9zbrlOM5PGnA4laBWu6uPjrOeRP\n6I2FmlLev/gdWAKb3FnKxsM8xHdPMpw6QeKCZlvJXaIgFcl7sTtVdmE8hTGL0Vk9\nvebOlq44f09/M38PmxQjBXhDvR3d7XTeKFbGxwJMZws6/SoRESYWmtzvqFwGvQOZ\nXl8tFBwjCnNK0Tahx++ceyNvtq0D3gge4gjeQihcFxmAHetIujrE7boIvDO2iSsD\nBVgklmQ/9EMopFw9XZ0r2fLgVwRbPMi2livRIcpp9QeyBbnhjDjoYVfUGQ==\n-----END CERTIFICATE-----\n"
	key="-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAmAdHzwJW4m3/K0iZhbYnkqwc0FwQpvMdhIxVMxHquFVm6Dkk\nCIh5yf5p2sxoLNuYjmaFD6ujNrOojocaxbx1WTVHYCPM9uqGw5eeG+shYJ0TLyZ0\nukLpauGfnBcqN0601U8Nm3Opdhbg1J9EfdQxsemWFpkqyLT7sqW1U5DR3WeSxqER\njs+0BVZ14q9xsSIv0CfSSTv+XggrhBOXSxA5dOgSn7xqH9YBlTIezzTbK9tWmlzc\nCQEsgtDUvqSz/8R/QiEKYVB9C4xDJcX3dm9758Ca2lEJxLUVg/t9sIAEKahYYJwY\noOtR1q/q376xY/AmWsFXZcrhlMrmmwMfJj2UUQIDAQABAoIBACX2Y3MIxYNjLj2z\nskpWUwloGwXYh3v4510K4deP2MnQ4ZKJejr7QVY0RmCRfE3/Q9gDN6TTGj11nViB\n2iiYR8FP8ZPLPMpHhAFhYeTc0QfcEUdL1ggQ31KGJqIGks8ewz8kr2Qq/Jz1V40g\nfCzjiMlBZ+4arzlRICza/i1w00byzIb4Mc+WpLigLAW+z7dvkwuPjLDgkERTzGO3\nBSh7dMyu+kPFtd1peokfV2lFveVsl6aZ0TldV6Yj+vn3JMyLRvYTSHJMT7+sLVKV\nbSdtNgZnqu6rPfBy0LsdFSL47dzroCa+CwHjxytMHsTEPrjPfGCaEme2i8z300se\n2ELkIJUCgYEAx20s6Vnra45mJMxuFbhbXUbTNBlDaMYoFxrN4pJfJ0iAxwbTkBsb\nsx1j17vd7CR9psWrsjRoVy8ciw3//GrQsfLJW7eQDFAf2TuKorpr8s1TkvL+F/15\n07WDHwoMRbCaGD3tbYB7TOS3p31fvxT9sK5TB8hWa/9NBP+YasZxSHcCgYEAwyfy\nkp1ZegSSyrqy19CNt6jPla+9f3m+36+Q75Amk+5dxPnxwh0kLhvLdBpAX6gjdj7S\nzIQqsJLchSqEau9yXNyMoEspU67c+1pdn/bvkV5qcclpXgCZRP+zu0W6xDdcyyp8\nNSWubz05cnjNjTC0ppKY2DbeMuuk9CkkOeyAg3cCgYEAuSxZmG+jFmLJ239q9IYT\nZ+AcunJ+0im1DgcYhzosWniLAsMG7PcO2FkA3U+W9+7GMXW4QKdC/zlCgqfEOugf\n0ivC5oPioFxBxl7wEruRAom/AWLZpwL4/Si8kLSuwoMCAmJ2NMgpNTPmiAH1RqNC\nEK09cauG+466QFrodrS+eW8CgYEAtqAedyd6gR3ghIicLivDQvhdcfVofu9eOJNi\nLV0XTN4Gr/s6Su3cWq22AetDDIEa1l/KAos4K87zQEbebfevbMkTbDmJ7f64Wxtg\ne/6oh7/0GpFh4g7rX09uUTTapx4r3w/d9hpSF1US+lWj/ZlzpGGRWNxXqQ0tazNI\n7E+un9cCgYBnfr1+jHkmHZT0RRCGA7d5Wq4er4xAgQAvOlT9cOYpBhyAZtyWrnMx\nP15kSH8xJfwmnyuR/RyMLiaYFCYig5dS33d7NX0WKrW6717D2L1xn8f6Xq0GdO7V\nM++5vFJs1xyzlnWPHjmorb04Oo8/zNbvbs+2EL/0uf/hBwyaUTTpxg==\n-----END RSA PRIVATE KEY-----\n"
	banner()
	flag.Parse()
	if *flagHelp  {
		fmt.Printf("Usage: \n\n")
		flag.PrintDefaults()
		return
	}
	http.Handle("/",http.FileServer(http.Dir(*flagDir)))
	if *flagSSL {
		if *flagLis != "0.0.0.0:80" {
			fmt.Println("Staring Listen ",*flagLis)
		}else {
			*flagLis="0.0.0.0:443"
			fmt.Println("Staring Listen ",*flagLis)
		}
		if *flagCert=="cert.pem" || *flagKey=="key.pem" {
			cer :=[]byte(cert)
			ke  :=[]byte(key)
			ioutil.WriteFile(*flagCert,cer,0644)
			ioutil.WriteFile(*flagKey,ke,0644)
		}
		err :=http.ListenAndServeTLS(*flagLis, *flagCert,
			*flagKey, nil)
		if err!=nil {
			log.Fatal(err)
		}
	}else {
		fmt.Println("Staring Listen ",*flagLis)
		err := http.ListenAndServe(*flagLis,nil)
		if err!=nil {
			log.Fatal(err)
		}
	}
}

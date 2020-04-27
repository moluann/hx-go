package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var srcfileName string
var outFileName string
var encrypt bool
func init() {
	flag.StringVar(&srcfileName,"i","","input file")
	flag.StringVar(&outFileName,"o","","output file")
	flag.BoolVar(&encrypt,"e",true,"ture=encrypt,false=decrypt")

}

func main() {
	flag.Parse()
	if  srcfileName =="" || outFileName == ""{
		fmt.Fprint(os.Stderr,"Please type input and output name!")
	}
	file, err := ioutil.ReadFile(srcfileName)
	if err != nil {
		fmt.Fprint(os.Stderr,err.Error())
		return
	}
	key:=[]byte("weifeng201812345")
	e := make([]byte,20)
	if encrypt {
		e = encryptAES(file, key)
	}else {
		e = decryptAES(file,key)
	}

	//f,_ :=os.OpenFile(outFileName,os.O_CREATE,0644)
	//f.Write(e)
	ioutil.WriteFile(outFileName,e,0777)

}



func padding(src []byte,blocksize int) []byte {
	padnum:=blocksize-len(src)%blocksize
	pad:=bytes.Repeat([]byte{byte(padnum)},padnum)
	return append(src,pad...)
}

func unpadding(src []byte) []byte {
	n:=len(src)
	unpadnum:=int(src[n-1])
	return src[:n-unpadnum]
}

func encryptAES(src []byte,key []byte) []byte {
	block,_:=aes.NewCipher(key)
	src=padding(src,block.BlockSize())
	blockmode:=cipher.NewCBCEncrypter(block,key)
	blockmode.CryptBlocks(src,src)
	return src
}

func decryptAES(src []byte,key []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockmode := cipher.NewCBCDecrypter(block, key)
	blockmode.CryptBlocks(src, src)
	src = unpadding(src)
	return src

}
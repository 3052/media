package paramount

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/base64"
   "encoding/hex"
)

// 16.0.0
const ComCbsCa AppSecret = "6c68178445de8138"

// 16.0.0
const ComCbsApp AppSecret = "9fc14cb03691c342"

func (a AppSecret) At() (At, error) {
   key, err := hex.DecodeString(secret_key)
   if err != nil {
      return "", err
   }
   block, err := aes.NewCipher(key)
   if err != nil {
      return "", err
   }
   var iv [aes.BlockSize]byte
   data := []byte{'|'}
   data = append(data, a...)
   data = pad(data)
   cipher.NewCBCEncrypter(block, iv[:]).CryptBlocks(data, data)
   data1 := []byte{0, aes.BlockSize}
   data1 = append(data1, iv[:]...)
   data1 = append(data1, data...)
   return At(base64.StdEncoding.EncodeToString(data1)), nil
}

type At string

func pad(data []byte) []byte {
   length := aes.BlockSize - len(data)%aes.BlockSize
   for high := byte(length); length >= 1; length-- {
      data = append(data, high)
   }
   return data
}

const secret_key = "302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"

type AppSecret string

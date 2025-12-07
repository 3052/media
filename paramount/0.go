package paramount

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/base64"
   "encoding/binary"
   "encoding/hex"
   "slices"
)

const secret_key = "302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"

// pkcs7Pad appends padding to the data according to the PKCS #7 standard.
func pkcs7Pad(data []byte, blockSize int) []byte {
   // Calculate the number of padding bytes needed.
   // If data is already a multiple of blockSize, this results in a full block of padding.
   paddingLen := blockSize - (len(data) % blockSize)
   // Create a padding byte (the value is the length of the padding)
   padByte := byte(paddingLen)
   // Append the padding byte 'paddingLen' times
   for i := 0; i < paddingLen; i++ {
      data = append(data, padByte)
   }
   return data
}

func GetAt(appSecret string) (string, error) {
   // 1. Decode hex secret key
   key, err := hex.DecodeString(secret_key)
   if err != nil {
      return "", err
   }
   // 2. Create aes cipher with key
   block, err := aes.NewCipher(key)
   if err != nil {
      return "", err
   }
   // 3 & 4. Create payload: "|" + appSecret
   data := []byte("|" + appSecret)
   // 5. Apply PKCS7 Padding (Separate Function)
   data = pkcs7Pad(data, aes.BlockSize)
   // Prepare Empty IV (16 bytes of zeros)
   iv := make([]byte, aes.BlockSize)
   // 6. CBC encrypt with empty IV
   // We encrypt 'data' in place
   cipher.NewCBCEncrypter(block, iv).CryptBlocks(data, data)
   // 8. Create Header for block size (uint16)
   sizeHeader := make([]byte, 2)
   binary.BigEndian.PutUint16(sizeHeader, uint16(aes.BlockSize))
   // 7 & 8. Combine [Size] + [IV] + [Encrypted Data]
   finalBuffer := slices.Concat(sizeHeader, iv, data)
   // 9. Return result base64 encoded
   return base64.StdEncoding.EncodeToString(finalBuffer), nil
}

type Provider struct {
   AppSecret string
   Version   string
}

var ComCbsApp = Provider{
   AppSecret: "9fc14cb03691c342",
   Version:   "16.0.0",
}

var ComCbsCa = Provider{
   AppSecret: "6c68178445de8138",
   Version:   "16.0.0",
}

package wallet

import (
	"crypto/ecdsa"
)

const version = byte(0x00)

// A structure that represents a wallet to access the blockchain
type Wallet struct {
	PrivateKey ecdsa.PrivateKey // Represetns the Public Key of the Wallet
	PublicKey  []byte           // Represents the Private Key of the Wallet
}

// A constructor function that generates and returns a Wallet
func NewWallet() *Wallet {
	// Generate private-public key pair
	private, public := GenerateWalletKeys()
	// Assign the keys to the wallet fields
	wallet := Wallet{PrivateKey: private, PublicKey: public}
	// Return the wallet
	return &wallet
}

/*
A method of Wallet that generates and returns the address
of the wallet as a slice of bytes.

The address generation is done by generating the hash of the public key
and the check sum of this public hash and concatentating these with a version.
The total hash is then encoded into base58 which creates the address.
*/
func (w *Wallet) GenerateAddress() []byte {
	// Generate the hash of the public key
	publickeyhash := GeneratePublicKeyHash(w.PublicKey)
	// Generate the checksum hash
	checksumhash := GeneratePublicKeyCheckSum(publickeyhash)

	// Add the version byte, public key hash and checksum hash
	finalhash := append([]byte{version}, publickeyhash...)
	finalhash = append(finalhash, checksumhash...)

	// Encode the final hash to base58
	address := Base58Encode(finalhash)
	// Return the address
	return address
}

/*
A function that generates the hash of the public key.

The hash of the public key is generated by double hashing the public key with
SHA256 and then truncated to 20 bytes. This equivalent to the Bitcoin method of
using SHA256 followed by RIPEMD160 and is equally safe.
The RIPEMD160 algorithm was rejected because of its deprecated nature in the golang
crypto library along with the fact that truncated SHA2 algortithms are equally safe
Hence the choice to use truncated SHA256 over the SHA1 which also outputs 20 bytes.

References:

https://bitcoin.stackexchange.com/questions/16543/why-was-the-ripemd-160-hash-algorithms-chosen-before-sha-1

https://crypto.stackexchange.com/questions/3153/sha-256-vs-any-256-bits-of-sha-512-which-is-more-secure
*/
func GeneratePublicKeyHash(publickey []byte) []byte {
	// Generate the double hash of the public key
	publickeyhash := GenerateDoubleHash(publickey)
	// Truncate the public key hash to 20 bytes
	publickeyhash = TruncateHash(publickeyhash, 20)
	// Return the hash of the public key
	return publickeyhash
}

/*
A function that generates the checksum of the public key hash.

The checksum is obtained by double hashing the hashed public key with SHA256
and truncated to its first 4 bytes. This method is consistent with the methods
suggested in the original Bitcoin papers.
*/
func GeneratePublicKeyCheckSum(publickeyhash []byte) []byte {
	// Generate the checksum hash by double hashing the puclic key hash
	checksum := GenerateDoubleHash(publickeyhash)
	// Truncate the checksum hash to 4 bytes
	checksum = TruncateHash(checksum, 4)
	// Return the checksum of the public key hash
	return checksum
}
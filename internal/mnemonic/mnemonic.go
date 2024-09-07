package mnemonic

import (
	"github.com/tyler-smith/go-bip39"
)

func Generate() string {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	// Display mnemonic and keys
	return mnemonic
}

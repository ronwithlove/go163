package BLC

import (
	"fmt"
	"testing"
)

func TestNewWallt(t *testing.T) {
	wallet:=NewWallt()
	fmt.Printf("private key :%v\n",wallet.PrivateKey)
	fmt.Printf("public key :%v\n",wallet.PublicKey)
	fmt.Printf("wallet :%v\n",wallet)
}
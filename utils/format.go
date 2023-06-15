package utils

import (
	"fmt"
	"strings"
	"math/big"
)


func HexToBigInt(hexStr string) (*big.Int, error) {
	bigInt := new(big.Int)
	_, success := bigInt.SetString(hexStr[2:], 16)
	if !success {
		return nil, fmt.Errorf("failed to convert hex to big.Int")
	}
	return bigInt, nil
}


func IntToHex(num uint64) string {
	return fmt.Sprintf("0x%x", num)
}


func CheckAddressIsExists(stringAddressFrom string, stringAddressTo string, arrayAddressString []string) string {

	for _, v := range arrayAddressString {
		if strings.EqualFold(v, stringAddressFrom) || strings.EqualFold(v, stringAddressTo) {
			return v
		}
	}

	return ""

}
package main

import (
	"encoding/hex"
	"fmt"
	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/btcec"
	"github.com/huhongjia/bitcoin-abc-tech/config"
	"github.com/sirupsen/logrus"
)

func main() {
	for i := 0; i < 3; i++ {
		pk, err := btcec.NewPrivateKey(btcec.S256())
		if err != nil {
			logrus.Errorf("NewPrivateKey failed: %v", err)
			return
		}

		fmt.Println("------------------------------------")

		wif, err := cashutil.NewWIF((*btcec.PrivateKey)(pk),
			config.GetChainParam(), true)
		fmt.Println(wif.String())

		fmt.Println(hex.EncodeToString(wif.PrivKey.PubKey().SerializeCompressed()))

		w, _ := cashutil.DecodeWIF(wif.String())
		h := cashutil.Hash160(w.PrivKey.PubKey().SerializeCompressed())
		f, err := cashutil.NewAddressPubKeyHash(h, config.GetChainParam())
		fmt.Println(f.EncodeAddress(true))
	}
}

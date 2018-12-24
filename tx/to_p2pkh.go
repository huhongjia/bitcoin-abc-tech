package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/txscript"
	"github.com/huhongjia/bitcoin-abc-tech/basic"
	"github.com/huhongjia/bitcoin-abc-tech/config"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var conf = config.GetConf()

func main() {
	amount := decimal.NewFromFloat(0.1).Mul(decimal.NewFromFloat32(basic.UNIT))
	//FeeRate Unit fee/b= (fee/kb) * 1e8 / 1000
	feeRate := decimal.NewFromFloat(0.00002).Mul(decimal.NewFromFloat32(basic.UNIT)).Div(decimal.NewFromFloat32(1000))
	to := "bchtest:qqvk2qjqudp687azp2ln0nrdh7afehq0zgkp0kph25"
	privKey := "cUfp6zE37nwcqUSZjRtsuSTJ4No6gdLHF9A5LGLD4KZDiTXY5KTU"

	wif, err := cashutil.DecodeWIF(privKey)
	//P2PKH
	pkHash := cashutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed())
	from, err := cashutil.NewAddressPubKeyHash(pkHash, config.GetChainParam())

	utxos, err := basic.Query(from)
	if err != nil {
		logrus.Errorf("Query Utxo Error:%v", err)
		return
	}

	tx, usedUtxo, err := basic.AssemblyTx(utxos, amount, feeRate, from, to)
	if err != nil {
		logrus.Errorf("AssemblyTx Error:%v", err)
		return
	}

	script, _ := txscript.PayToAddrScript(from)
	tx, err = basic.SignTx(tx, wif, script, usedUtxo)
	if err != nil {
		logrus.Errorf("signTx Error:%v", err)
		return
	}

	buf := bytes.NewBuffer(nil)
	err = tx.Serialize(buf)
	txHex := hex.EncodeToString(buf.Bytes())

	fmt.Println(txHex)
}

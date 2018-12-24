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

func main() {

	wifStr := "cSrAZuRWUj2v4vXuByyQcD17eUXVNf1YjhEER3cnEe1gXhn4CeKz"
	w, _ := cashutil.DecodeWIF(wifStr)

	script, err := txscript.NewScriptBuilder().AddData(w.PrivKey.PubKey().SerializeCompressed()).AddOp(txscript.OP_CHECKSIG).Script()
	if err != nil {
		fmt.Printf("Build P2SH Error:%v", err)
		return
	}
	fmt.Println("////------------------------------------///")
	address, _ := cashutil.NewAddressScriptHash(script, config.GetChainParam())
	fromAddr := address.EncodeAddress(true)
	fmt.Println(fromAddr)

	amount := decimal.NewFromFloat(0.025).Mul(decimal.NewFromFloat32(basic.UNIT))
	//FeeRate Unit fee/b= (fee/kb) * 1e8 / 1000
	feeRate := decimal.NewFromFloat(0.00002).Mul(decimal.NewFromFloat32(basic.UNIT)).Div(decimal.NewFromFloat32(1000))
	to := "bchtest:qqvk2qjqudp687azp2ln0nrdh7afehq0zgkp0kph25"

	from, err := cashutil.DecodeAddress(fromAddr, config.GetChainParam())

	utxos, err := basic.QueryAddress(from.EncodeAddress(false))
	if err != nil {
		logrus.Errorf("Query Utxo Error:%v", err)
		return
	}

	tx, usedUtxo, err := basic.AssemblyTx(utxos, amount, feeRate, from, to)
	if err != nil {
		logrus.Errorf("AssemblyTx Error:%v", err)
		return
	}

	basic.SignP2SH(tx, w, script, usedUtxo)

	buf := bytes.NewBuffer(nil)
	err = tx.Serialize(buf)
	txHex := hex.EncodeToString(buf.Bytes())

	fmt.Println(txHex)
}

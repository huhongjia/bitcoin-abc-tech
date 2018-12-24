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

var wifStr = []string{
	"cSrAZuRWUj2v4vXuByyQcD17eUXVNf1YjhEER3cnEe1gXhn4CeKz",
	"cS4XRm39DZdAqCRwurdrNku1uXzbAh48A4PYBi5Sn1GWsW1fvGZf",
	"cS8aKjLhgXuZbWXQZu46DagsUqUyAj22njn3YjYHKoxEcJjkyiFG",
}

func main() {

	var pubKeys [][]byte
	var wifs []*cashutil.WIF
	for i := 0; i < len(wifStr); i++ {
		w, _ := cashutil.DecodeWIF(wifStr[i])

		wifs = append(wifs, w)

		pubKey := w.PrivKey.PubKey().SerializeCompressed()
		pubKeys = append(pubKeys, pubKey)
	}

	script, err := txscript.NewScriptBuilder().AddOp(txscript.OP_2).
		AddData(pubKeys[0]).AddData(pubKeys[1]).AddData(pubKeys[2]).
		AddOp(txscript.OP_3).AddOp(txscript.OP_CHECKMULTISIG).Script()
	if err != nil {
		fmt.Printf("Build P2SH Error:%v", err)
		return
	}
	fmt.Println("////------------------------------------///")
	address, _ := cashutil.NewAddressScriptHash(script, config.GetChainParam())
	fmt.Println(address.EncodeAddress(true))

	amount := decimal.NewFromFloat(0.025).Mul(decimal.NewFromFloat32(basic.UNIT))
	//FeeRate Unit fee/b= (fee/kb) * 1e8 / 1000
	feeRate := decimal.NewFromFloat(0.00002).Mul(decimal.NewFromFloat32(basic.UNIT)).Div(decimal.NewFromFloat32(1000))
	to := "bchtest:qqvk2qjqudp687azp2ln0nrdh7afehq0zgkp0kph25"
	fromAddr := "bchtest:prhasgpwhtpay9gznt6uwu8ge943dj7keu9cc7cjvx"

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

	basic.SignMultiTx(tx, wifs, script, usedUtxo)

	buf := bytes.NewBuffer(nil)
	err = tx.Serialize(buf)
	txHex := hex.EncodeToString(buf.Bytes())

	fmt.Println(txHex)

}

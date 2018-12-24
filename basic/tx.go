package basic

import (
	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/chaincfg/chainhash"
	"github.com/bcext/gcash/txscript"
	"github.com/bcext/gcash/wire"
	"github.com/copernet/go-electrum/electrum"
	"github.com/huhongjia/bitcoin-abc-tech/config"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const (
	defaultSequence      = 0xffffffff
	defaultOutputSize    = 34
	defaultSignatrueSize = 107
	dustSatoshi          = 546
)

func SignP2SH(tx *wire.MsgTx, wif *cashutil.WIF, script []byte, utxos []*electrum.Transaction) (*wire.MsgTx, error) {
	for idx, in := range tx.TxIn {
		sig, err := txscript.RawTxInSignature(tx, idx, script, cashutil.Amount(utxos[idx].Value),
			txscript.SigHashAll|txscript.SigHashForkID, wif.PrivKey)

		if err != nil {
			return nil, err
		}

		//txscript.NewScriptBuilder().AddOp(txscript.OP_1)
		sig, err = txscript.NewScriptBuilder().AddData(sig).AddData(script).Script()
		if err != nil {
			return nil, err
		}

		in.SignatureScript = sig

		// check whether signature is ok or not.
		engine, err := txscript.NewEngine(script, tx, idx, txscript.StandardVerifyFlags,
			nil, nil, utxos[idx].Value)
		if err != nil {
			return nil, err
		}
		// execution the script in stack
		err = engine.Execute()
		if err != nil {
			return nil, err
		}
	}

	return tx, nil
}

func SignMultiTx(tx *wire.MsgTx, wifs []*cashutil.WIF, pkScript []byte, utxos []*electrum.Transaction) (*wire.MsgTx, error) {
	for idx, in := range tx.TxIn {
		sig1, err := txscript.RawTxInSignature(tx, idx, pkScript, cashutil.Amount(utxos[idx].Value),
			txscript.SigHashAll|txscript.SigHashForkID, wifs[0].PrivKey)

		sig2, err := txscript.RawTxInSignature(tx, idx, pkScript, cashutil.Amount(utxos[idx].Value),
			txscript.SigHashAll|txscript.SigHashForkID, wifs[1].PrivKey)

		if err != nil {
			return nil, err
		}

		//txscript.NewScriptBuilder().AddOp(txscript.OP_1)
		sig, err := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(sig1).AddData(sig2).AddData(pkScript).Script()
		if err != nil {
			return nil, err
		}

		in.SignatureScript = sig

		// check whether signature is ok or not.
		engine, err := txscript.NewEngine(pkScript, tx, idx, txscript.StandardVerifyFlags,
			nil, nil, utxos[idx].Value)
		if err != nil {
			return nil, err
		}
		// execution the script in stack
		err = engine.Execute()
		if err != nil {
			return nil, err
		}
	}

	return tx, nil
}

func SignTx(tx *wire.MsgTx, wif *cashutil.WIF, pkScript []byte, utxos []*electrum.Transaction) (*wire.MsgTx, error) {
	for idx, in := range tx.TxIn {
		sig, err := txscript.RawTxInSignature(tx, idx, pkScript, cashutil.Amount(utxos[idx].Value),
			txscript.SigHashAll|txscript.SigHashForkID, wif.PrivKey)
		if err != nil {
			return nil, err
		}
		sig, err = txscript.NewScriptBuilder().AddData(sig).Script()
		if err != nil {
			return nil, err
		}

		pk, err := txscript.NewScriptBuilder().AddData(wif.PrivKey.PubKey().SerializeCompressed()).Script()
		if err != nil {
			return nil, err
		}
		sig = append(sig, pk...)
		in.SignatureScript = sig

		// check whether signature is ok or not.
		engine, err := txscript.NewEngine(pkScript, tx, idx, txscript.StandardVerifyFlags,
			nil, nil, utxos[idx].Value)
		if err != nil {
			return nil, err
		}
		// execution the script in stack
		err = engine.Execute()
		if err != nil {
			return nil, err
		}
	}

	return tx, nil

}

func AssemblyTx(utxos []*electrum.Transaction, amount decimal.Decimal, feeRate decimal.Decimal, from cashutil.Address, to string) (*wire.MsgTx, []*electrum.Transaction, error) {
	var tx wire.MsgTx
	tx.Version = 1 //TODO Version
	tx.LockTime = 0
	//Build Pay Output
	ad, err := cashutil.DecodeAddress(to, config.GetChainParam())
	if err != nil {
		return nil, nil, err
	}

	script, _ := txscript.PayToAddrScript(ad)
	tx.TxOut = append(tx.TxOut, &wire.TxOut{PkScript: script, Value: amount.IntPart()})

	//Build Input
	var totalAmount = int64(0)
	var fee = int64(0)
	usedUtxo := make([]*electrum.Transaction, 0)
	for i, u := range utxos {
		hash, _ := chainhash.NewHashFromStr(u.Hash)
		txIn := wire.TxIn{
			PreviousOutPoint: *wire.NewOutPoint(hash, uint32(u.Pos)),
			Sequence:         defaultSequence, //TODO Sequence
		}
		tx.TxIn = append(tx.TxIn, &txIn)
		totalAmount = totalAmount + u.Value

		//Cal miner fee
		fee = feeRate.IntPart() * int64(tx.SerializeSize()+defaultSignatrueSize*(i+1)+defaultOutputSize)
		usedUtxo = append(usedUtxo, u)
		if totalAmount-(amount.IntPart()+fee) > 0 {
			break
		}
	}

	diff := totalAmount - (amount.IntPart() + fee)
	if diff <= 0 {
		return nil, nil, errors.New("Utxo is not enough")
	}

	if diff > dustSatoshi {
		//Build Back output
		script, _ := txscript.PayToAddrScript(from)
		tx.TxOut = append(tx.TxOut, &wire.TxOut{PkScript: script, Value: diff})
	}

	return &tx, usedUtxo, nil
}

package basic

import (
	"github.com/bcext/cashutil"
	"github.com/copernet/go-electrum/electrum"
	"github.com/huhongjia/bitcoin-abc-tech/config"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	UNIT = 1e8
)

var conf = config.GetConf()

func Query(addr *cashutil.AddressPubKeyHash) ([]*electrum.Transaction, error) {
	addrStr := addr.EncodeAddress(false)
	return getNode().BlockchainAddressListUnspent(addrStr)
}

func QueryAddress(address string) ([]*electrum.Transaction, error) {
	return getNode().BlockchainAddressListUnspent(address)
}

func getNode() *electrum.Node {

	n := electrum.NewNode()
	if err := n.ConnectTCP(conf.Electron.Host + ":" + conf.Electron.Port); err != nil {
		logrus.Errorf("create connection to electrum error: %v", err)
		os.Exit(1)
	}

	return n
}

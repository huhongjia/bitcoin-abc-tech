package block

import (
	"fmt"
	"github.com/bcext/gcash/wire"
	"os"
	"testing"
)

func TestHashHeader(t *testing.T) {

	filePath := "/Users/hongjia.hu/Documents/bitmain/data/file"
	file, _ := os.Open(filePath)
	//For block init
	header, _ := loadCurrHeader(file)
	fmt.Printf("Version:%d", header.Version)
	fmt.Printf("Timestamp:%d", header.Timestamp.Unix())
	fmt.Printf("Bits:%d", header.Bits)
	fmt.Printf("Nonce:%d", header.Nonce)
	fmt.Printf("MerkleRoot:%s", header.MerkleRoot.String())
	fmt.Printf("PrevBlock:%s", header.PrevBlock.String())
}

func loadCurrHeader(file *os.File) (*wire.BlockHeader, error) {

	fi, err := file.Stat()
	if err != nil {
		// Could not obtain stat, handle error
		return nil, err
	}

	_, err = file.Seek(fi.Size()-80, 0)
	if err != nil {
		return nil, err
	}

	var header wire.BlockHeader
	err = header.Deserialize(file)
	if err != nil {
		return nil, err
	}

	return &header, nil
}

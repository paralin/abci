package testsuite

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"

	crypto "github.com/libp2p/go-libp2p-crypto"
	abcicli "github.com/tendermint/abci/client"
	"github.com/tendermint/abci/types"
	gdata "github.com/tendermint/go-wire/data"
	cmn "github.com/tendermint/tmlibs/common"
)

func InitChain(client abcicli.Client) error {
	total := 10
	vals := make([]*types.Validator, total)
	for i := 0; i < total; i++ {
		r := rand.New(rand.NewSource(int64(1035132 + (i * 1337))))
		_, pubKey, err := crypto.GenerateEd25519Key(r)
		if err != nil {
			return err
		}

		pubKeyBin, err := pubKey.Bytes()
		if err != nil {
			return err
		}

		pubKeyStr, err := gdata.Encoder.Marshal(pubKeyBin)
		if err != nil {
			return err
		}

		power := cmn.RandInt()
		vals[i] = &types.Validator{string(pubKeyStr), int64(power)}
	}
	_, err := client.InitChainSync(types.RequestInitChain{Validators: vals})
	if err != nil {
		fmt.Printf("Failed test: InitChain - %v\n", err)
		return err
	}
	fmt.Println("Passed test: InitChain")
	return nil
}

func SetOption(client abcicli.Client, key, value string) error {
	res, err := client.SetOptionSync(types.RequestSetOption{Key: key, Value: value})
	log := res.GetLog()
	if err != nil {
		fmt.Println("Failed test: SetOption")
		fmt.Printf("setting %v=%v: \nlog: %v\n", key, value, log)
		fmt.Println("Failed test: SetOption")
		return err
	}
	fmt.Println("Passed test: SetOption")
	return nil
}

func Commit(client abcicli.Client, hashExp []byte) error {
	res, err := client.CommitSync()
	_, data := res.Code, res.Data
	if err != nil {
		fmt.Println("Failed test: Commit")
		fmt.Printf("committing %v\nlog: %v\n", res.GetLog(), err)
		return err
	}
	if !bytes.Equal(data, hashExp) {
		fmt.Println("Failed test: Commit")
		fmt.Printf("Commit hash was unexpected. Got %X expected %X\n",
			data.Bytes(), hashExp)
		return errors.New("CommitTx failed")
	}
	fmt.Println("Passed test: Commit")
	return nil
}

func DeliverTx(client abcicli.Client, txBytes []byte, codeExp uint32, dataExp []byte) error {
	res, _ := client.DeliverTxSync(txBytes)
	code, data, log := res.Code, res.Data, res.Log
	if code != codeExp {
		fmt.Println("Failed test: DeliverTx")
		fmt.Printf("DeliverTx response code was unexpected. Got %v expected %v. Log: %v\n",
			code, codeExp, log)
		return errors.New("DeliverTx error")
	}
	if !bytes.Equal(data, dataExp) {
		fmt.Println("Failed test: DeliverTx")
		fmt.Printf("DeliverTx response data was unexpected. Got %X expected %X\n",
			data, dataExp)
		return errors.New("DeliverTx error")
	}
	fmt.Println("Passed test: DeliverTx")
	return nil
}

func CheckTx(client abcicli.Client, txBytes []byte, codeExp uint32, dataExp []byte) error {
	res, _ := client.CheckTxSync(txBytes)
	code, data, log := res.Code, res.Data, res.Log
	if code != codeExp {
		fmt.Println("Failed test: CheckTx")
		fmt.Printf("CheckTx response code was unexpected. Got %v expected %v. Log: %v\n",
			code, codeExp, log)
		return errors.New("CheckTx")
	}
	if !bytes.Equal(data, dataExp) {
		fmt.Println("Failed test: CheckTx")
		fmt.Printf("CheckTx response data was unexpected. Got %X expected %X\n",
			data, dataExp)
		return errors.New("CheckTx")
	}
	fmt.Println("Passed test: CheckTx")
	return nil
}

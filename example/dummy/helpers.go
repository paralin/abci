package dummy

import (
	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/tendermint/abci/types"
	"github.com/tendermint/go-wire/data"
	cmn "github.com/tendermint/tmlibs/common"
	"math/rand"
)

// RandVal creates one random validator, with a key derived
// from the input value
func RandVal(i int) *types.Validator {
	// keyBin := []byte(cmn.Fmt("test%d", i))
	r := rand.New(rand.NewSource(int64(1035132 + (i * 1337))))
	_, pubKey, err := crypto.GenerateEd25519Key(r)
	if err != nil {
		panic(err)
	}

	pubKeyBin, err := pubKey.Bytes()
	if err != nil {
		panic(err)
	}

	pubKeyStr, err := data.Encoder.Marshal(pubKeyBin)
	if err != nil {
		panic(err)
	}

	power := cmn.RandUint16() + 1
	return &types.Validator{string(pubKeyStr), int64(power)}
}

// RandVals returns a list of cnt validators for initializing
// the application. Note that the keys are deterministically
// derived from the index in the array, while the power is
// random (Change this if not desired)
func RandVals(cnt int) []*types.Validator {
	res := make([]*types.Validator, cnt)
	for i := 0; i < cnt; i++ {
		res[i] = RandVal(i)
	}
	return res
}

// InitDummy initializes the dummy app with some data,
// which allows tests to pass and is fine as long as you
// don't make any tx that modify the validator state
func InitDummy(app *PersistentDummyApplication) {
	app.InitChain(types.RequestInitChain{RandVals(1)})
}

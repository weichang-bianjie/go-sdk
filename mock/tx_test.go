package mock

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/binance-chain/go-sdk/client"
	"github.com/binance-chain/go-sdk/client/rpc"
	"github.com/binance-chain/go-sdk/client/transaction"
	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	"math"
	"strings"
	"testing"
	"time"
)

//var (
//	networkType = types.ProdNetwork
//	dexUrl      = "dex.binance.org"
//	nodeUrl     = "tcp://seed1.longevito.io:80"
//)

var (
	networkType = types.TestNetwork
	dexUrl      = "testnet-dex.binance.org"
	nodeUrl     = "tcp://data-seed-pre-0-s1.binance.org:80"
)

var (
	ksFilePath = "./BD-testnet-ks_2_Admin@123456.txt"
	ksAuth     = "Admin@123456"
	dexClient  client.DexClient
	rpcClient  rpc.Client
)

func TestMain(m *testing.M) {
	if km, err := keys.NewKeyStoreKeyManager(ksFilePath, ksAuth); err != nil {
		panic(err)
	} else {
		if c, err := client.NewDexClient(dexUrl, networkType, km); err != nil {
			panic(err)
		} else {
			fmt.Printf("km address is %s\n", km.GetAddr().String())
			dexClient = c
		}

		rpcClient = rpc.NewRPCClient(nodeUrl, networkType)
		if _, err := rpcClient.Status(); err != nil {
			fmt.Printf("init rpc client fail, err is %s\n", err.Error())
			panic(err)
		}
	}
	m.Run()
}

func TestRecoverFromKeyStore(t *testing.T) {
	file := ""
	if file == "" {
		file = ksFilePath
	}
	auth := "Admin@123456"
	if km, err := keys.NewKeyStoreKeyManager(file, auth); err != nil {
		t.Fatal(err)
	} else {
		content := []byte("Testing")
		if signedBytes, err := km.GetPrivKey().Sign(content); err != nil {
			t.Fatal(err)
		} else {
			if km.GetPrivKey().PubKey().VerifyBytes(content, signedBytes) {
				t.Log("verify signed bytes success")
			}

			if phrase, err := km.ExportAsPrivateKey(); err != nil {
				t.Fatal(err)
			} else {
				t.Logf("phrase is: %s\n", phrase)
			}
		}
	}

}

func TestHTLTRefund(t *testing.T) {
	swapId := "aa9c9eb100d7b0340c6f35949398b1103a9052ef90afac8c4ade011f3ca29074"
	swapIdBytes, err := hex.DecodeString(swapId)
	if err != nil {
		t.Fatal(err)
	}
	option := transaction.WithMemo("")
	sendResult, err := dexClient.RefundHTLT(swapIdBytes, true, option)
	if err != nil {
		t.Fatalf("refund htlt tx failed, swapId: %s, err: %s\n", swapId, err.Error())
	} else {
		if sendResult.Ok {
			fmt.Printf("refund htlt tx success, txHash is %s\n", sendResult.Hash)
		} else {
			fmt.Printf("refund htlt tx fail, txHash is %s, log is %s\n",
				sendResult.Hash, sendResult.Log)
		}
	}
}

func TestHTLTCreate(t *testing.T) {
	//recipient := ""
	//recipientOnOtherChain := ""
	//senderOnOtherChain := ""
	//hashLock := ""
	//timestamp := 0
	//amt := ""
	//expectedIncoing := ""
	//heightSpan := 0
	//crossChain := true
	//sync := true
	//option := transaction.WithMemo("")
	//
	//getAccAddr := func(addr string) types.AccAddress {
	//	if v, err := types.AccAddressFromBech32(addr); err != nil {
	//		t.Fatalf("get acc addr failed, err: %s\n", err.Error())
	//		return nil
	//	} else {
	//		return v
	//	}
	//}
	//
	//genHashLock := func() {
	//
	//}
	//
	//dexClient.HTLT(getAccAddr(recipient), getAccAddr(recipientOnOtherChain), getAccAddr(senderOnOtherChain),
	//	)
}

func TestSendToken(t *testing.T) {
	var (
		msgs []msg.Transfer
	)
	denom := "XRP-FB8"

	receivers := map[string]float64{
		"tbnb1la9uz527hwv78tz7jrdmheu6fx3m3gevkqc2zm": 1,
	}

	for k, v := range receivers {
		if toAddr, err := types.AccAddressFromBech32(k); err != nil {
			t.Fatalf("invalid addr: %s\n", toAddr)
		} else {
			var coins []types.Coin
			coin := types.Coin{
				Denom:  denom,
				Amount: int64(v * math.Pow10(8)),
			}
			coins = append(coins, coin)

			msgSendToken := msg.Transfer{
				ToAddr: toAddr,
				Coins:  coins,
			}
			msgs = append(msgs, msgSendToken)
		}
	}

	for _, v := range msgs {
		msgSendToken := []msg.Transfer{v}
		option := transaction.WithMemo("congratulation, you got airdrop from IRISnet")

		if sendResult, err := dexClient.SendToken(msgSendToken, true, option); err != nil {
			fmt.Printf("send token occur error, toAddr is %s, err is %s\n", v.ToAddr.String(), err.Error())
		} else {
			if sendResult.Ok {
				fmt.Printf("send token success, toAddr is %s, txHash is %s\n", v.ToAddr.String(), sendResult.Hash)
			} else {
				fmt.Printf("send token fail, toAddr is %s, txHash is %s, log is %s\n", v.ToAddr.String(),
					sendResult.Hash, sendResult.Log)
			}
		}
		fmt.Println("now sleep 5 seconds")
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func TestGetTxDetail(t *testing.T) {
	txHash := "A8A5A57DCB0AF1FCEB2A4114F88067AB48E11D45C5E7DE66504969AD437DF4FB"
	txHashBytes, err := hex.DecodeString(txHash)
	if err != nil {
		t.Fatal(err)
	}
	resultTx, err := rpcClient.Tx(txHashBytes, false)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("tx result is:\n %s\n", marshalJsonIgnoreError(resultTx))
	}

	txDataStr := getTxDataStr(resultTx.Tx.String())
	txData, err := hex.DecodeString(txDataStr)
	if err != nil {
		t.Fatal(err)
	}

	if stdTx, err := parseTxToStdTx(txData); err != nil {
		t.Fatal(err)
	} else {
		msgs := stdTx.Msgs
		if len(msgs) > 0 {
			txMsg := msgs[0]
			switch txMsg.(type) {
			case msg.SendMsg:
				txMsg := txMsg.(msg.SendMsg)
				t.Logf("tx msg is:\n %s\n", marshalJsonIgnoreError(txMsg))
				break
			default:
				t.Log("unknown tx msg")
			}
		}
	}
}

func getTxDataStr(txStr string) string {
	prefix := "Tx{"
	suffix := "}"
	return strings.TrimSuffix(strings.TrimPrefix(txStr, prefix), suffix)
}

func parseTxToStdTx(txBytes []byte) (tx.StdTx, error) {
	var txInfo tx.StdTx
	txStructure, err := rpc.ParseTx(tx.Cdc, txBytes)
	if err != nil {
		return txInfo, err
	}

	switch txStructure.(type) {
	case tx.StdTx:
		txInfo = txStructure.(tx.StdTx)
		return txInfo, nil
	default:
		return txInfo, fmt.Errorf("unkonwn txStructure")
	}
}

func marshalJsonIgnoreError(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

package mock

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/binance-chain/go-sdk/client"
	"github.com/binance-chain/go-sdk/client/rpc"
	"github.com/binance-chain/go-sdk/client/transaction"
	"github.com/binance-chain/go-sdk/common/types"
	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	"io"
	"math"
	"net/http"
	"strings"
	"testing"
	//"time"
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
	//ksFilePath = "/Users/user/Downloads/testnet_tbnb1cegl4x48qy6mq5vg5wtryk806n2vjtyhk3sj6v.json"
	//ksAuth     = "1234567890"

	ksFilePath = "/Users/user/Downloads/QA_Binance_Withdraw_Admin@123456_keystore.txt"
	ksAuth     = "Admin@123456"
	dexClient  client.DexClient
	rpcClient  rpc.Client
)

func TestMain(m *testing.M) {
	if km, err := keys.NewKeyStoreKeyManager(ksFilePath, ksAuth); err != nil {
		panic(err)
	} else {
		//if c, err := client.NewDexClient(dexUrl, networkType, km); err != nil {
		//	panic(err)
		//} else {
		//	fmt.Printf("km address is %s\n", km.GetAddr().String())
		//	dexClient = c
		//}

		rpcClient = rpc.NewRPCClient(nodeUrl, networkType)
		rpcClient.SetKeyManager(km)
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
	swapId := "a5ae994b5f7d0a49a690ae995c9321bb8c556715501f34fdca7215f4ef47ceec"
	swapIdBytes, err := hex.DecodeString(swapId)
	if err != nil {
		t.Fatal(err)
	}
	option := transaction.WithMemo("")
	sendResult, err := rpcClient.RefundHTLT(swapIdBytes, rpc.Sync, option)
	if err != nil {
		t.Fatalf("refund htlt tx failed, swapId: %s, err: %s\n", swapId, err.Error())
	} else {
		//if sendResult.Ok {
		fmt.Printf("refund htlt tx success, txHash is %s\n", sendResult.Hash)
		//} else {
		//	fmt.Printf("refund htlt tx fail, txHash is %s, log is %s\n",
		//		sendResult.Hash, sendResult.Log)
		//}
	}
}

func Post(url string, body string) (bz []byte, err error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode != 200, code:%v", resp.StatusCode)
	}
	defer resp.Body.Close()

	bz, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return
}

func TestHTLTCreate(t *testing.T) {

	var data struct {
		Data struct {
			Id        uint64 `json:"id"`
			Hashlock  string `json:"hashlock"`
			Timestamp uint64 `json:"timestamp"`
		} `json:"data"`
	}
	bytedata, err := Post("http://localhost:8081/htlts", "{\"address\":\"iaa1eqvkfthtrr93g4p9qspp54w6dtjtrn27ar7rpw\",\"amount\":{\"denom\":\"BNB\",\"amount\":100000},\"sc_chain\":\"binance\",\"dc_chain\":\"iris\"}")
	if err != nil {
		t.Error(err.Error())
		return
	}

	err = json.Unmarshal(bytedata, &data)
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(string(bytedata))

	//amount := ctypes.Coin{"BUSD-BAF", 100000}
	//deputy := "iaa1l7qf5gjxsygp9a59dm2xsjjwdelrf7yne5rl20"

	amount := ctypes.Coin{"BNB", 100000}
	deputy := "iaa1ljemm0yznz58qxxs8xyak7fashcfxf5lgl4zjx"

	recipient, _ := ctypes.AccAddressFromBech32("tbnb1reve5wwfes5u3r8hq8ukm6jwl5wtpfawxgshpr")
	recipientOtherChain := "iaa1eqvkfthtrr93g4p9qspp54w6dtjtrn27ar7rpw"
	senderOtherChain := deputy
	randomNumberHash, _ := hex.DecodeString(data.Data.Hashlock)
	timestamp := int64(data.Data.Timestamp)
	expectedIncome := fmt.Sprintf("%d%s", amount.Amount, amount.Denom)
	heightSpan := int64(360)

	val, err := dexClient.HTLT(recipient, recipientOtherChain, senderOtherChain, randomNumberHash, timestamp, ctypes.Coins{amount}, expectedIncome, heightSpan, true, true)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(val.Hash)
	//_, err = Post("http://localhost:8081/txs", fmt.Sprintf("{\"htlt_id\":%d,\"tx\":\"%s\",\"chain\":\"binance\"}", data.Data.Id, hex.EncodeToString(val)))
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
}

func TestSendToken(t *testing.T) {
	var (
		msgs []msg.Transfer
	)
	if data, err := rpcClient.Status(); err == nil {
		fmt.Println(data.NodeInfo.Network, data.SyncInfo.LatestBlockHeight)
	}
	denom := "BNB"

	receivers := map[string]float64{
		"tbnb1la9uz527hwv78tz7jrdmheu6fx3m3gevkqc2zm": 0.001200,
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

		if sendResult, err := rpcClient.SendToken(msgSendToken, rpc.Sync, option); err != nil {
			fmt.Printf("send token occur error, toAddr is %s, err is %s\n", v.ToAddr.String(), err.Error())
		} else {
			//if sendResult.Ok {
			fmt.Printf("send token success, toAddr is %s, txHash is %s\n", v.ToAddr.String(), sendResult.Hash)
			//} else {
			//	fmt.Printf("send token fail, toAddr is %s, txHash is %s, log is %s\n", v.ToAddr.String(),
			//		sendResult.Hash, sendResult.Log)
			//}
		}
		fmt.Println("now sleep 5 seconds")
		//time.Sleep(time.Duration(5) * time.Second)
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

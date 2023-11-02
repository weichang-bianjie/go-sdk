package e2e

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/binance-chain/go-sdk/types/tx"
	"io"
	"math/rand"
	"net/http"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tmquery "github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/types"

	"github.com/binance-chain/go-sdk/client/rpc"
	"github.com/binance-chain/go-sdk/client/transaction"
	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types/msg"
)

var (
	nodeAddr           = "tcp://data-seed-pre-2-s1.bnbchain.org:80"
	badAddr            = "tcp://127.0.0.1:80"
	testTxHash         = "073D00B80A516AA0B1F522F76C5FAED486223A4FA432AF37B51C28B63CCD5C11"
	testTxHeight       = int64(47905085)
	testAddress        = "tbnb1cegl4x48qy6mq5vg5wtryk806n2vjtyhk3sj6v"
	testDelAddr        = "tbnb12hlquylu78cjylk5zshxpdj6hf3t0tahwjt3ex"
	testTradePair      = "PPC-00A_BNB"
	testTradeSymbol    = "000-0E1"
	testTxStr          = "db01f0625dee0a63ce6dc0430a14813e4939f1567b219704ffc2ad4df58bde010879122b383133453439333946313536374232313937303446464332414434444635384244453031303837392d34341a0d5a454252412d3136445f424e422002280130c0843d38904e400112700a26eb5ae9872102139bdd95de72c22ac2a2b0f87853b1cca2e8adf9c58a4a689c75d3263013441a124015e99f7a686529c76ccc2d70b404af82ca88dfee27c363439b91ea0280571b2731c03b902193d6a5793baf64b54bcdf3f85e0d7cf657e1a1077f88143a5a65f518d2e518202b"
	mnemonic           = "test mnemonic"
	onceClient         = sync.Once{}
	testClientInstance *rpc.HTTP
)

func startBnbchaind(t *testing.T) *exec.Cmd {
	cmd := exec.Command("bnbchaind", "start", "--home", "testnoded")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Start()
	assert.NoError(t, err)
	// wait for completely start
	time.Sleep(15 * time.Second)
	return cmd
}

func defaultClient() *rpc.HTTP {
	onceClient.Do(func() {
		testClientInstance = rpc.NewRPCClient(nodeAddr, ctypes.TestNetwork)
	})
	return testClientInstance
}

func TestRPCGetProposals(t *testing.T) {
	c := defaultClient()
	statuses := []ctypes.ProposalStatus{
		ctypes.StatusDepositPeriod,
		ctypes.StatusVotingPeriod,
		ctypes.StatusPassed,
		ctypes.StatusRejected,
	}
	for _, s := range statuses {
		proposals, err := c.GetProposals(s, 100)
		assert.NoError(t, err)
		for _, p := range proposals {
			assert.Equal(t, p.GetStatus(), s)
		}
		bz, err := json.Marshal(proposals)
		fmt.Println(string(bz))
	}
}
func TestRPCGetTimelocks(t *testing.T) {
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	records, err := c.GetTimelocks(acc)
	assert.NoError(t, err)
	fmt.Println(len(records))
	for _, record := range records {
		fmt.Println(record)
	}
}

func TestRPCGetTimelock(t *testing.T) {
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	record, err := c.GetTimelock(acc, 1)
	assert.NoError(t, err)
	fmt.Println(record)

}

func TestRPCGetProposal(t *testing.T) {
	c := defaultClient()
	proposal, err := c.GetProposal(int64(1))
	assert.NoError(t, err)
	bz, err := json.Marshal(proposal)
	fmt.Println(string(bz))
}

func TestRPCStatus(t *testing.T) {
	c := defaultClient()
	status, err := c.Status()
	assert.NoError(t, err)
	bz, err := json.Marshal(status)
	fmt.Println(string(bz))
}

func TestRPCABCIInfo(t *testing.T) {
	c := defaultClient()
	info, err := c.ABCIInfo()
	assert.NoError(t, err)
	bz, err := json.Marshal(info)
	fmt.Println(string(bz))
}

func TestUnconfirmedTxs(t *testing.T) {
	c := defaultClient()
	txs, err := c.UnconfirmedTxs(10)
	assert.NoError(t, err)
	bz, err := json.Marshal(txs)
	fmt.Println(string(bz))
}

func TestNumUnconfirmedTxs(t *testing.T) {
	c := defaultClient()
	numTxs, err := c.NumUnconfirmedTxs()
	assert.NoError(t, err)
	bz, err := json.Marshal(numTxs)
	fmt.Println(string(bz))
}

func TestNetInfo(t *testing.T) {
	c := defaultClient()
	netInfo, err := c.NetInfo()
	assert.NoError(t, err)
	bz, err := json.Marshal(netInfo)
	fmt.Println(string(bz))
}

func TestDumpConsensusState(t *testing.T) {
	c := defaultClient()
	state, err := c.DumpConsensusState()
	assert.NoError(t, err)
	bz, err := json.Marshal(state)
	fmt.Println(string(bz))
}

func TestConsensusState(t *testing.T) {
	c := defaultClient()
	state, err := c.ConsensusState()
	assert.NoError(t, err)
	bz, err := json.Marshal(state)
	fmt.Println(string(bz))
}

func TestHealth(t *testing.T) {
	c := defaultClient()
	health, err := c.Health()
	assert.NoError(t, err)
	bz, err := json.Marshal(health)
	fmt.Println(string(bz))
}

func TestBlockchainInfo(t *testing.T) {
	c := defaultClient()
	blockInfos, err := c.BlockchainInfo(1, 5)
	assert.NoError(t, err)
	bz, err := json.Marshal(blockInfos)
	fmt.Println(string(bz))
}

func TestGenesis(t *testing.T) {
	c := defaultClient()
	genesis, err := c.Genesis()
	assert.NoError(t, err)
	bz, err := json.Marshal(genesis)
	fmt.Println(string(bz))
}

func TestBlock(t *testing.T) {
	c := defaultClient()
	block, err := c.Block(nil)
	assert.NoError(t, err)
	bz, err := json.Marshal(block)
	fmt.Println(string(bz))
}

func TestBlockResults(t *testing.T) {
	c := defaultClient()
	block, err := c.BlockResults(&testTxHeight)
	assert.NoError(t, err)
	bz, err := json.Marshal(block)
	fmt.Println(string(bz))
}

func TestCommit(t *testing.T) {
	c := defaultClient()
	commit, err := c.Commit(nil)
	assert.NoError(t, err)
	bz, err := json.Marshal(commit)
	fmt.Println(string(bz))
}

func TestTx(t *testing.T) {
	c := defaultClient()
	bz, err := hex.DecodeString(testTxHash)
	assert.NoError(t, err)

	tx, err := c.Tx(bz, false)
	assert.NoError(t, err)
	bz, err = json.Marshal(tx)
	fmt.Println(string(bz))
}

func TestReconnection(t *testing.T) {
	repeatNum := 10
	c := defaultClient()

	// Find error
	time.Sleep(1 * time.Second)
	for i := 0; i < repeatNum; i++ {
		_, err := c.Status()
		assert.Error(t, err)
	}

	// Reconnect and find no error
	cmd := startBnbchaind(t)

	for i := 0; i < repeatNum; i++ {
		status, err := c.Status()
		assert.NoError(t, err)
		bz, err := json.Marshal(status)
		fmt.Println(string(bz))
	}

	// kill process
	err := cmd.Process.Kill()
	assert.NoError(t, err)
	err = cmd.Process.Release()
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)

	// Find error
	for i := 0; i < repeatNum; i++ {
		_, err := c.Status()
		assert.Error(t, err)
	}

	// Restart bnbchain
	cmd = startBnbchaind(t)

	// Find no error
	for i := 0; i < repeatNum; i++ {
		status, err := c.Status()
		assert.NoError(t, err)
		bz, _ := json.Marshal(status)
		fmt.Println(string(bz))
	}

	// Stop bnbchain
	cmd.Process.Kill()
	cmd.Process.Release()
}

func TestTxSearch(t *testing.T) {
	c := defaultClient()
	tx, err := c.TxInfoSearch(fmt.Sprintf("tx.height=%d", testTxHeight), false, 1, 10)
	assert.NoError(t, err)
	bz, err := json.Marshal(tx)
	fmt.Println(string(bz))
}

func TestValidators(t *testing.T) {
	c := defaultClient()
	validators, err := c.Validators(nil)
	assert.NoError(t, err)
	bz, err := json.Marshal(validators)
	fmt.Println(string(bz))
}

func TestBadNodeAddr(t *testing.T) {
	c := rpc.NewRPCClient(badAddr, ctypes.TestNetwork)
	_, err := c.Validators(nil)
	assert.Error(t, err, "context deadline exceeded")
}

func TestSetTimeOut(t *testing.T) {
	c := rpc.NewRPCClient(badAddr, ctypes.TestNetwork)
	c.SetTimeOut(1 * time.Second)
	before := time.Now()
	_, err := c.Validators(nil)
	duration := time.Now().Sub(before).Seconds()
	assert.True(t, duration > 1)
	assert.True(t, duration < 2)
	assert.Error(t, err, "context deadline exceeded")
}

func TestSubscribeEvent(t *testing.T) {
	c := defaultClient()
	query := "tm.event = 'CompleteProposal'"
	_, err := tmquery.New(query)
	assert.NoError(t, err)
	out, err := c.Subscribe(query, 10)
	assert.NoError(t, err)
	noMoreEvent := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case o := <-out:
				bz, err := json.Marshal(o)
				assert.NoError(t, err)
				fmt.Println(string(bz))
			case <-noMoreEvent:
				fmt.Println("no more event after")
			}
		}
	}()
	time.Sleep(10 * time.Second)
	err = c.Unsubscribe(query)
	noMoreEvent <- struct{}{}
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)
}

func TestSubscribeEventTwice(t *testing.T) {
	c := defaultClient()
	query := "tm.event = 'CompleteProposal'"
	_, err := tmquery.New(query)
	assert.NoError(t, err)
	_, err = c.Subscribe(query, 10)
	assert.NoError(t, err)
	_, err = c.Subscribe(query, 10)
	assert.Error(t, err)
}

func TestReceiveWithRequestId(t *testing.T) {
	c := defaultClient()
	c.SetTimeOut(1 * time.Second)
	w := sync.WaitGroup{}
	w.Add(2000)
	testCases := []func(t *testing.T){
		TestRPCStatus,
		TestRPCABCIInfo,
		TestUnconfirmedTxs,
		TestNumUnconfirmedTxs,
		TestNetInfo,
		TestDumpConsensusState,
		TestConsensusState,
		TestHealth,
		TestBlockchainInfo,
		TestGenesis,
		TestBlock,
		//TestBlockResults,
		TestCommit,
		//TestTx,
		//TestTxSearch,
		TestValidators,
	}
	for i := 0; i < 2000; i++ {
		testFuncIndex := rand.Intn(len(testCases))
		go func() {
			testCases[testFuncIndex](t)
			w.Done()
		}()
	}
	w.Wait()
}

func TestListAllTokens(t *testing.T) {
	c := defaultClient()
	tokens, err := c.ListAllTokens(1, 10)
	assert.NoError(t, err)
	bz, err := json.Marshal(tokens)
	fmt.Println(string(bz))
}

func TestGetTokenInfo(t *testing.T) {
	c := defaultClient()
	token, err := c.GetTokenInfo("BNB")
	assert.NoError(t, err)
	bz, err := json.Marshal(token)
	fmt.Println(string(bz))
}

func TestGetAccount(t *testing.T) {
	ctypes.Network = ctypes.TestNetwork
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	account, err := c.GetAccount(acc)
	assert.NoError(t, err)
	bz, err := json.Marshal(account)
	fmt.Println(string(bz))
	fmt.Println(hex.EncodeToString(account.GetAddress().Bytes()))

}

func TestGetBalances(t *testing.T) {
	ctypes.Network = ctypes.TestNetwork
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	balances, err := c.GetBalances(acc)
	assert.NoError(t, err)
	bz, err := json.Marshal(balances)
	fmt.Println(string(bz))
}

func TestGetBalance(t *testing.T) {
	ctypes.Network = ctypes.TestNetwork
	c := defaultClient()
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	balance, err := c.GetBalance(acc, "BNB")
	assert.NoError(t, err)
	bz, err := json.Marshal(balance)
	fmt.Println(string(bz))
}

func TestGetFees(t *testing.T) {
	c := defaultClient()
	fees, err := c.GetFee()
	assert.NoError(t, err)
	bz, err := json.Marshal(fees)
	fmt.Println(string(bz))
}

func TestGetOpenOrder(t *testing.T) {
	ctypes.Network = ctypes.TestNetwork
	acc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	c := defaultClient()
	openorders, err := c.GetOpenOrders(acc, testTradePair)
	assert.NoError(t, err)
	bz, err := json.Marshal(openorders)
	assert.NoError(t, err)
	fmt.Println(string(bz))
}

func TestGetTradePair(t *testing.T) {
	c := defaultClient()
	trades, err := c.GetTradingPairs(0, 10)
	assert.NoError(t, err)
	bz, err := json.Marshal(trades)
	fmt.Println(string(bz))
}

func TestGetDepth(t *testing.T) {
	c := defaultClient()
	depth, err := c.GetDepth(testTradePair, 2)
	assert.NoError(t, err)
	bz, err := json.Marshal(depth)
	fmt.Println(string(bz))
}

func TestSendToken(t *testing.T) {
	c := defaultClient()

	testacc, err := ctypes.AccAddressFromBech32(testAddress)
	assert.NoError(t, err)
	res, err := c.SendToken([]msg.Transfer{{testacc, []ctypes.Coin{{"BNB", 10000}}}}, rpc.Sync, transaction.WithMemo("123"))
	assert.NoError(t, err)
	bz, err := json.Marshal(res)
	fmt.Println(string(bz))
}

func TestCreateOrder(t *testing.T) {
	c := defaultClient()
	ctypes.Network = ctypes.TestNetwork
	keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
	assert.NoError(t, err)
	c.SetKeyManager(keyManager)
	createOrderResult, err := c.CreateOrder(testTradeSymbol, "BNB", msg.OrderSide.BUY, 100000000, 100000000, rpc.Commit, transaction.WithSource(100), transaction.WithMemo("test memo"))

	assert.NoError(t, err)
	bz, err := json.Marshal(createOrderResult)
	fmt.Println(string(bz))
	fmt.Println(createOrderResult.Hash.String())

	type commitData struct {
		OrderId string `json:"order_id"`
	}
	var cdata commitData
	err = json.Unmarshal([]byte(createOrderResult.Data), &cdata)
	assert.NoError(t, err)
	cancleOrderResult, err := c.CancelOrder(testTradeSymbol, "BNB", cdata.OrderId, rpc.Commit)
	assert.NoError(t, err)
	bz, _ = json.Marshal(cancleOrderResult)
	fmt.Println(string(bz))
}

func TestBroadcastTxCommit(t *testing.T) {
	c := defaultClient()
	txbyte, err := hex.DecodeString(testTxStr)
	assert.NoError(t, err)
	res, err := c.BroadcastTxCommit(types.Tx(txbyte))
	assert.NoError(t, err)
	fmt.Println(res)
}

func TestGetStakeValidators(t *testing.T) {
	c := defaultClient()
	ctypes.Network = ctypes.TestNetwork
	vals, err := c.GetStakeValidators()
	assert.NoError(t, err)
	bz, err := json.Marshal(vals)
	fmt.Println(string(bz))
}

func TestGetDelegatorUnbondingDelegations(t *testing.T) {
	c := defaultClient()
	ctypes.Network = ctypes.TestNetwork
	acc, err := ctypes.AccAddressFromBech32(testDelAddr)
	assert.NoError(t, err)
	vals, err := c.GetDelegatorUnbondingDelegations(acc)
	assert.NoError(t, err)
	bz, err := json.Marshal(vals)
	fmt.Println(string(bz))
}

func TestNoRequestLeakInBadNetwork(t *testing.T) {
	c := rpc.NewRPCClient(badAddr, ctypes.TestNetwork)
	c.SetTimeOut(1 * time.Second)
	w := sync.WaitGroup{}
	w.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			c.GetFee()
			w.Done()
		}()
	}
	w.Wait()
	assert.Equal(t, c.PendingRequest(), 0)
}

func TestNoRequestLeakInGoodNetwork(t *testing.T) {
	c := defaultClient()
	c.SetTimeOut(1 * time.Second)
	w := sync.WaitGroup{}
	w.Add(3000)
	for i := 0; i < 3000; i++ {
		go func() {
			_, err := c.Block(nil)
			assert.NoError(t, err)
			//bz, err := json.Marshal(fees)
			//fmt.Println(string(bz))
			w.Done()
		}()
	}
	w.Wait()
	assert.Equal(t, c.PendingRequest(), 0)
}
func parseTxToStdTx(txBytes []byte) (tx.StdTx, error) {
	var parsedTx tx.StdTx
	err := tx.Cdc.UnmarshalBinaryLengthPrefixed(txBytes, &parsedTx)

	if err != nil {
		return parsedTx, err
	}

	return parsedTx, nil
}
func TestBroadCast(t *testing.T) {
	databytes, _ := hex.DecodeString("c402f0625dee0acb01b33f9a240a14a71cd5db2a70a5ea65178ef901cf3e479b48157e12141e599a39c9cc29c88cf701f96dea4efd1cb0a7ae1a2a696161316571766b667468747272393367347039717370703534773664746a74726e3237617237727077222a696161316c6a656d6d30797a6e7a353871787873387879616b376661736863667866356c676c347a6a782a209aaacb5bbffc253827390bd624330164474f6cbb82d4d407d0437ceb958172323080bc8daa063a090a03424e4210a08d064209313030303030424e4248e802500112700a26eb5ae987210387b265b8848792259340637060e4be296eae8e9d0e6d7c1a5a6d4e4f197b40bc1240a1766282cd4ffdaa744763dd0d1acf3b8d86d8e444a1d3a5bc5ca64dd705b9e4042016177c72b898d6e203bd6ace377559bf5ff4b081e8654b45b752f1535c4218bcf4012025")
	stdTx, err := parseTxToStdTx(databytes)
	if err != nil {
		t.Fatal(err)
	}
	retbytes, _ := json.Marshal(stdTx)
	t.Log(string(retbytes))

	var (
		ksFilePath = "/Users/user/Downloads/QA_Binance_Withdraw_Admin@123456_keystore.txt"
		ksAuth     = "Admin@123456"
		rpcClient  *rpc.HTTP
	)

	if km, err := keys.NewKeyStoreKeyManager(ksFilePath, ksAuth); err != nil {
		panic(err)
	} else {

		rpcClient = rpc.NewRPCClient(nodeAddr, ctypes.TestNetwork)
		rpcClient.SetKeyManager(km)
		if _, err := rpcClient.Status(); err != nil {
			fmt.Printf("init rpc client fail, err is %s\n", err.Error())
			panic(err)
		}
		fmt.Println("Km:", km.GetAddr().String())
	}
	data, err := rpcClient.BroadcastTxSync(databytes)
	if err != nil {
		t.Fatal(err)
	}
	retbytes, _ = json.Marshal(data)
	t.Log(string(retbytes))
	fmt.Println(string(retbytes))
	fmt.Println(data.Hash.String())
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

func TestHtlt(t *testing.T) {
	var (
		ksFilePath = "/Users/user/Downloads/QA_Binance_Withdraw_Admin@123456_keystore.txt"
		ksAuth     = "Admin@123456"
		rpcClient  *rpc.HTTP
	)
	km, err := keys.NewKeyStoreKeyManager(ksFilePath, ksAuth)
	if err != nil {
		panic(err)
	} else {

		rpcClient = rpc.NewRPCClient(nodeAddr, ctypes.TestNetwork)
		rpcClient.SetKeyManager(km)
		if _, err := rpcClient.Status(); err != nil {
			fmt.Printf("init rpc client fail, err is %s\n", err.Error())
			panic(err)
		}
		fmt.Println("Km:", km.GetAddr().String())
	}

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
	if !rpcClient.IsRunning() {
		t.Fatal("rpc port is abnormal")
	}
	rpcClient.SetKeyManager(km)

	val, err := rpcClient.SignHTLT(recipient, recipientOtherChain, senderOtherChain, randomNumberHash, timestamp, ctypes.Coins{amount}, expectedIncome, heightSpan, true, rpc.Sync)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(hex.EncodeToString(val))
	_, err = Post("http://localhost:8081/txs", fmt.Sprintf("{\"htlt_id\":%d,\"tx\":\"%s\",\"chain\":\"binance\"}", data.Data.Id, hex.EncodeToString(val)))
	if err != nil {
		t.Error(err.Error())
		return
	}
}

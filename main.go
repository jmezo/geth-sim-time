package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var client = &http.Client{Timeout: 70 * time.Second}

func main() {
	godotenv.Load()
	url := os.Getenv("WEB3_URL")
        if len(url) == 0 {
                panic("missing WEB3_URL env var")
        }
	numOfCalls := flag.Int("calls", 4, "number of calls")
	sameHolder := flag.Bool("same", true, "use same holder addr for each call")
	flag.Parse()
	if *numOfCalls > len(holders) {
		panic("calls flag must be less than " + fmt.Sprint(len(holders)))
	}
	datas := getDatas(*numOfCalls, *sameHolder)
	start := time.Now()
	callXtimes(*numOfCalls, datas, url)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

func getDatas(numOfCalls int, sameHolder bool) []string {
	// res := make([]string, 0)
	res := []string{}
	for i := 0; i < numOfCalls; i++ {
		if sameHolder {
			res = append(res, getData(0))
		} else {
			res = append(res, getData(i))
		}
	}
	return res
}

func getData(index int) string {
	var holder = holders[index]
	return ` { "to": "0xF24F35e5Ed0338175DeD0D972DaFD0e6B56E6F2B", "data": "0x7762a9d80000000000000000000000006c6d9a2ac42ad2601725234e25f30fb49968836200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000` + holder + `" } `
}

func callXtimes(numOfCalls int, datas []string, url string) {
	results := make(chan string)
	for i := 0; i < numOfCalls; i++ {
		ii := i
		go func() { results <- callDebugTrace(client, url, datas[ii]) }()
	}
	for i := 0; i < numOfCalls; i++ {
		<-results
	}
}

func callDebugTrace(client *http.Client, url string, data string) string {
	var jsonStr = []byte(`{"jsonrpc":"2.0","id":1,"method":"debug_traceCall","params":[` + data + `, "latest", {"tracer": "` + tracer + `", "timeout": "60s"}]}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	fmt.Println("resp status:", resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

var holders = []string{
	"82e8936b187d83fd6eb2b7dab5b19556e9deff1c",
	"d85782de3a7bff8d30b8f7b7ae4feb6fbf0600bd",
	"208c78c4491b3aaa18025894dca4571e8efbb0d6",
	"a0a9c60b968b9781eaae2ca83765b21f1848238f",
	"780ad0cfd2292ced8eb44ba4a0ea2f70b4551c8b",
	"23be060093db74f38b1a3daf57afdc1a23db0077",
	"aff561e8f736b78036c083a82b3ef131730effce",
	"7bfee91193d9df2ac0bfe90191d40f23c773c060",
	"8f5db667276e9d805bf5adb315374f8fa299699e",
	"10f7ce43ad3779a5313b913ceb421417c4993950",
	"8e3f5e7578b9e0f4b9ab65c85568c8f80aec95ae",
	"7c6b59f52578a1bea4f8d750c5b4bb044669b5cb",
	"449a72de50d19e5de3572744f892f75ee1855a9a",
	"96427109835d2cb6ba483a351c576b127cb28a41",
	"ef7d9be5a88af3c6bc4d87606b2747695485e50e",
	"487502f921ba3dadacf63dbf7a57a978c241b72c",
	"7608ae02b28232f564d9018783a56a98fe5038b0",
	"b06f7fb7bc17cdaf2bcdec9b9d869a37a42e05a2",
	"212694d63a75124bbb898092d1f022f46fd0b6d3",
	"45b522b0c2f7fed988f10e0eb14bd935d8872b59",
	"2fe9d8dffc83d207395db34c1393f4fc6e64785",
	"4467554da9e6e79ef5f90f0c0fcbf7d645c394cf",
	"6cfe9755269786f6681518c00bd22801f98f9e57",
	"9f8322bbc6d512f431a0a2aca2d732956c62de80",
	"2be68e381eaad342b9892961beb822270ef1fbd4",
	"3de1820d8d3b7f6c61c34dfd74f941c88cb27143",
	"c3d9d96da22e25499ea3a5667bde39430a20a74b",
	"42a1de863683f3230568900ba23f86991d012f42",
	"5b68e207687884a2ca4b3e6cc5885f626fcde69b",
}

const tracer = "{" +
	"retVal: []," +
	"afterSload: false," +
	"step: function(log,db) {" +
	"   if(log.op.toNumber() == 0x55) " +
	"     this.retVal.push(log.getPC() + \\\": SSTORE \\\" +" +
	"        log.stack.peek(0).toString(16) + \\\" <- \\\" +" +
	"        log.stack.peek(1).toString(16));" +
	"}," +
	"fault: function(log,db) {this.retVal.push(\\\"FAULT: \\\" + JSON.stringify(log))}," +
	"result: function(ctx,db) {return this.retVal}" +
	"}"

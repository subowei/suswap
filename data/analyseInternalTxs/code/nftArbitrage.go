package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"io/ioutil"
	"os"

	// "sort"
	"strconv"
	// "math/big"
	// "fmt"
)

var mapEvents = make(map[string]map[string]string)

func load(data interface{}, filename string) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		panic(err)
	}
}

//套利模式1:同区块套利，在一个区块中先买后卖来赚取差价。买卖是两笔交易
func getArbitrageData_1() {
	f1, err1 := os.Create("../data/csv/arbitrageERC721Txs_mode1.csv") //创建文件
	if err1 != nil {
		panic(err1)
	}
	defer f1.Close()
	f1.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w1 := csv.NewWriter(f1)        //创建一个新的写入文件流
	data1 := [][]string{
		{"arbitrageAddr", "valueEarn", "tokenAddress", "tokenId", "transactionBuy", "transactionSell", "blockNumBuy", "fromAddrBuy", "toAddrBuy", "fromAddrSell", "toAddrSell", "valueBuy", "valueSell"},
	}
	w1.WriteAll(data1) //写入数据
	w1.Flush()

	for _, value0 := range mapEvents {
		blockNum0, _ := strconv.ParseInt(value0["blockNum"], 10, 64)
		for _, value1 := range mapEvents {
			blockNum1, _ := strconv.ParseInt(value1["blockNum"], 10, 64)
			if blockNum1 < blockNum0 {
				continue
			}
			toAddrEvent0 := value0["eventToAddr"]
			tokenAddress0 := value0["tokenAddress"]
			tokenId0 := value0["tokenId"]
			valueNft0, _ := strconv.ParseInt(value0["value"], 10, 64)

			fromAddrEvent1 := value1["eventFromAddr"]
			tokenAddress1 := value1["tokenAddress"]
			tokenId1 := value1["tokenId"]
			valueNft1, _ := strconv.ParseInt(value1["value"], 10, 64)

			//此时是将nft代理给交易所
			// if valueNft0 == 0 {
			// 	continue
			// }

			if toAddrEvent0 == fromAddrEvent1 && tokenAddress0 == tokenAddress1 && tokenId0 == tokenId1 && blockNum0 == blockNum1 && value0["transactionHash"] != value1["transactionHash"] {
				valueEarn := valueNft1 - valueNft0
				valueNft0_string := strconv.FormatInt(valueNft0, 10)
				valueNft1_string := strconv.FormatInt(valueNft1, 10)
				valueEarn_string := strconv.FormatInt(valueEarn, 10)

				data1 := [][]string{
					{toAddrEvent0, valueEarn_string, tokenAddress0, tokenId0, value0["transactionHash"], value1["transactionHash"], value0["blockNum"], value0["eventFromAddr"], value0["eventToAddr"], value1["eventFromAddr"], value1["eventToAddr"], valueNft0_string, valueNft1_string},
				}
				w1.WriteAll(data1) //写入数据
				w1.Flush()

			}
		}
	}
}

func main() {
	load(&mapEvents, "../../output/output_event_erc721/1130To/mapEventsERC721_11439999")
	// for key, value := range mapEvents {
	// 	fmt.Println("key:",key)
	// 	fmt.Println("value:",value)
	// 	break;
	// }

	//套利模式1:同区块套利，在一个区块中先买后卖
	getArbitrageData_1()

}

package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	// "math/big"
	// "strconv"
)

var mapAddrs = make(map[string]map[string]string)

var mapTxsForArbitrage = make(map[string]map[string]string)

func store(data interface{}, filename string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, buffer.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

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

type TxInfo struct {
	address string
	info    map[string]string
}

type TxInfos []TxInfo

func getGethOriDataSortByBlockNum() {
	var txInfos TxInfos
	for key, value := range mapAddrs {
		if value["tokenIdOwnerAddrOriginal"] == value["tokenIdOwnerAddrEdited"] {
			continue
		}
		txInfos = append(txInfos, TxInfo{key, value})
	}
	// 按照blockNum从小到大排序
	sort.SliceStable(txInfos, func(i, j int) bool {
		i_blockNum, _ := strconv.Atoi(txInfos[i].info["blockNum"])
		j_blockNum, _ := strconv.Atoi(txInfos[j].info["blockNum"])
		i_position, _ := strconv.Atoi(txInfos[i].info["positionOriginal"])
		j_position, _ := strconv.Atoi(txInfos[j].info["positionOriginal"])

		if i_blockNum < j_blockNum {
			return true
		}
		if (i_blockNum == j_blockNum) && (i_position < j_position) {
			return true
		}
		return false
	})

	f, err := os.Create("../data/csv/addrMapERC721ResSorted.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)         //创建一个新的写入文件流
	data := [][]string{
		{"keyAddr", "keyAddrEventType", "blockNum", "transactionHash", "value", "fromAddrEvent", "toAddrEvent", "tokenAddress", "tokenId", "tokenIdOwnerAddrOriginal", "tokenIdOwnerAddrEdited", "positionOriginal", "positionEdited"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()
	for _, item := range txInfos {
		key := item.address
		value := item.info

		data := [][]string{
			{key[0:42], value["keyAddrEventType"], value["blockNum"], value["transactionHash"], value["value"], value["fromAddrEvent"], value["toAddrEvent"], value["tokenAddress"], value["tokenId"], value["tokenIdOwnerAddrOriginal"], value["tokenIdOwnerAddrEdited"], value["positionOriginal"], value["positionEdited"]},
		}
		w.WriteAll(data) //写入数据
		w.Flush()
	}
}

func getGethOriData() {
	f0, err0 := os.Create("../data/csv/addrMapERC721Res.csv") //创建文件
	if err0 != nil {
		panic(err0)
	}
	defer f0.Close()
	f0.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w0 := csv.NewWriter(f0)        //创建一个新的写入文件流
	data0 := [][]string{
		{"keyAddr", "keyAddrEventType", "blockNum", "transactionHash", "fromAddrEvent", "toAddrEvent", "tokenAddress", "tokenId", "tokenIdOwnerAddrOriginal", "tokenIdOwnerAddrEdited", "positionOriginal", "positionEdited"},
	}
	w0.WriteAll(data0) //写入数据
	w0.Flush()

	// 遍历原始addr map
	for key, value := range mapAddrs {
		if value["tokenIdOwnerAddrOriginal"] != value["tokenIdOwnerAddrEdited"] {
			data0 := [][]string{
				{key[0:42], value["keyAddrEventType"], value["blockNum"], value["transactionHash"], value["value"], value["fromAddrEvent"], value["toAddrEvent"], value["tokenAddress"], value["tokenId"], value["tokenIdOwnerAddrOriginal"], value["tokenIdOwnerAddrEdited"], value["positionOriginal"], value["positionEdited"]},
			}
			w0.WriteAll(data0) //写入数据
			w0.Flush()
		}

		mapTxsForArbitrage_value := map[string]string{"blockNum": value["blockNum"], "value": value["value"], "fromAddrEvent": value["fromAddrEvent"], "toAddrEvent": value["toAddrEvent"], "tokenAddress": value["tokenAddress"], "tokenId": value["tokenId"], "tokenIdOwner": value["tokenIdOwner"]}
		mapTxsForArbitrage[value["transactionHash"]] = mapTxsForArbitrage_value
	}
}

func getArbitrageData() {
	//输出1：统计套利地址················································
	//用于保存套利交易
	f1, err1 := os.Create("../data/csv/arbitrageERC721Txs.csv") //创建文件
	if err1 != nil {
		panic(err1)
	}
	defer f1.Close()
	f1.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w1 := csv.NewWriter(f1)        //创建一个新的写入文件流
	data1 := [][]string{
		{"arbitrageAddr", "valueEarn", "tokenAddress", "tokenId", "transactionBuy", "transactionSell", "blockNumBuy", "blockNumSell", "fromAddrBuy", "toAddrBuy", "fromAddrSell", "toAddrSell", "valueBuy", "valueSell"},
	}
	w1.WriteAll(data1) //写入数据
	w1.Flush()

	//遍历erc721的交易，找出存套利的行为
	for key0, value0 := range mapTxsForArbitrage {
		blockNum0, _ := strconv.ParseInt(value0["blockNum"], 10, 64)
		for key1, value1 := range mapTxsForArbitrage {
			blockNum1, _ := strconv.ParseInt(value1["blockNum"], 10, 64)
			if blockNum1 < blockNum0 {
				continue
			}
			toAddrEvent0 := value0["toAddrEvent"]
			tokenAddress0 := value0["tokenAddress"]
			tokenId0 := value0["tokenId"]
			valueNft0, _ := strconv.ParseInt(value0["value"], 10, 64)

			fromAddrEvent1 := value1["fromAddrEvent"]
			tokenAddress1 := value1["tokenAddress"]
			tokenId1 := value1["tokenId"]
			valueNft1, _ := strconv.ParseInt(value1["value"], 10, 64)

			//此时是将nft代理给交易所
			if valueNft0 == 0 {
				continue
			}

			if toAddrEvent0 == fromAddrEvent1 && tokenAddress0 == tokenAddress1 && tokenId0 == tokenId1 && valueNft0 < valueNft1 {
				valueEarn := valueNft1 - valueNft0
				valueNft0_string := strconv.FormatInt(valueNft0, 10)
				valueNft1_string := strconv.FormatInt(valueNft1, 10)
				valueEarn_string := strconv.FormatInt(valueEarn, 10)

				data1 := [][]string{
					{toAddrEvent0, valueEarn_string, tokenAddress0, tokenId0, key0, key1, value0["blockNum"], value1["blockNum"], value0["fromAddrEvent"], value0["toAddrEvent"], value1["fromAddrEvent"], value1["toAddrEvent"], valueNft0_string, valueNft1_string},
				}
				w1.WriteAll(data1) //写入数据
				w1.Flush()

			}
		}
	}
}

func main() {
	//mapTx := map[string]string{"blockNum", "transactionHash", "tokenAddress", "tokenBalanceOriginal", "tokenBalanceEdited","positionOriginal","positionEdited","addrType"}

	// load(&mapAddrs, "/home/dell/sbw/gethv7/data/output_addrInfo_erc721/1000To/mapAddrsERC721_10614999")
	load(&mapAddrs, "../../output/output_addrInfo_erc721/1130To/mapAddrsERC721_11473999")
	// load(&mapAddrs, "../../output_beifen/mapAddrsERC721_10504299")

	// getGethOriDataSortByBlockNum()
	for key, value := range mapAddrs {
		fmt.Println("!")
		if value["tokenIdOwnerAddrEdited"] != value["tokenIdOwnerAddrOriginal"] {
			fmt.Println("key", key)
			fmt.Println("value", value)
			break
		}

	}

}

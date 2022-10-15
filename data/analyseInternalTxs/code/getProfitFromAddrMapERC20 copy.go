package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	// "math/big"
	// "fmt"
	"strconv"
)

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
	transactionHash string
	info            map[string]string
}

type TxInfos []TxInfo

var mapAddrs = make(map[string]map[string]string)

var mapTxsForArbitrage = make(map[string]map[string]string)

var mapBlockProfit = make(map[string]map[string]int64)

//输出0：得到geth原始输出信息·············································
func getGethOriData() {
	f0, err0 := os.Create("../data/csv/addrMapERC20Res.csv") //创建文件
	if err0 != nil {
		panic(err0)
	}
	defer f0.Close()
	f0.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w0 := csv.NewWriter(f0)        //创建一个新的写入文件流
	data0 := [][]string{
		{"addr", "blockNum", "transactionHash", "tokenAddress", "tokenBalanceDiff", "tokenBalanceOriginal", "tokenBalanceEdited", "positionOriginal", "positionEdited", "addrType"},
	}
	w0.WriteAll(data0) //写入数据
	w0.Flush()

	for key, value := range mapAddrs {
		// fmt.Println(value)
		blockNum := value["blockNum"]
		tokenAddress := value["tokenAddress"]
		tokenBalanceOriginal, _ := strconv.ParseInt(value["tokenBalanceOriginal"], 10, 64)
		tokenBalanceEdited, _ := strconv.ParseInt(value["tokenBalanceEdited"], 10, 64)
		tokenBalanceDiff := tokenBalanceOriginal - tokenBalanceEdited
		if tokenBalanceOriginal == tokenBalanceEdited {
			continue
		}
		data0 := [][]string{
			{key, value["blockNum"], value["transactionHash"], value["tokenAddress"], strconv.FormatInt(tokenBalanceDiff, 10), value["tokenBalanceOriginal"], value["tokenBalanceEdited"], value["positionOriginal"], value["positionEdited"], value["addrType"]},
		}
		w0.WriteAll(data0) //写入数据
		w0.Flush()

		mapBlockProfit_value, ok := mapBlockProfit[blockNum]
		if ok == true {
			mapBlockProfit_value[tokenAddress] += tokenBalanceDiff
		} else {
			mapBlockProfit_value = make(map[string]int64)
			mapBlockProfit_value[tokenAddress] = tokenBalanceDiff
		}
		mapBlockProfit[blockNum] = mapBlockProfit_value

		// mapTxsForArbitrage_value := map[string]string{"blockNum": value["blockNum"], "value": value["value"], "fromAddrEvent": value["fromAddrEvent"], "toAddrEvent": value["toAddrEvent"], "tokenAddress": value["tokenAddress"], "tokenId": value["tokenId"], "tokenIdOwner": value["tokenIdOwner"]}

	}

}

//输出1：统计区块代币余额差·············································
func getTokenBalanceDiffData() {
	f1, err1 := os.Create("../data/csv/blockProfitFromAddrMapERC20Res.csv") //创建文件
	if err1 != nil {
		panic(err1)
	}
	defer f1.Close()
	f1.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w1 := csv.NewWriter(f1)        //创建一个新的写入文件流
	data1 := [][]string{
		{"blockNum", "tokenAddress", "tokenBalanceDifference"},
	}
	w1.WriteAll(data1) //写入数据
	w1.Flush()

	for key, value := range mapBlockProfit {
		for key1, value1 := range value {
			data1 := [][]string{
				{key, key1, strconv.FormatInt(value1, 10)},
			}
			w1.WriteAll(data1) //写入数据
			w1.Flush()
		}

	}
}

//输出2：统计代币套利·············································
func getArbitrageData() {
	f2, err2 := os.Create("../data/csv/arbitrageERC20Txs.csv") //创建文件
	if err2 != nil {
		panic(err2)
	}
	defer f2.Close()
	f2.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w2 := csv.NewWriter(f2)        //创建一个新的写入文件流
	data2 := [][]string{
		{"blockNum", "tokenAddress", "tokenBalanceDifference"},
	}
	w2.WriteAll(data2) //写入数据
	w2.Flush()
}

//输出3：每个地址在不同token上面的余额差总和·············································
func getGethSpecialDefiData() {
	//{addr: { "WETH":每种token对应的tokenDiff总和 }  }
	var addrDefi = make(map[string]map[string]int64)

	tokenAddr2name := map[string]string{"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2": "WETH",
		"0xdac17f958d2ee523a2206206994597c13d831ec7": "USDT",
		"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48": "USDC",
		"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599": "WBTC",
		"0x6b175474e89094c44da98b954eedeac495271d0f": "DAI",
		"0x514910771af9ca656af840dff83e8264ecf986ca": "LINK",
		"0x0d8775f648430679a709e98d2b0cb6250d2887ef": "BAT"}

	for key, value := range mapAddrs {
		key = key[0:42]
		key = strings.ToLower(key)
		// fmt.Println(value)
		tokenAddress := value["tokenAddress"]
		tokenAddress = strings.ToLower(tokenAddress)
		tokenBalanceOriginal, _ := strconv.ParseInt(value["tokenBalanceOriginal"], 10, 64)
		tokenBalanceEdited, _ := strconv.ParseInt(value["tokenBalanceEdited"], 10, 64)
		//打印特定数据
		if value["transactionHash"] == "0x9dfc5312cd1f095ac04a452afd5393b881afa8050b97792a5e6164acbe548446" {
			fmt.Println("")
			fmt.Println("key_!!!!", key)
			fmt.Println("value_!!!!", value)
			fmt.Println("")
		}

		tokenBalanceDiff := tokenBalanceOriginal - tokenBalanceEdited
		if tokenBalanceOriginal == tokenBalanceEdited {
			continue
		}
		//打印所有数据
		// if value["transactionHash"]=="0xa2d197e3c9d944f1722409d14e16c9bbad385bd6e9a421e12771abbf3de241aa"{
		// 	fmt.Println("")
		// 	fmt.Println("key_!!!!", key)
		// 	fmt.Println("value_!!!!", value)
		// 	fmt.Println("")
		// }

		addrDefi_value, addrDefi_ok := addrDefi[key]
		if addrDefi_ok == false {
			addrDefi_value = map[string]int64{"WETH": 0, "USDT": 0, "USDC": 0, "WBTC": 0, "DAI": 0, "LINK": 0, "BAT": 0}
		}
		tokenAddr2name_value, tokenAddr2name_ok := tokenAddr2name[tokenAddress]
		if tokenAddr2name_ok == true {
			addrDefi_value[tokenAddr2name_value] += tokenBalanceDiff
			addrDefi[key] = addrDefi_value

			fmt.Println("key", key)
			fmt.Println("value", value)
		}
	}

	//输出到外存
	f, err := os.Create("../data/csv/addrMapSpecialDefiERC20Res.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)         //创建一个新的写入文件流
	data := [][]string{
		{"addr", "WETH", "USDT", "USDC", "WBTC", "DAI", "LINK", "BAT"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()

	for key, value := range addrDefi {
		// tokenNameArray := [...]string{"WETH", "USDT", "USDC", "WBTC", "DAI"}
		// var writeData []string
		// writeData = append(writeData, key)
		// for _, tokenNameArray_value := range tokenNameArray {
		// 	writeData = append(writeData, strconv.FormatInt(value[tokenNameArray_value], 10))
		// }

		writeData := []string{key,
			strconv.FormatInt(value["WETH"], 10),
			strconv.FormatInt(value["USDT"], 10),
			strconv.FormatInt(value["USDC"], 10),
			strconv.FormatInt(value["WBTC"], 10),
			strconv.FormatInt(value["DAI"], 10),
			strconv.FormatInt(value["LINK"], 10),
			strconv.FormatInt(value["BAT"], 10)}

		w.Write(writeData) //写入数据
		w.Flush()
	}

}

func main() {
	load(&mapAddrs, "../../output/output_addrInfo_erc20/600To/mapAddrsERC20_6995999")
	for key, value := range mapAddrs {
		if value["balanceEdited"]!="0"{
			fmt.Println("key: ",key)
			fmt.Println("value: ",value)
			break
		}

	}
	// getGethSpecialDefiData()
}

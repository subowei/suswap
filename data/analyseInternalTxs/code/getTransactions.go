package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	// "math/big"
	// "fmt"
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

func getTransactions() {
	testBlockNum := "13499115"

	//store
	f, err := os.Create("../data/csv/txs_" + testBlockNum + ".csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)         //创建一个新的写入文件流
	data := [][]string{
		{"transaction", "blockNum", "fromAddr", "toAddr", "gasPrice", "positionOriginal", "positionEdited"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()

	var txsMap = make(map[string]map[string]string)
	path := "/mnt/4t3/geth/output/output_tx/1300To/mapTxs13499999"
	load(&txsMap, path)

	for key, value := range txsMap {
		blockNum := value["blockNum"]
		if blockNum == testBlockNum {
			// fmt.Println(key)
			// fmt.Println(value)
			writeData := []string{key, value["blockNum"], value["fromAddr"], value["toAddr"], value["gasPrice"], value["positionOriginal"], value["positionEdited"]}
			w.Write(writeData) //写入数据
			w.Flush()
		}
	}
}

func getNumOfTransactionsInBlock() {
	var blockMap = make(map[string]map[string]int)
	var txsMap = make(map[string]map[string]string)
	path := "../../output/output_tx/1190To/mapTxs12399999"
	load(&txsMap, path)

	for _, txsMap_value := range txsMap {
		blockNum := txsMap_value["blockNum"]
		positionOriginal := txsMap_value["positionOriginal"]
		positionEdited := txsMap_value["positionEdited"]
		numOfUnordered := 0
		if positionOriginal != positionEdited {
			numOfUnordered = 1
		}

		blockNum_int, _ := strconv.Atoi(blockNum)
		if blockNum_int >= 12960000 {
			continue
		}
		_, ok := blockMap[blockNum]
		if ok == false {
			blockMap_value := make(map[string]int)
			blockMap_value["num"] = 1
			blockMap_value["numOfUnordered"] = numOfUnordered
			blockMap[blockNum] = blockMap_value

		} else {
			blockMap_value := blockMap[blockNum]
			blockMap_value["num"] += 1
			blockMap_value["numOfUnordered"] += numOfUnordered
			blockMap[blockNum] = blockMap_value
		}
	}

	//找到交易数量最少的区块
	minNum := int(999999)
	var tempBlockNum string
	for blockMap_key, blockMap_value := range blockMap {
		if blockMap_value["numOfUnordered"] == 0 {
			continue
		}
		if blockMap_value["num"] < minNum && blockMap_value["num"] > 80 && blockMap_value["num"] < 100 {
			minNum = blockMap_value["num"]
			tempBlockNum = blockMap_key
		}
	}
	fmt.Println(tempBlockNum)
	fmt.Println(minNum)
}

// 12217458,12056255,12301100,12234931
func main() {
	getTransactions()
	// getNumOfTransactionsInBlock()
}

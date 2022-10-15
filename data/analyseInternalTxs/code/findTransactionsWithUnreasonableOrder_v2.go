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

func getTransactionsInSpecialBlock() {
	// 1. find designated txs,and put into struct
	for key, value := range txsMap {
		//mapTxs13299999
		blockNum := value["blockNum"]
		blockNum_int, _ := strconv.Atoi(blockNum)
		// if blockNum_int == 12399017 {
		// 	fmt.Println("12399017 exit")
		// }
		if blockNum_int < 13299899 {
			continue
		}
		// if blockNum_int > 12399217 {
		// 	break
		// }
		// _, ok := blcokNumMap[blockNum]
		// if ok == false {
		// 	continue
		// }
		txInfos_oriOrder = append(txInfos_oriOrder, TxInfo{key, value})
		txInfos_editOrder = append(txInfos_editOrder, TxInfo{key, value})
	}

	// 2. sort txs by blockNum and position
	sort.SliceStable(txInfos_oriOrder, func(i, j int) bool {
		i_blockNum, _ := strconv.Atoi(txInfos_oriOrder[i].info["blockNum"])
		j_blockNum, _ := strconv.Atoi(txInfos_oriOrder[j].info["blockNum"])
		i_position, _ := strconv.Atoi(txInfos_oriOrder[i].info["positionOriginal"])
		j_position, _ := strconv.Atoi(txInfos_oriOrder[j].info["positionOriginal"])

		if i_blockNum < j_blockNum {
			return true
		}
		if (i_blockNum == j_blockNum) && (i_position < j_position) {
			return true
		}
		return false
	})

	sort.SliceStable(txInfos_editOrder, func(i, j int) bool {
		i_blockNum, _ := strconv.Atoi(txInfos_editOrder[i].info["blockNum"])
		j_blockNum, _ := strconv.Atoi(txInfos_editOrder[j].info["blockNum"])
		i_position, _ := strconv.Atoi(txInfos_editOrder[i].info["positionEdited"])
		j_position, _ := strconv.Atoi(txInfos_editOrder[j].info["positionEdited"])

		if i_blockNum < j_blockNum {
			return true
		}
		if (i_blockNum == j_blockNum) && (i_position < j_position) {
			return true
		}
		return false
	})

}

func findReasonableTxs() {
	for i := 0; i < len(txInfos_editOrder); i++ {
		for j := i + 1; j < len(txInfos_editOrder); j++ {
			if txInfos_editOrder[i].info["blockNum"] != txInfos_editOrder[j].info["blockNum"] {
				continue
			}
			if judgePositionReasonable(i, j) == true {
				// 保存顺序合理的交易哈希
				for k := i; k <= j; k++ {
					txsReasonableMap[txInfos_editOrder[k].transactionHash] = true
				}
			}
		}
	}
}

// 根据排序前后的交易位置，来判断数组i-j的交易的交易位置是否合理
func judgePositionReasonable(start_edit int, end_edit int) bool {
	startTransactionHash := txInfos_editOrder[start_edit].transactionHash
	start_ori := 0
	for k := 0; k < len(txInfos_oriOrder); k++ {
		if txInfos_oriOrder[k].transactionHash == startTransactionHash {
			start_ori = k
			break
		}
	}
	for k := start_edit; k <= end_edit; k++ {
		if start_ori+k-start_edit >= len(txInfos_oriOrder) {
			return false
		}
		if txInfos_editOrder[k].transactionHash != txInfos_oriOrder[start_ori+k-start_edit].transactionHash {
			return false
		}
	}
	return true
}

func outputUnreasonableOrderTxs() {
	f, err := os.Create("../data/csv/txsWithUnreasonablePosition_v2.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	data := [][]string{
		{"transactionHash", "blockNum", "positionOriginal", "positionEdited"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()

	for _, item := range txInfos_editOrder {
		transactionHash := item.transactionHash
		_, ok := txsReasonableMap[transactionHash]
		if ok == true {
			continue
		}
		value := item.info
		data := [][]string{
			{item.transactionHash, value["blockNum"], value["positionOriginal"], value["positionEdited"]},
		}
		w.WriteAll(data) //写入数据
		w.Flush()
	}
}

// 0xb1f0a1a80d4e69bd39f73f613d8baf252f9f4ec8fe8ad3464937055caf9699e0
// map[blockNum: fromAccountRootEdited: fromAccountRootOriginal: fromAddr: fromBalanceEdited: fromBalanceOriginal: gasPrice: method: miner: positionEdited: positionOriginal: toAccountRootEdited: toAccountRootOriginal: toAddr: toBalanceEdited: toBalanceOriginal: value:]

type TxInfo struct {
	transactionHash string
	info            map[string]string
}

type TxInfos []TxInfo

var txInfos_oriOrder TxInfos
var txInfos_editOrder TxInfos
var txInfos_res TxInfos

var txsReasonableMap = make(map[string]bool)

var txsMap = make(map[string]map[string]string)

var blcokNumMap = make(map[string]bool)

func main() {
	path := "../../output/output_tx/1326Tov1/mapTxs13299999"
	load(&txsMap, path)
	// for key, value := range txsMap {
	// 	fmt.Println(key)
	// 	fmt.Println(value)
	// 	break
	// }
	// 将待研究的区块号数组存入map
	// 得到上述区块号中的交易
	getTransactionsInSpecialBlock()
	// 找到合理顺序的交易
	findReasonableTxs()
	// 输出不合理顺序的交易
	outputUnreasonableOrderTxs()

	fmt.Println("txsReasonable", len(txsReasonableMap))
	fmt.Println("txsAll", len(txInfos_editOrder))
}

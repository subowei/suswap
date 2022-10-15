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

func generateBlcokNumMapFromArray() {
	for _, item := range blockNumArray {
		blcokNumMap[strconv.Itoa(item)] = true
	}
}

func getTransactionsInSpecialBlock() {
	// 1. find designated txs,and put into struct
	for key, value := range txsMap {
		blockNum := value["blockNum"]
		_, ok := blcokNumMap[blockNum]
		if ok == false {
			continue
		}
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
	f, err := os.Create("../data/csv/txsWithUnreasonablePosition.csv") //创建文件
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
var blockNumArray = [...]int{
	12972381,
	12972361,
	12976306,
	12420617,
	12446810,
	12427229,
	12276186,
	11673416,
	12600970,
	12381509,
	12426553,
	12192851,
	12192987,
	12446739,
	12179275,
	12160830,
	12192997,
	12108265,
	12743367,
	12192902,
	12064244,
	12282609,
	12210983,
	12977813,
	12977039,
	12648002,
	12972415,
	12972979,
	11914907,
	12105446,
	12355619,
	12977980,
	12602068,
	12973360,
	12517961,
	12308779,
	12074001,
	12985411,
	12968310,
	12492031,
	12895521,
	12829087,
	12448865,
	11907351,
	12984715,
	12985207,
	12988625,
	12510755,
	12767316,
	12981130,
	12818818,
	12513370,
	12284999,
	12977980,
	12511083,
	12977995,
	12359127,
	12543809,
	12510739,
	12465538,
	12953787,
	12929922,
	12765723,
	12898203,
	12900717,
	12884356,
	12977021,
	12977535,
	12989835,
	12225821,
	12987003,
	12988378,
	12984901,
	12988485,
	12984959,
	12986777,
	12984074,
	12982656,
	12985395,
	12984626,
	12929121,
	12507046,
	11907351,
	12298713,
	12465836,
	12684412,
	12447516,
	12468095,
	12983444,
	11907351,
	12566270,
	12507046,
	12985419,
	12321951,
	12497822,
	12322928,
	12721075,
	12501354,
	12986549,
	12969762,
	12789356,
	12312761,
	12473557,
	12100241,
	12674701,
	11673416,
	12924267,
	12608033,
	12660185,
	12749122,
	12977309,
	12973702,
	12980127,
	12815381,
	12349664,
	12200820,
	12970370,
	12459924,
	12967740,
	12669850}

func main() {
	path := "../../output/output_tx/1280To/mapTxs12989999"
	load(&txsMap, path)
	// for key, value := range txsMap {
	// 	fmt.Println(key)
	// 	fmt.Println(value)
	// 	break
	// }
	// 将待研究的区块号数组存入map
	generateBlcokNumMapFromArray()
	// 得到上述区块号中的交易
	getTransactionsInSpecialBlock()
	// 找到合理顺序的交易
	findReasonableTxs()
	// 输出不合理顺序的交易
	outputUnreasonableOrderTxs()

	fmt.Println("txsReasonable", len(txsReasonableMap))
	fmt.Println("txsAll", len(txInfos_editOrder))
}

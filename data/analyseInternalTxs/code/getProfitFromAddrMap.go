package main

import (

	// "fmt"
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"sort"
	"strings"
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

func getBlockData() {
	// 遍历原始addr map
	for key, addr_value := range mapAddrs {
		addr_key := strings.ToLower(key[0:42])
		fromAddr := strings.ToLower(addr_value["fromAddr"])
		toAddr := strings.ToLower(addr_value["toAddr"])

		blockNum := addr_value["blockNum"]
		balanceDifference_int, _ := new(big.Int).SetString(addr_value["balanceDifference"], 10)
		serviceChargeDiff_int, _ := new(big.Int).SetString(addr_value["serviceChargeDiff"], 10)

		mapBlockProfit_value, ok := mapBlockProfit[blockNum]
		if ok == true {
			mapBlockProfit_value["balanceDifference"].Add(mapBlockProfit_value["balanceDifference"], balanceDifference_int)
		} else {
			mapBlockProfit_value = make(map[string]*big.Int)
			mapBlockProfit_value["balanceDifference"] = balanceDifference_int
			mapBlockProfit_value["serviceChargeDiff"] = serviceChargeDiff_int
			mapBlockProfit_value["fromBalanceDifference"] = big.NewInt(0)
			mapBlockProfit_value["toBalanceDifference"] = big.NewInt(0)
			mapBlockProfit_value["fromAndToBalanceDifference"] = big.NewInt(0)
		}
		if addr_key == fromAddr {
			mapBlockProfit_value["fromBalanceDifference"].Add(mapBlockProfit_value["fromBalanceDifference"], balanceDifference_int)
			mapBlockProfit_value["fromAndToBalanceDifference"].Add(mapBlockProfit_value["fromAndToBalanceDifference"], balanceDifference_int)
		} else if addr_key == toAddr {
			mapBlockProfit_value["toBalanceDifference"].Add(mapBlockProfit_value["toBalanceDifference"], balanceDifference_int)
			mapBlockProfit_value["fromAndToBalanceDifference"].Add(mapBlockProfit_value["fromAndToBalanceDifference"], balanceDifference_int)
		}
		mapBlockProfit[blockNum] = mapBlockProfit_value

		//test
		// if blockNum == "12382959" && balanceDifference_int.Cmp(big.NewInt(0)) != 0 && (addr_key == fromAddr || addr_key == toAddr) {
		// 	fmt.Println("addr_key", addr_key)
		// 	fmt.Println("addr_value", addr_value)
		// }
		// blockNum_int, _ := strconv.ParseInt(blockNum, 10, 64)
		// if blockNum_int > 12382959 {
		// 	return
		// }

	}
}

func storeBlockDataAfterSorting() {
	var blockInfos BlockInfos
	for blockNum, mapBlockProfit_value := range mapBlockProfit {
		isCorrect := "false"

		balanceDifference_int, _ := mapBlockProfit_value["balanceDifference"]
		serviceChargeDiff_int, _ := mapBlockProfit_value["serviceChargeDiff"]

		serviceAndBalanceDiff_int := new(big.Int).Add(balanceDifference_int, serviceChargeDiff_int)

		if balanceDifference_int.Cmp(big.NewInt(0)) == 0 || serviceAndBalanceDiff_int.Cmp(big.NewInt(0)) == 0 {
			isCorrect = "true"
		}
		if isCorrect == "true" {
			blockInfos = append(blockInfos, BlockInfo{blockNum, mapBlockProfit_value})
		}
	}

	// sort
	sort.SliceStable(blockInfos, func(i, j int) bool {
		i_balanceDiff, _ := blockInfos[i].info["fromAndToBalanceDifference"]
		j_balanceDiff, _ := blockInfos[j].info["fromAndToBalanceDifference"]

		if i_balanceDiff.Cmp(j_balanceDiff) == -1 {
			return false
		}
		return true
	})

	// store
	f, err := os.Create("../data/csv/blockProfitFromAddrMap.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流

	data := [][]string{
		{"blockNum", "fromAndToBalanceDifference", "fromBalanceDifference", "toBalanceDifference"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()
	for _, item := range blockInfos {
		key := item.blockNum
		value := item.info
		data := [][]string{
			{key, value["fromAndToBalanceDifference"].String(), value["fromBalanceDifference"].String(), value["toBalanceDifference"].String()},
		}
		w.WriteAll(data) //写入数据
		w.Flush()
	}
}

func storeBlock() {
	// 遍历block map
	f, err := os.Create("../data/csv/blockProfitFromAddrMap.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	// data := [][]string{
	// 	{"blockNum", "fromBalanceDifference", "toBalanceDifference", "fromAndToBalanceDifference", "serviceAndBalanceDiff", "balanceDifference", "serviceChargeDiff", "isCorrect"},
	// }
	data := [][]string{
		{"blockNum", "fromBalanceDifference", "toBalanceDifference", "fromAndToBalanceDifference", "isCorrect"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()
	for blockNum, mapBlockProfit_value := range mapBlockProfit {
		isCorrect := "false"

		balanceDifference_int, _ := mapBlockProfit_value["balanceDifference"]
		serviceChargeDiff_int, _ := mapBlockProfit_value["serviceChargeDiff"]
		fromBalanceDifference_int, _ := mapBlockProfit_value["fromBalanceDifference"]
		toBalanceDifference_int, _ := mapBlockProfit_value["toBalanceDifference"]
		fromAndToBalanceDifference_int, _ := mapBlockProfit_value["fromAndToBalanceDifference"]

		serviceAndBalanceDiff_int := new(big.Int).Add(balanceDifference_int, serviceChargeDiff_int)

		if balanceDifference_int.Cmp(big.NewInt(0)) == 0 || serviceAndBalanceDiff_int.Cmp(big.NewInt(0)) == 0 {
			isCorrect = "true"
		}
		if isCorrect == "true" {
			data := [][]string{
				{blockNum, fromBalanceDifference_int.String(), toBalanceDifference_int.String(), fromAndToBalanceDifference_int.String(), isCorrect},
			}
			w.WriteAll(data) //写入数据
			w.Flush()
		}
	}
}

func getAddrData() {
	// 遍历原始addr map
	for key, addr_value := range mapAddrs {
		addr_key := strings.ToLower(key[0:42])
		fromAddr := strings.ToLower(addr_value["fromAddr"])
		toAddr := strings.ToLower(addr_value["toAddr"])
		blockNum_bigInt, _ := new(big.Int).SetString(addr_value["blockNum"], 10)
		balanceDifference_bigInt, _ := new(big.Int).SetString(addr_value["balanceDifference"], 10)
		if balanceDifference_bigInt.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		mapAddrBalanceDiff_value, ok := mapAddrBalanceDiff[addr_key]
		if ok == false {
			mapAddrBalanceDiff_value = make(map[string]*big.Int)
			mapAddrBalanceDiff_value["maxBlockNum"] = blockNum_bigInt
			mapAddrBalanceDiff_value["minBlockNum"] = blockNum_bigInt
			mapAddrBalanceDiff_value["balanceDifference"] = big.NewInt(0)
			mapAddrBalanceDiff_value["isFrom"] = big.NewInt(0)
			mapAddrBalanceDiff_value["isTo"] = big.NewInt(0)
			mapAddrBalanceDiff_value["num"] = big.NewInt(0)
		}
		// 更新blockNum
		maxBlockNum_bigInt, _ := mapAddrBalanceDiff_value["maxBlockNum"]
		// fmt.Println(maxBlockNum_bigInt)
		// fmt.Println(blockNum_bigInt)
		if blockNum_bigInt == nil {
			continue
		}
		if maxBlockNum_bigInt.Cmp(blockNum_bigInt) == -1 {
			mapAddrBalanceDiff_value["maxBlockNum"] = blockNum_bigInt
		}
		// 更新balanceDiff
		mapAddrBalanceDiff_value["balanceDifference"].Add(mapAddrBalanceDiff_value["balanceDifference"], balanceDifference_bigInt)
		// 更新isfrom
		if addr_key == fromAddr {
			mapAddrBalanceDiff_value["isFrom"] = big.NewInt(1)
		}
		// 更新isto
		if addr_key == toAddr {
			mapAddrBalanceDiff_value["isTo"] = big.NewInt(1)
		}
		// 更新num
		mapAddrBalanceDiff_value["num"].Add(mapAddrBalanceDiff_value["num"], big.NewInt(1))

		mapAddrBalanceDiff[addr_key] = mapAddrBalanceDiff_value
	}
}
func storeAddrDataAfterSorting() {
	var addrInfos AddrInfos
	for address, mapAddrBalanceDiff_value := range mapAddrBalanceDiff {
		addrInfos = append(addrInfos, AddrInfo{address, mapAddrBalanceDiff_value})
	}
	// sort
	sort.SliceStable(addrInfos, func(i, j int) bool {
		i_balanceDiff, _ := addrInfos[i].info["balanceDifference"]
		j_balanceDiff, _ := addrInfos[j].info["balanceDifference"]
		// i_num, _ := addrInfos[i].info["num"]
		// j_num, _ := addrInfos[j].info["num"]

		if i_balanceDiff.Cmp(j_balanceDiff) == -1 {
			return false
		}
		// if i_num.Cmp(j_num) == -1 {
		// 	return false
		// }
		return true
	})

	f, err := os.Create("../data/csv/addrBalanceDiff_sortedByDiff.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	data := [][]string{
		{"address", "balanceDifference", "num", "maxBlockNum", "minBlockNum", "isFrom", "isTo"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()

	for _, item := range addrInfos {
		addr_key := item.address
		addr_value := item.info
		data := [][]string{
			{addr_key, addr_value["balanceDifference"].String(), addr_value["num"].String(), addr_value["maxBlockNum"].String(), addr_value["minBlockNum"].String(), addr_value["isFrom"].String(), addr_value["isTo"].String()},
		}
		w.WriteAll(data) //写入数据
		w.Flush()
	}
}

func storeAddrData() {
	f, err := os.Create("../data/csv/addrBalanceDiff.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	data := [][]string{
		{"address", "balanceDifference", "maxBlockNum", "minBlockNum", "isFrom", "isTo"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()

	for addr_key, addr_value := range mapAddrBalanceDiff {
		data := [][]string{
			{addr_key, addr_value["balanceDifference"].String(), addr_value["maxBlockNum"].String(), addr_value["minBlockNum"].String(), addr_value["isFrom"].String(), addr_value["isTo"].String()},
		}
		w.WriteAll(data) //写入数据
		w.Flush()
	}

}

func getTransactionMapWithBalanceDifference() {
	for key, addr_value := range mapAddrs {
		transactionHash := strings.ToLower(addr_value["transactionHash"])
		balanceDifference_bigInt, _ := new(big.Int).SetString(addr_value["balanceDifference"], 10)
		if transactionHash == "0x08406a6c1390043d4fc49d3d4bd3a27ee39edc06e67db76a26b82c6f994807c5" {
			fmt.Println("key", key)
			fmt.Println("addr_value", addr_value)
		}
		if balanceDifference_bigInt.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		_, ok := txsMap[transactionHash]
		if ok == false {
			txsMap[transactionHash] = 1
		}
	}
}

func storeTransactionMapData() {
	f, err := os.Create("../data/csv/txs.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	data := [][]string{
		{"transaction"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()

	for key, _ := range txsMap {
		data := [][]string{
			{key},
		}
		w.WriteAll(data) //写入数据
		w.Flush()
	}

}

type BlockInfo struct {
	blockNum string
	info     map[string]*big.Int
}

type AddrInfo struct {
	address string
	info    map[string]*big.Int
}

type BlockInfos []BlockInfo
type AddrInfos []AddrInfo

var txsMap = make(map[string]int)
var mapAddrs = make(map[string]map[string]string)
var mapBlockProfit = make(map[string]map[string]*big.Int)
var mapAddrBalanceDiff = make(map[string]map[string]*big.Int)

func main() {
	// var paths []string = []string{"../../output/output_addrInfo/1100To/mapAddrs11409999", "../../output/output_addrInfo/1130To/mapAddrs11981999", "../../output/output_addrInfo/1190To/mapAddrs12399999", "../../output/output_addrInfo/1200To/mapAddrs12959999", "../../output/output_addrInfo/1280To/mapAddrs12989999"}
	// for i, path := range paths {
	// 	load(&mapAddrs, path)
	// 	getBlockData()
	// 	mapAddrs = make(map[string]map[string]string)
	// 	fmt.Println(i)
	// }
	// storeBlockDataAfterSorting()

	// var paths []string = []string{"../../output/output_addrInfo/1100To/mapAddrs11409999", "../../output/output_addrInfo/1130To/mapAddrs11981999", "../../output/output_addrInfo/1190To/mapAddrs12399999", "../../output/output_addrInfo/1200To/mapAddrs12959999", "../../output/output_addrInfo/1280To/mapAddrs12989999"}
	// for i, path := range paths {
	// 	load(&mapAddrs, path)
	// 	getAddrData()
	// 	mapAddrs = make(map[string]map[string]string)
	// 	fmt.Println(i)
	// }
	// storeAddrDataAfterSorting()

	// get txsMap
	var paths []string = []string{"../../output/output_addrInfo/1304ToV1/mapAddrs13079999"}
	for i, path := range paths {
		load(&mapAddrs, path)
		getTransactionMapWithBalanceDifference()
		fmt.Println(i)
	}
	storeTransactionMapData()

}

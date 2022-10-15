package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"sort"
	"strings"

	// "math/big"

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

func mergeMap_0(map0 map[string]map[string]string, map1 map[string]map[string]string) {
	for k, v := range map1 {
		map0[k] = v
	}
}

func mergeMap_1(map0 map[string]string, map1 map[string]string) {
	for k, v := range map1 {
		map0[k] = v
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

//··················································································

type AddrInfo struct {
	address              string
	blockNum             string
	transactionHash      string
	positionOriginal     string
	positionEdited       string
	gasUsedOriginal      string
	gasUsedEdited        string
	totalFeeDiff         string
	addrType             string
	effectiveFeeOriginal string

	info map[string]string
}

type AddrInfos []AddrInfo

func uniqueArr(arr []string) []string {
	newArr := make([]string, 0)
	tempArr := make(map[string]bool, len(newArr))
	for _, v := range arr {
		if tempArr[v] == false {
			tempArr[v] = true
			newArr = append(newArr, v)
		}
	}
	return newArr
}

func sumArr(arr []string) string {
	sum := big.NewInt(0)
	for _, v := range arr {
		value_v_big, _ := new(big.Int).SetString(v, 10)
		sum = new(big.Int).Add(sum, value_v_big)
	}
	return sum.String()
}

// 多个代币同时增加或者减少的情况(并且要2个以及以上),但是不局限于指定代币
func getAddrBalanceDiff_specialTokens_afterSorting_all() {
	var addrBalanceDiffMap = make(map[string]map[string]string)
	// fileDirList := [...]string{"/mnt/4t3/geth/output/output_addrInfo_erc20/1300To/", "/mnt/4t3/geth/output/output_addrInfo_erc20/1300To_v1/"}
	// for _, fileDir := range fileDirList {
	// 	fileInfoList, _ := ioutil.ReadDir(fileDir)
	// 	for i := range fileInfoList {
	// 		filePath := fileDir + fileInfoList[i].Name()
	// 		load(&mapAddrs, filePath)
	// 		for key, value := range mapAddrs {
	// 			// fmt.Println("value", value)
	// 			// return
	// 			if value["tokenBalanceEdited"] == "none" {
	// 				value["tokenBalanceEdited"] = "0"
	// 			}
	// 			if value["tokenBalanceOriginal"] == "none" {
	// 				value["tokenBalanceOriginal"] = "0"
	// 			}
	// 			if value["balanceOriginal"] == "none" {
	// 				value["balanceOriginal"] = "0"
	// 			}
	// 			if value["balanceEdited"] == "none" {
	// 				value["balanceEdited"] = "0"
	// 			}

	// 			addr := strings.ToLower(key[0:42])
	// 			blockNum := value["blockNum"]
	// 			// if blockNum == "1326210069" {
	// 			// 	fmt.Println("1326210069 exits")
	// 			// }
	// 			// continue
	// 			transactionHash := value["transactionHash"]
	// 			fromAddr := strings.ToLower(value["fromAddr"])
	// 			toAddr := strings.ToLower(value["toAddr"])
	// 			addrType := strings.ToLower(value["addrType"])
	// 			if addr == fromAddr {
	// 				addrType = "eoaFrom"
	// 			}
	// 			if addr == toAddr {
	// 				addrType = "eoaTo"
	// 			}

	// 			//单独处理miner
	// 			if transactionHash == "miner" {
	// 				addrKey := addr + "_" + blockNum
	// 				addrBalanceDiffMap_value, ok := addrBalanceDiffMap[addrKey]
	// 				balanceEdited_bigInt, _ := new(big.Int).SetString(value["balanceEdited"], 10)
	// 				balanceOriginal_bigInt, _ := new(big.Int).SetString(value["balanceOriginal"], 10)
	// 				balanceDiff_bigInt := new(big.Int).Sub(balanceOriginal_bigInt, balanceEdited_bigInt)

	// 				if ok == false {
	// 					addrBalanceDiffMap_value = make(map[string]string)
	// 					addrBalanceDiffMap_value["transactionHash"] = transactionHash
	// 					addrBalanceDiffMap_value["blockNum"] = blockNum
	// 					addrBalanceDiffMap_value["address"] = addr
	// 					addrBalanceDiffMap_value["addrType"] = "miner"
	// 					addrBalanceDiffMap_value["positionEdited"] = "-1"
	// 					addrBalanceDiffMap_value["positionOriginal"] = "-1"
	// 					addrBalanceDiffMap_value["gasUsedOriginal"] = "0"
	// 					addrBalanceDiffMap_value["gasUsedEdited"] = "0"
	// 					addrBalanceDiffMap_value["totalFeeDiff"] = "0"
	// 					addrBalanceDiffMap_value["effectiveFeeOriginal"] = "0"
	// 				}
	// 				if balanceDiff_bigInt.Cmp(big.NewInt(0)) != 0 {
	// 					addrBalanceDiffMap_value["ether"] = balanceDiff_bigInt.String()
	// 				}
	// 				addrBalanceDiffMap[addrKey] = addrBalanceDiffMap_value
	// 				continue
	// 			}

	// 			tokenAddress := strings.ToLower(value["tokenAddress"])
	// 			balanceEdited_bigInt, _ := new(big.Int).SetString(value["balanceEdited"], 10)
	// 			balanceOriginal_bigInt, _ := new(big.Int).SetString(value["balanceOriginal"], 10)
	// 			tokenbalanceOriginal_bigInt, _ := new(big.Int).SetString(value["tokenBalanceOriginal"], 10)
	// 			tokenbalanceEdited_bigInt, _ := new(big.Int).SetString(value["tokenBalanceEdited"], 10)
	// 			balanceDiff_bigInt := new(big.Int).Sub(balanceOriginal_bigInt, balanceEdited_bigInt)
	// 			tokenbalanceDiff_bigInt := new(big.Int).Sub(tokenbalanceOriginal_bigInt, tokenbalanceEdited_bigInt)
	// 			baseFee_bigInt, _ := new(big.Int).SetString(value["baseFee"], 10)
	// 			effectiveGasTipValue_bigInt, _ := new(big.Int).SetString(value["effectiveGasTipValue"], 10)
	// 			//计算手续费差
	// 			gasUsedOriginal_bigInt, _ := new(big.Int).SetString(value["gasUsedOriginal"], 10)
	// 			gasUsedEdited_bigInt, _ := new(big.Int).SetString(value["gasUsedEdited"], 10)

	// 			totalFeeOriginal_bigInt := new(big.Int).Mul(gasUsedOriginal_bigInt, baseFee_bigInt)
	// 			totalFeeEdited_bigInt := new(big.Int).Mul(gasUsedEdited_bigInt, baseFee_bigInt)
	// 			totalFeeDiff_bigInt := new(big.Int).Sub(totalFeeOriginal_bigInt, totalFeeEdited_bigInt)

	// 			effectiveFeeOriginal_bigInt := new(big.Int).Mul(gasUsedOriginal_bigInt, effectiveGasTipValue_bigInt)

	// 			addrKey := addr + "_" + transactionHash
	// 			addrBalanceDiffMap_value, ok := addrBalanceDiffMap[addrKey]
	// 			if ok == false {
	// 				addrBalanceDiffMap_value = make(map[string]string)
	// 				addrBalanceDiffMap_value["transactionHash"] = transactionHash
	// 				addrBalanceDiffMap_value["blockNum"] = blockNum
	// 				addrBalanceDiffMap_value["address"] = addr
	// 				addrBalanceDiffMap_value["addrType"] = addrType
	// 				addrBalanceDiffMap_value["positionEdited"] = value["positionEdited"]
	// 				addrBalanceDiffMap_value["positionOriginal"] = value["positionOriginal"]
	// 				addrBalanceDiffMap_value["gasUsedOriginal"] = value["gasUsedOriginal"]
	// 				addrBalanceDiffMap_value["gasUsedEdited"] = value["gasUsedEdited"]
	// 				addrBalanceDiffMap_value["totalFeeDiff"] = totalFeeDiff_bigInt.String()
	// 				addrBalanceDiffMap_value["effectiveFeeOriginal"] = effectiveFeeOriginal_bigInt.String()
	// 			}
	// 			if balanceDiff_bigInt.Cmp(big.NewInt(0)) != 0 {
	// 				if addrType != "eoaFrom" {
	// 					addrBalanceDiffMap_value["ether"] = balanceDiff_bigInt.String()
	// 				} else if addrType == "eoaFrom" {
	// 					addrBalanceDiffMap_value["ether"] = new(big.Int).Add(balanceDiff_bigInt, totalFeeDiff_bigInt).String()

	// 					_, ok = addrBalanceDiffMap_value["ether"]
	// 					if ok == false {
	// 						tempStr := balanceDiff_bigInt.String() + ";" + totalFeeDiff_bigInt.String()
	// 						addrBalanceDiffMap_value["ether"] = tempStr
	// 					} else {
	// 						tempStr := ";" + totalFeeDiff_bigInt.String()
	// 						addrBalanceDiffMap_value["ether"] += tempStr
	// 					}

	// 				}
	// 			}
	// 			if tokenbalanceDiff_bigInt.Cmp(big.NewInt(0)) != 0 {
	// 				addrBalanceDiffMap_value[tokenAddress] = tokenbalanceDiff_bigInt.String()
	// 			}
	// 			addrBalanceDiffMap[addrKey] = addrBalanceDiffMap_value
	// 		}
	// 		mapAddrs = make(map[string]map[string]string)

	// 		fmt.Println(i)
	// 	}
	// }
	// //整合ether这一项
	// for _, value := range addrBalanceDiffMap {
	// 	ether := value["ether"]
	// 	arr := strings.Split(ether, ";")
	// 	if len(arr) == 1 {
	// 		continue
	// 	}
	// 	//对arr中的值去重并且相加
	// 	uniqueArr := uniqueArr(arr)
	// 	sum := sumArr(uniqueArr)
	// 	value["ether"] = sum
	// }
	// store(addrBalanceDiffMap, "../data/dict/erc20/addrBalanceDiffMap.map")

	// return

	load(&addrBalanceDiffMap, "../data/dict/erc20/addrBalanceDiffMap.map")
	fmt.Println("load success")
	//存入结构体······················
	var addrInfos AddrInfos
	txSignMap := make(map[string]string)
	for _, value := range addrBalanceDiffMap {
		address := value["address"]
		blockNum := value["blockNum"]
		transactionHash := value["transactionHash"]
		positionOriginal := value["positionOriginal"]
		positionEdited := value["positionEdited"]
		gasUsedOriginal := value["gasUsedOriginal"]
		gasUsedEdited := value["gasUsedEdited"]
		totalFeeDiff := value["totalFeeDiff"]
		addrType := value["addrType"]
		effectiveFeeOriginal := value["effectiveFeeOriginal"]

		// numOfZero := 0
		var tempInfoMap = make(map[string]string)
		numOfIncrease := 0
		numOfDecrease := 0

		// 第一遍执行：判断哪些交易包含同向变化的地址，记录到txSignMap
		for value_k, value_v := range value {
			if value_k == "effectiveFeeOriginal" || value_k == "addrType" || value_k == "totalFeeDiff" || value_k == "address" || value_k == "blockNum" || value_k == "transactionHash" || value_k == "positionOriginal" || value_k == "positionEdited" || value_k == "gasUsedOriginal" || value_k == "gasUsedEdited" {
				continue
			}
			value_v_big, _ := new(big.Int).SetString(value_v, 10)

			if value_v_big.Cmp(big.NewInt(0)) == 1 {
				numOfIncrease += 1
			}
			if value_v_big.Cmp(big.NewInt(0)) == -1 {
				numOfDecrease += 1
			}

			tempInfoMap[value_k] = value_v
		}
		if (numOfDecrease != 0 && numOfIncrease != 0) || addrType == "miner" {
			continue
		}
		txSignMap[transactionHash] = "1"

		// 第二遍执行：保存结果
		for value_k, value_v := range value {
			if value_k == "effectiveFeeOriginal" || value_k == "addrType" || value_k == "totalFeeDiff" || value_k == "address" || value_k == "blockNum" || value_k == "transactionHash" || value_k == "positionOriginal" || value_k == "positionEdited" || value_k == "gasUsedOriginal" || value_k == "gasUsedEdited" {
				continue
			}

			value_v_big, _ := new(big.Int).SetString(value_v, 10)

			if value_v_big.Cmp(big.NewInt(0)) == 1 {
				numOfIncrease += 1
			}
			if value_v_big.Cmp(big.NewInt(0)) == -1 {
				numOfDecrease += 1
			}

			tempInfoMap[value_k] = value_v
		}
		if (numOfDecrease != 0 && numOfIncrease != 0) || addrType == "miner" {
			_, ok := txSignMap[transactionHash]
			if ok == false {
				continue
			}
		}
		addrInfos = append(addrInfos, AddrInfo{address, blockNum, transactionHash, positionOriginal, positionEdited, gasUsedOriginal, gasUsedEdited, totalFeeDiff, addrType, effectiveFeeOriginal, tempInfoMap})
	}

	//排序结构体······················
	sort.SliceStable(addrInfos, func(i, j int) bool {

		i_blockNum, _ := new(big.Int).SetString(addrInfos[i].blockNum, 10)
		j_blockNum, _ := new(big.Int).SetString(addrInfos[j].blockNum, 10)
		i_positionOriginal, _ := new(big.Int).SetString(addrInfos[i].positionOriginal, 10)
		j_positionOriginal, _ := new(big.Int).SetString(addrInfos[j].positionOriginal, 10)
		if i_blockNum.Cmp(j_blockNum) == -1 {
			return false
		}
		if i_blockNum.Cmp(j_blockNum) == 0 && i_positionOriginal.Cmp(j_positionOriginal) == 1 {
			return false
		}
		return true
	})

	//存入外存······················
	f, err := os.Create("../data/csv/erc20/addrBalanceDiff_allTokens_syndromicChange_NEW.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)         //创建一个新的写入文件流
	data := [][]string{
		{"address", "addrType", "blockNum", "transactionHash", "positionOriginal", "positionEdited", "totalFeeDiff", "effectiveFeeOriginal", "gasUsedOriginal"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()

	for _, item := range addrInfos {
		value := item.info

		data := []string{item.address, item.addrType, item.blockNum, item.transactionHash, item.positionOriginal, item.positionEdited, item.totalFeeDiff, item.effectiveFeeOriginal, item.gasUsedOriginal}
		for value_k, value_v := range value {
			data = append(data, value_k)
			data = append(data, value_v)
		}
		w.Write(data) //写入数据
		w.Flush()

	}
}

func getBlockDiff() {
	blockDiffMap := make(map[string]map[string]*big.Int)
	fileDir := "/mnt/4t3/geth/output/output_addrInfo_erc20/1300To/"
	fileInfoList, err := ioutil.ReadDir(fileDir)

	for i := range fileInfoList {
		filePath := fileDir + fileInfoList[i].Name()
		load(&mapAddrs, filePath)
		// 遍历原始addr map
		for key, value := range mapAddrs {
			addr_key := strings.ToLower(key[0:42])
			blockNum := value["blockNum"]
			transactionHash := value["transactionHash"]

			// if blockNum=="13315330"{
			// 	fmt.Println("key",key)
			// 	fmt.Println("value",value)
			// }

			// totalFeeDiff:=big.NewInt(0)
			baseFeeDiff_bigInt := big.NewInt(0)
			if transactionHash != "miner" {
				//1. 计算交易手续费差
				gasUsedOriginal_bigInt, _ := new(big.Int).SetString(value["gasUsedOriginal"], 10)
				gasUsedEdited_bigInt, _ := new(big.Int).SetString(value["gasUsedEdited"], 10)
				// gasPrice_bigInt,_:=new(big.Int).SetString(value["gasPrice"], 10)
				// effectiveGasTipValue_bigInt,_:=new(big.Int).SetString(value["effectiveGasTipValue"], 10)
				baseFee_bigInt, _ := new(big.Int).SetString(value["baseFee"], 10)
				// validGasPrice_bigInt:=new(big.Int).Add(effectiveGasTipValue_bigInt,baseFee_bigInt)
				// if validGasPrice_bigInt.Cmp(gasPrice_bigInt)==1{
				// 	validGasPrice_bigInt=gasPrice_bigInt
				// }
				// totalFeeOriginal_bigInt:=new(big.Int).Mul(gasUsedOriginal_bigInt,validGasPrice_bigInt)
				// totalFeeEdited_bigInt:=new(big.Int).Mul(gasUsedEdited_bigInt,validGasPrice_bigInt)
				// totalFeeDiff=new(big.Int).Sub(totalFeeOriginal_bigInt,totalFeeEdited_bigInt)

				gasUsedDiff_bigInt := new(big.Int).Sub(gasUsedOriginal_bigInt, gasUsedEdited_bigInt)
				baseFeeDiff_bigInt = new(big.Int).Mul(baseFee_bigInt, gasUsedDiff_bigInt)

				// fmt.Println("transactionHash: ",transactionHash)
				// fmt.Println("totalFeeOriginal_bigInt: ",totalFeeOriginal_bigInt)
				// return
			}
			//2. 计算地址余额差
			if value["balanceEdited"] == "none" {
				value["balanceEdited"] = "0"
			}
			if value["balanceOriginal"] == "none" {
				value["balanceOriginal"] = "0"
			}
			balanceEdited_bigInt, _ := new(big.Int).SetString(value["balanceEdited"], 10)
			balanceOriginal_bigInt, _ := new(big.Int).SetString(value["balanceOriginal"], 10)
			balanceDiff_bigInt := new(big.Int).Sub(balanceOriginal_bigInt, balanceEdited_bigInt)

			blockDiffMap_value, ok := blockDiffMap[blockNum]
			if ok == false {
				blockDiffMap_value = make(map[string]*big.Int)
			}
			if transactionHash != "miner" {
				blockDiffMap_value["tx_"+transactionHash] = baseFeeDiff_bigInt
				blockDiffMap_value["addr_"+addr_key] = balanceDiff_bigInt
			} else {
				blockDiffMap_value["addr_"+addr_key] = balanceDiff_bigInt
			}

			blockDiffMap[blockNum] = blockDiffMap_value
		}
		mapAddrs = make(map[string]map[string]string)
	}

	resBlockMap := make(map[string]map[string]string)
	for key, value := range blockDiffMap {
		txsDiff := big.NewInt(0)
		txsNum := 0
		addrsDiff := big.NewInt(0)
		addrsNum := 0
		for kv, vv := range value {
			if find := strings.Contains(kv, "tx"); find {
				txsDiff = new(big.Int).Add(txsDiff, vv)
				txsNum += 1
			}
			if find := strings.Contains(kv, "addr"); find {
				addrsDiff = new(big.Int).Add(addrsDiff, vv)
				addrsNum += 1
			}
		}
		resBlockMap_value := map[string]string{"txsDiff": txsDiff.String(), "addrsDiff": addrsDiff.String(), "txsNum": strconv.Itoa(txsNum), "addrsNum": strconv.Itoa(addrsNum)}

		resBlockMap[key] = resBlockMap_value
	}

	//存入外存······················
	f, err := os.Create("../data/csv/erc20/test_blockDiff.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)         //创建一个新的写入文件流
	data := [][]string{
		{"blockNum", "txsDiff", "addrsDiff", "txsNum", "addrsNum", "same"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()

	for key, value := range resBlockMap {
		data := []string{key, value["txsDiff"], value["addrsDiff"], value["txsNum"], value["addrsNum"]}
		w.Write(data)
		w.Flush()
	}
}

// 修复map中漏掉的交易
// 1. 对原始mapAddrs进行处理，原始键为addr_string + "_ether_" + transactionHash的数据改为 addr_string + "_ether_" + blockNum并保存
//	得到mapAddrsERC20_blockNum
// 2. 遍历xblock内外部交易，如果该交易的内外部地址与blockNum构成的键存在于mapAddrsERC20_blockNum中，则保存
//	新的key为: addr_string + "_ether_" + transactionHash
// 3. 新的交易使用完整交易信息补齐

func convertKeyForOriMap() (int, int) {
	fileDir := "/mnt/4t3/geth/output/output_addrInfo_erc20/1300To/"
	fileInfoList, _ := ioutil.ReadDir(fileDir)
	var newMapAddrs = make(map[string]map[string]string)
	minBlockNum := 99999999
	maxBlockNum := 0

	for i := range fileInfoList {
		filePath := fileDir + fileInfoList[i].Name()
		load(&mapAddrs, filePath)
		for key, value := range mapAddrs {
			arr := strings.Split(key, "_")
			if arr[1] != "ether" {
				continue
			}
			addr := key[0:42]
			blockNum := value["blockNum"]
			newKey := addr + "_ether_" + blockNum
			newMapAddrs[newKey] = value

			blockNum_int, _ := strconv.Atoi(blockNum)
			if blockNum_int < minBlockNum {
				minBlockNum = blockNum_int
			}
			if blockNum_int > maxBlockNum {
				maxBlockNum = blockNum_int
			}
		}
		mapAddrs = make(map[string]map[string]string)
		fmt.Println(i)
	}
	store(newMapAddrs, "../data/dict/repair/mapAddrsERC20_blockNum.map")

	return minBlockNum, maxBlockNum
}

func dealNormalTransaction(oldMapAddrs map[string]map[string]string) {
	//[blockNumber timestamp transactionHash from to toCreate fromIsContract toIsContract
	//value gasLimit gasPrice gasUsed callingFunction isError eip2718type baseFeePerGas
	//maxFeePerGas maxPriorityFeePerGas]
	var newMapAddrs = make(map[string]map[string]string)

	fileDir := "/mnt/4t2/bowei/sbw/xblock/"
	files := [...]string{
		"13000000to13249999_BlockTransaction",
		"13250000to13499999_BlockTransaction",
		"13500000to13749999_BlockTransaction",
		"13750000to13999999_BlockTransaction",
	}
	for _, file := range files {
		fmt.Println("file: ", file)
		filePath := fileDir + file + ".zip"
		//打开并读取压缩文件中的内容
		fr, err := zip.OpenReader(filePath)
		if err != nil {
			panic(err)
		}
		defer fr.Close()
		for _, file := range fr.Reader.File {
			//判断文件该目录文件是否为文件夹
			if file.FileInfo().IsDir() {
				err := os.MkdirAll(file.Name, 0644)
				if err != nil {
					fmt.Println(err)
				}
				continue
			}
			if strings.Contains(file.Name, "csv") == false {
				continue
			}
			//为文件时，打开文件
			fs, _ := file.Open()
			r := csv.NewReader(fs)
			//针对大文件，一行一行的读取文件
			i := 0
			for {
				i += 1
				if i%100000 == 0 {
					fmt.Println(i)
				}
				row, err := r.Read()
				if err != nil && err != io.EOF {
					log.Fatalf("can not read, err is %+v", err)
				}
				if err == io.EOF {
					break
				}
				blockNumber := row[0]
				transactionHash := row[2]
				from := row[3]
				to := row[4]

				if blockNumber == "blockNumber" {
					continue
				}
				oldFromKey := from + "_ether_" + blockNumber
				oldFromValue, okFrom := oldMapAddrs[oldFromKey]
				if okFrom == true {
					newFromKey := from + "_ether_" + transactionHash
					newMapAddrs[newFromKey] = map[string]string{"addrType": "from", "isContract": oldFromValue["isContract"], "tokenAddress": "none", "tokenBalanceOriginal": "none", "tokenBalanceEdited": "none", "balanceOriginal": oldFromValue["balanceOriginal"], "balanceEdited": oldFromValue["balanceEdited"], "transactionHash": transactionHash, "lackData": "true"}
				}

				oldToKey := to + "_ether_" + blockNumber
				oldToValue, okTo := oldMapAddrs[oldToKey]
				if okTo == true {
					newToKey := to + "_ether_" + transactionHash
					newMapAddrs[newToKey] = map[string]string{"addrType": "to", "isContract": oldToValue["isContract"], "tokenAddress": "none", "tokenBalanceOriginal": "none", "tokenBalanceEdited": "none", "balanceOriginal": oldToValue["balanceOriginal"], "balanceEdited": oldToValue["balanceEdited"], "transactionHash": transactionHash, "lackData": "true"}
				}
			}
		}
		store(newMapAddrs, "../data/dict/repair/repairMap_normalTxs.map")
	}
}

func dealInterTransaction(oldMapAddrs map[string]map[string]string) {
	//[blockNumber timestamp transactionHash from to toCreate fromIsContract toIsContract
	//value gasLimit gasPrice gasUsed callingFunction isError eip2718type baseFeePerGas
	//maxFeePerGas maxPriorityFeePerGas]
	var newMapAddrs = make(map[string]map[string]string)

	fileDir := "/mnt/4t2/bowei/sbw/xblock/"
	files := [...]string{
		"13000000to13249999_InternalTransaction",
		"13250000to13499999_InternalTransaction",
		"13500000to13749999_InternalTransaction",
		"13750000to13999999_InternalTransaction",
	}
	for _, file := range files {
		fmt.Println("file: ", file)
		filePath := fileDir + file + ".zip"
		//打开并读取压缩文件中的内容
		fr, err := zip.OpenReader(filePath)
		if err != nil {
			panic(err)
		}
		defer fr.Close()
		for _, file := range fr.Reader.File {
			//判断文件该目录文件是否为文件夹
			if file.FileInfo().IsDir() {
				err := os.MkdirAll(file.Name, 0644)
				if err != nil {
					fmt.Println(err)
				}
				continue
			}
			if strings.Contains(file.Name, "csv") == false {
				continue
			}
			//为文件时，打开文件
			fs, _ := file.Open()
			r := csv.NewReader(fs)
			//针对大文件，一行一行的读取文件
			i := 0
			for {
				i += 1
				if i%100000 == 0 {
					fmt.Println(i)
				}
				row, err := r.Read()
				if err != nil && err != io.EOF {
					log.Fatalf("can not read, err is %+v", err)
				}
				if err == io.EOF {
					break
				}
				blockNumber := row[0]
				transactionHash := row[2]
				from := row[4]
				to := row[5]

				if blockNumber == "blockNumber" {
					continue
				}
				oldFromKey := from + "_ether_" + blockNumber
				oldFromValue, okFrom := oldMapAddrs[oldFromKey]
				if okFrom == true {
					newFromKey := from + "_ether_" + transactionHash
					newMapAddrs[newFromKey] = map[string]string{"addrType": "from", "isContract": oldFromValue["isContract"], "tokenAddress": "none", "tokenBalanceOriginal": "none", "tokenBalanceEdited": "none", "balanceOriginal": oldFromValue["balanceOriginal"], "balanceEdited": oldFromValue["balanceEdited"], "transactionHash": transactionHash, "lackData": "true"}
				}

				oldToKey := to + "_ether_" + blockNumber
				oldToValue, okTo := oldMapAddrs[oldToKey]
				if okTo == true {
					newToKey := to + "_ether_" + transactionHash
					newMapAddrs[newToKey] = map[string]string{"addrType": "to", "isContract": oldToValue["isContract"], "tokenAddress": "none", "tokenBalanceOriginal": "none", "tokenBalanceEdited": "none", "balanceOriginal": oldToValue["balanceOriginal"], "balanceEdited": oldToValue["balanceEdited"], "transactionHash": transactionHash, "lackData": "true"}
				}
			}
		}
		store(newMapAddrs, "../data/dict/repair/repairMap_interTxs.map")
	}
}

func dealWithROGTxData(oldMapAddrs map[string]map[string]string) {
	var mapTxs = make(map[string]map[string]string)
	fileDir := "/mnt/4t3/geth/output/output_tx/1300To/"
	fileInfoList, _ := ioutil.ReadDir(fileDir)

	for i := range fileInfoList {
		var tempMapTxs = make(map[string]map[string]string)
		filePath := fileDir + fileInfoList[i].Name()
		load(&tempMapTxs, filePath)
		mergeMap_0(mapTxs, tempMapTxs)
		tempMapTxs = make(map[string]map[string]string)

		fmt.Println(i)

		if i%50 == 0 || i == len(fileInfoList)-1 {
			for _, value := range oldMapAddrs {
				_, lackLabel := value["lackData"]
				if lackLabel == false {
					continue
				}
				transactionHash := value["transactionHash"]
				tempValue, ok := mapTxs[transactionHash]
				if ok == true {
					mergeMap_1(value, tempValue)
				}
			}
			store(oldMapAddrs, "../data/dict/repair/repairMap_ROGTx"+strconv.Itoa(i)+".map")
			mapTxs = make(map[string]map[string]string)
			// oldMapAddrs = make(map[string]map[string]string)
		}

	}

}

func addXblockData() {
	// var newMapAddrs map[string]map[string]string
	// load(&newMapAddrs, "../data/dict/repair/mapAddrsERC20_blockNum.map")

	// fmt.Print("dealNormalTransaction")
	// dealNormalTransaction(newMapAddrs)
	// fmt.Print("dealInterTransaction")
	// dealInterTransaction(newMapAddrs)

	var tempMapAddrs0 map[string]map[string]string
	load(&tempMapAddrs0, "../data/dict/repair/repairMap_normalTxs.map")
	var tempMapAddrs1 map[string]map[string]string
	load(&tempMapAddrs1, "../data/dict/repair/repairMap_interTxs.map")
	mergeMap_0(tempMapAddrs0, tempMapAddrs1)

	fmt.Print("dealWithROGTxData")
	dealWithROGTxData(tempMapAddrs0)

}

func repairMap() {
	// minBlockNum, maxBlockNum := convertKeyForOriMap()
	// fmt.Println("minBlockNum", minBlockNum)
	// fmt.Println("maxBlockNum", maxBlockNum)

	addXblockData()
}

//·······························································

//·······························································

var tokenAddrMap = make(map[string]string)

func main() {
	// getAndStoreAllTokenAddress()
	// getAddrBalanceDiff()
	// getBlockBalanceDiff()
	// getAddrBalanceDiff_specialTokens_afterSorting()

	// getAddrBalanceDiff_specialTokens_afterSorting_mutl()
	getAddrBalanceDiff_specialTokens_afterSorting_all()
	// getTxs1()
	// test()
	// getBlockDiff()
	// getBlockDiff_simple()
	// miaomi()

	// repairMap()
	// repairMap()
}

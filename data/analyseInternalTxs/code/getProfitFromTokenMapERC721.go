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
	"strings"
	// "math/big"
	// "strconv"
)

var mapAddrs = make(map[string]map[string]string)

var mapTxsForArbitrage = make(map[string]map[string]string)

// 统计每种nft受交易重排序的次数
var tokenMap = make(map[string]map[string]int)

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

type AddrInfo struct {
	address string
	info    map[string]string
}

type AddrInfos []AddrInfo

func combineMap(map0 map[string]map[string]string, map1 map[string]map[string]string) map[string]map[string]string {
	for key, value := range map1 {
		map0[key] = value
	}
	return map0
}

// 获取原始信息
func getGethOriDataSortByBlockNum() {
	var addrInfos AddrInfos
	// var resMapAddrs = make(map[string]map[string]string)
	fileDir := "/mnt/4t3/geth/output/output_tokenInfo_erc721/1300To_v1/"
	fileInfoList, err := ioutil.ReadDir(fileDir)

	for i := range fileInfoList {
		filePath := fileDir + fileInfoList[i].Name()
		load(&mapAddrs, filePath)
		for key, value := range mapAddrs {
			if value["blockNum"] == "" || value["tokenIdOwnerAddrOriginal"] == "" || value["tokenIdOwnerAddrEdited"] == "" {
				continue
			}
			addrInfos = append(addrInfos, AddrInfo{key, value})
		}
		mapAddrs = make(map[string]map[string]string)
		fmt.Println(i)
		// if i == 100 {
		// 	break
		// }
	}

	fmt.Println("addrInfos len", len(addrInfos))
	// 按照blockNum从小到大排序
	sort.SliceStable(addrInfos, func(i, j int) bool {
		i_blockNum, _ := strconv.Atoi(addrInfos[i].info["blockNum"])
		j_blockNum, _ := strconv.Atoi(addrInfos[j].info["blockNum"])
		i_position, _ := strconv.Atoi(addrInfos[i].info["positionOriginal_0"])
		j_position, _ := strconv.Atoi(addrInfos[j].info["positionOriginal_0"])
		i_tokenAddress := addrInfos[i].info["tokenAddress"]
		j_tokenAddress := addrInfos[j].info["tokenAddress"]
		i_tokenId, _ := strconv.Atoi(addrInfos[i].info["tokenId"])
		j_tokenId, _ := strconv.Atoi(addrInfos[j].info["tokenId"])

		if i_blockNum < j_blockNum {
			return true
		}
		if (i_blockNum == j_blockNum) && (i_position < j_position) {
			return true
		}
		if (i_blockNum == j_blockNum) && (i_position == j_position) && (i_tokenAddress == j_tokenAddress) && (i_tokenId < j_tokenId) {
			return true

		}
		return false
	})

	f, err := os.Create("../data/csv/erc721/tokenMapERC721ResSorted.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)         //创建一个新的写入文件流
	data := [][]string{
		{"blockNum", "tokenAddress", "tokenId", "tokenIdOwnerAddrOriginal", "tokenIdOwnerAddrEdited", "endAddr_0", "transactionHash_0", "positionOriginal_0", "positionEdited_0", "endAddr_1", "transactionHash_1", "positionOriginal_1", "positionEdited_1"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()
	for _, item := range addrInfos {
		value := item.info

		data := [][]string{
			{value["blockNum"], value["tokenAddress"], value["tokenId"], value["tokenIdOwnerAddrOriginal"], value["tokenIdOwnerAddrEdited"], value["endAddr_0"], value["transactionHash_0"], value["positionOriginal_0"], value["positionEdited_0"], value["endAddr_1"], value["transactionHash_1"], value["positionOriginal_1"], value["positionEdited_1"]},
		}
		w.WriteAll(data) //写入数据
		w.Flush()
	}
}

// 统计每种nft受影响的次数
func getNftCount() {
	tokenCountMap := make(map[string]map[string]int)
	fileDir := "/mnt/4t3/geth/output/output_tokenInfo_erc721/1300To_v1/"
	fileInfoList, err := ioutil.ReadDir(fileDir)

	for i := range fileInfoList {
		filePath := fileDir + fileInfoList[i].Name()
		load(&mapAddrs, filePath)
		for key, value := range mapAddrs {
			if value["blockNum"] == "" || value["tokenIdOwnerAddrOriginal"] == "" || value["tokenIdOwnerAddrEdited"] == "" {
				continue
			}
			blockNum_int, _ := strconv.Atoi(value["blockNum"])
			arr := strings.Split(key, "_")
			tokenAddr := arr[0]
			tokenCountMap_value, ok := tokenCountMap[tokenAddr]
			if ok == false {
				tokenCountMap_value = make(map[string]int)
				tokenCountMap_value["count"] = 1
				tokenCountMap_value["minBlockNum"] = blockNum_int
				tokenCountMap_value["maxBlockNum"] = blockNum_int
			} else {
				tokenCountMap_value["count"] += 1
				if tokenCountMap_value["minBlockNum"] > blockNum_int {
					tokenCountMap_value["minBlockNum"] = blockNum_int
				}
				if tokenCountMap_value["maxBlockNum"] < blockNum_int {
					tokenCountMap_value["maxBlockNum"] = blockNum_int
				}
			}
			tokenCountMap[tokenAddr] = tokenCountMap_value
		}
		mapAddrs = make(map[string]map[string]string)
		fmt.Println(i)
	}

	var tokenAddrInfos AddrInfos
	for key, value := range tokenCountMap {
		newValue := make(map[string]string)
		newValue["minBlockNum"] = strconv.Itoa(value["minBlockNum"])
		newValue["maxBlockNum"] = strconv.Itoa(value["maxBlockNum"])
		newValue["count"] = strconv.Itoa(value["count"])
		tokenAddrInfos = append(tokenAddrInfos, AddrInfo{key, newValue})
	}
	// 按照blockNum从小到大排序
	sort.SliceStable(tokenAddrInfos, func(i, j int) bool {
		i_count, _ := strconv.Atoi(tokenAddrInfos[i].info["count"])
		j_count, _ := strconv.Atoi(tokenAddrInfos[j].info["count"])
		i_minBlockNum, _ := strconv.Atoi(tokenAddrInfos[i].info["minBlockNum"])
		j_minBlockNum, _ := strconv.Atoi(tokenAddrInfos[j].info["minBlockNum"])

		if i_count < j_count {
			return false
		}
		if (i_count == j_count) && (i_minBlockNum < j_minBlockNum) {
			return true
		}
		return true
	})
	f, err := os.Create("../data/csv/erc721/nftCount.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)         //创建一个新的写入文件流
	data := [][]string{
		{"tokenAddress", "count", "minBlockNum", "maxBlockNum"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()
	for _, item := range tokenAddrInfos {
		key := item.address
		value := item.info

		data := [][]string{
			{key, value["count"], value["minBlockNum"], value["maxBlockNum"]},
		}
		w.WriteAll(data) //写入数据
		w.Flush()
	}

}

func main() {
	// getGethOriDataSortByBlockNum()
	getNftCount()
}

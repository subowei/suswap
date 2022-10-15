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
	fileDir := "/mnt/4t3/geth/output/output_addrInfo_erc721/1300To_v1/"
	fileInfoList, err := ioutil.ReadDir(fileDir)

	for i := range fileInfoList {
		filePath := fileDir + fileInfoList[i].Name()
		load(&mapAddrs, filePath)
		for key, value := range mapAddrs {
			// addr := strings.ToLower(key[0:42])
			if value["tokenIdOwnerAddrOriginal"] == value["tokenIdOwnerAddrEdited"] {
				delete(mapAddrs, key)
			}
		}
		fmt.Println(i)
		if i == 60 {
			break
		}
	}

	for key, value := range mapAddrs {
		addrInfos = append(addrInfos, AddrInfo{key, value})
	}

	fmt.Println("addrInfos len", len(addrInfos))
	// 按照blockNum从小到大排序
	sort.SliceStable(addrInfos, func(i, j int) bool {
		i_blockNum, _ := strconv.Atoi(addrInfos[i].info["blockNum"])
		j_blockNum, _ := strconv.Atoi(addrInfos[j].info["blockNum"])
		i_position, _ := strconv.Atoi(addrInfos[i].info["positionOriginal"])
		j_position, _ := strconv.Atoi(addrInfos[j].info["positionOriginal"])
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

	f, err := os.Create("../data/csv/erc721/addrMapERC721ResSorted_lostToken_0to60.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)         //创建一个新的写入文件流
	data := [][]string{
		{"keyAddr", "tokenAddress", "tokenId", "tokenIdOwnerAddrOriginal", "tokenIdOwnerAddrEdited", "tokenBalanceOriginal", "tokenBalanceEdited", "blockNum", "transactionHash", "positionOriginal", "positionEdited"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()
	for _, item := range addrInfos {
		key := item.address
		// fmt.Println("key",key)
		value := item.info

		data := [][]string{
			{key[0:42], value["tokenAddress"], value["tokenId"], value["tokenIdOwnerAddrOriginal"], value["tokenIdOwnerAddrEdited"], value["tokenBalanceOriginal"], value["tokenBalanceEdited"], value["blockNum"], value["transactionHash"], value["positionOriginal"], value["positionEdited"]},
		}
		w.WriteAll(data) //写入数据
		w.Flush()
	}
}

// 统计每种nft受影响的次数
func getNftCount() {
	allCount := 0
	allMinBlockNum := 99999999999
	allMaxBlockNum := 0
	var paths []string = []string{"../../output/output_addrInfo_erc721/1200To/mapAddrsERC721_12959999", "../../output/output_addrInfo_erc721/1280To/mapAddrsERC721_13009999", "../../output/output_addrInfo_erc721/1300To/mapAddrsERC721_13069999", "../../output/output_addrInfo_erc721/1303To/mapAddrsERC721_13089999"}
	// var paths []string = []string{"../../output/output_addrInfo_erc721/1200To/mapAddrsERC721_12959999"}
	for _, path := range paths {
		load(&mapAddrs, path)
		for key, value := range mapAddrs {
			addr := strings.ToLower(key[0:42])
			if value["tokenIdOwnerAddrOriginal"] == value["tokenIdOwnerAddrEdited"] {
				continue
			}
			if addr != value["tokenIdOwnerAddrOriginal"] {
				continue
			}
			allCount += 1
			//存入tokenmap
			blockNum_int, _ := strconv.Atoi(value["blockNum"])
			if allMinBlockNum > blockNum_int {
				allMinBlockNum = blockNum_int
			}
			if allMaxBlockNum < blockNum_int {
				allMaxBlockNum = blockNum_int
			}

			tokenAddress := value["tokenAddress"]
			tokenMap_value, ok := tokenMap[tokenAddress]
			if ok == false {
				tokenMap_value = make(map[string]int)
				tokenMap_value["minBlockNum"] = blockNum_int
				tokenMap_value["maxBlockNum"] = blockNum_int
				tokenMap_value["count"] = 1

			} else {
				tokenMap_value["count"] += 1
				tokenMap_value["maxBlockNum"] = blockNum_int
			}
			tokenMap[tokenAddress] = tokenMap_value
		}
	}

	var tokenAddrInfos AddrInfos
	for key, value := range tokenMap {
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
	f, err := os.Create("../data/csv/nftCount.csv") //创建文件
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

	fmt.Println(allCount)
	fmt.Println(allMinBlockNum)
	fmt.Println(allMaxBlockNum)
}

func main() {
	getGethOriDataSortByBlockNum()
	// getNftCount()
	// getTxs0()
}

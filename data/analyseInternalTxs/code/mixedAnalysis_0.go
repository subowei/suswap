package main

import (
	// "fmt"
	"bytes"
	// "encoding/csv"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math/big"
	// "os"
	// "sort"
	// "strings"
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


// key:txhash	value:外部内部to地址，交易手续费
// txHash1, toA-toD, txFees1
func getData(){
	var addrsMap = make(map[string]map[string]string)
	var resMap = make(map[string]map[string]string)

	fileDir:="/mnt/sbw/reorder_geth-dev/data/output/output_addrInfo/1300To/"
	fileInfoList,_ := ioutil.ReadDir(fileDir)

	for i := range fileInfoList {
		filePath:=fileDir+fileInfoList[i].Name()
		load(&addrsMap, filePath)
		for key, value := range addrsMap {
			addr:=key[0:42]
			transactionHash := value["transactionHash"]
			resMap_value,ok:=resMap[transactionHash]
			if ok==false{
				resMap_value=make(map[string]string)
				effectiveGasTipValue:=value["effectiveGasTipValue"]
				gasUsedOriginal:=value["gasUsedOriginal"]
				addrType:=value["addrType"]
				if effectiveGasTipValue=="" || gasUsedOriginal=="" || addrType=="from"{
					continue
				}

				// resMap_value["effectiveGasTipValue"]=value["effectiveGasTipValue"]
				// resMap_value["gasPrice"]=value["gasPrice"]
				// resMap_value["gasUsedOriginal"]=value["gasUsedOriginal"]
				// resMap_value["baseFee"]=value["baseFee"]

				resMap_value["allToAddr"]=addr
				resMap_value["totalFeeOriginal"]=value["totalFeeOriginal"]
				
			}else{
				resMap_value["allToAddr"]=resMap_value["allToAddr"]+";"+addr
			}
			resMap[transactionHash]=resMap_value
		}

		for key, value := range resMap {
			fmt.Println("key",key)
			fmt.Println("value",value)
			return
		}

	}
}

// 函数：getData0:
// 输入：output_addrInfo_erc20
// 作用：获得dict{交易哈希:{余额发生变化的地址总数}}
var resTxMap = make(map[string]map[string]*big.Int)
func getData0(){
	var addrsMap = make(map[string]map[string]string)
	var tempTxsMap = make(map[string]int64)

	fileDir:="/mnt/sbw/reorder_geth-dev/data/output/output_addrInfo/1300To/"
	fileInfoList,_ := ioutil.ReadDir(fileDir)

	for i := range fileInfoList {
		filePath:=fileDir+fileInfoList[i].Name()
		load(&addrsMap, filePath)
}



// 函数：getData1:
// 输入：output_addrInfo
// 作用：获得dict{to地址:{涉及交易总数，交易费总数}}
func getData1(){
	var resAddrMap = make(map[string]map[string]string)
	var addrsMap = make(map[string]map[string]string)
	var tempAddrMap = make(map[string]map[string]*big.Int)

	fileDir:="/mnt/sbw/reorder_geth-dev/data/output/output_addrInfo/1300To/"
	fileInfoList,_ := ioutil.ReadDir(fileDir)

	for i := range fileInfoList {
		filePath:=fileDir+fileInfoList[i].Name()
		load(&addrsMap, filePath)
		for key, value := range addrsMap {
			addrKey:=key[0:42]
			effectiveGasTipValue:=value["effectiveGasTipValue"]
			gasUsedOriginal:=value["gasUsedOriginal"]
			addrType:=value["addrType"]
			totalFeeOriginal_bigInt,_:=new(big.Int).SetString(value["totalFeeOriginal"], 10)

			if effectiveGasTipValue=="" || gasUsedOriginal=="" || addrType=="from"{
				continue
			}

			tempAddrMap_value,ok:=tempAddrMap[addrKey]
			if ok==false{
				tempAddrMap_value=make(map[string]*big.Int)

				tempAddrMap_value["numOfTxs"]=big.NewInt(1)
				tempAddrMap_value["totalFee"]=totalFeeOriginal_bigInt
			}else{
				tempAddrMap_value["numOfTxs"]=new(big.Int).Add(tempAddrMap_value["numOfTxs"],big.NewInt(1))
				tempAddrMap_value["totalFee"]=new(big.Int).Add(tempAddrMap_value["totalFee"],totalFeeOriginal_bigInt)
			}
			tempAddrMap[addrKey]=tempAddrMap_value
		}
	}
	for key, value := range tempAddrMap {
		tempMap:=make(map[string]string)
		for vk, vv := range value{
			tempMap[vk]=vv.String()
		}
		resAddrMap[key]=tempMap
	}

	// for key, value := range resAddrMap {
	// 	fmt.Println("key",key)
	// 	fmt.Println("value",value)
	// 	return
	// }
}







func main() {
	getData0()
}
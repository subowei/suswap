package main

import (

	// "fmt"
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
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

type TxInfo struct {
	transactionHash string
	info map[string]string
}

type TxInfos []TxInfo

func main() {
	//mapTx := map[string]string{"blockNum", "miner", "fromAddrEdited", "toAddrEdited", "fromBalanceEdited", "toBalanceEdited": toBalance, "fromAccountRootEdited": fromAccountRootEdited, "toAccountRootEdited": toAccountRootEdited, "fromBalanceOriginal", "toBalanceOriginal", "fromAccountRootOriginal", "toAccountRootOriginal"}
	f, err := os.Create("../data/csv/gethInternalTxs_oriOrder_v4.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	data := [][]string{
		{"transactionHash", "blockNum", "miner", "fromAddr", "toAddr", "fromBalanceEdited", "toBalanceEdited", "fromAccountRootEdited", "toAccountRootEdited", "fromBalanceOriginal", "toBalanceOriginal", "fromAccountRootOriginal", "toAccountRootOriginal","gasPrice","positionOriginal","positionEdited","serviceChargeOriginal","serviceChargeEdited","isChanged"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()

	var mapTxs0 = make(map[string]map[string]string)

	// load(&mapTxs0, "../../output_internaltxs/400wToV3/mapTxs4899499")
	load(&mapTxs0, "../../output_internaltxs/500wToV3/mapTxs5087999")

	// 剔除不存在余额差的区块
	var validBlock = make(map[string]bool)
	for _, v := range mapTxs0 {
		if v["isChanged"]=="balance" {
			validBlock[ v["blockNum"] ]=true;
		}

	}

	// 加入结构体，并剔除空值
	var txInfos TxInfos
	for k, v := range mapTxs0 {
		// print(v["position"])
		// if v["gasPrice"]=="" {continue}
		blockNum:=v["blockNum"]
		_,ok:=validBlock[blockNum]
		if ok==false{
			continue
		}


		if v["positionOriginal"]==""{
			v["positionOriginal"]="999999"
		}
		if v["positionEdited"]==""{
			v["positionEdited"]="999999"
		}
		mapTxs0[k]=v
		txInfos = append(txInfos, TxInfo {k, v})
	}
	// 按照blockNum从小到大排序
	sort.SliceStable(txInfos, func(i, j int) bool {
		i_blockNum,_ := strconv.Atoi(txInfos[i].info["blockNum"])
		j_blockNum,_ := strconv.Atoi(txInfos[j].info["blockNum"])
		i_position,_ := strconv.Atoi(txInfos[i].info["positionOriginal"])
		j_position,_ := strconv.Atoi(txInfos[j].info["positionOriginal"])

		if i_blockNum < j_blockNum{
			return true
		}
		if (i_blockNum==j_blockNum) && (i_position<j_position){
			return true
		}
		return false
	})

	for _,item :=range txInfos{
		k:= item.transactionHash
		v:= item.info

		data := [][]string{
			{k, v["blockNum"], v["miner"], v["fromAddr"], v["toAddr"], v["fromBalanceEdited"], v["toBalanceEdited"], v["fromAccountRootEdited"], v["toAccountRootEdited"], v["fromBalanceOriginal"], v["toBalanceOriginal"], v["fromAccountRootOriginal"], v["toAccountRootOriginal"],v["gasPrice"],v["positionOriginal"],v["positionEdited"],v["serviceChargeOriginal"],v["serviceChargeEdited"],v["isChanged"]},
		}
		w.WriteAll(data) //写入数据
		w.Flush()
	}


	//````````````````````````````````````````````````````````````````````````````````````````````````
	f, err = os.Create("../data/csv/gethInternalTxs_editedOrder_v4.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w = csv.NewWriter(f) //创建一个新的写入文件流
	data = [][]string{
		{"transactionHash", "blockNum", "miner", "fromAddr", "toAddr", "fromBalanceEdited", "toBalanceEdited", "fromAccountRootEdited", "toAccountRootEdited", "fromBalanceOriginal", "toBalanceOriginal", "fromAccountRootOriginal", "toAccountRootOriginal","gasPrice","positionOriginal","positionEdited","serviceChargeOriginal","serviceChargeEdited","isChanged"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()

	sort.SliceStable(txInfos, func(i, j int) bool {
		i_blockNum,_ := strconv.Atoi(txInfos[i].info["blockNum"])
		j_blockNum,_ := strconv.Atoi(txInfos[j].info["blockNum"])
		i_position,_ := strconv.Atoi(txInfos[i].info["positionEdited"])
		j_position,_ := strconv.Atoi(txInfos[j].info["positionEdited"])

		if i_blockNum < j_blockNum{
			return true
		}
		if (i_blockNum==j_blockNum) && (i_position<j_position) {
			return true
		}
		return false
	})

	for _,item :=range txInfos{
		k:= item.transactionHash
		v:= item.info

		data := [][]string{
			{k, v["blockNum"], v["miner"], v["fromAddr"], v["toAddr"], v["fromBalanceEdited"], v["toBalanceEdited"], v["fromAccountRootEdited"], v["toAccountRootEdited"], v["fromBalanceOriginal"], v["toBalanceOriginal"], v["fromAccountRootOriginal"], v["toAccountRootOriginal"],v["gasPrice"],v["positionOriginal"],v["positionEdited"],v["serviceChargeOriginal"],v["serviceChargeEdited"],v["isChanged"]},
		}
		w.WriteAll(data) //写入数据
		w.Flush()
	}

}

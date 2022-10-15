package main

import (

	// "fmt"
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"io/ioutil"
	"os"
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

func main() {
	//mapTx := map[string]string{"blockNum", "miner", "fromAddrEdited", "toAddrEdited", "fromBalanceEdited", "toBalanceEdited": toBalance, "fromAccountRootEdited": fromAccountRootEdited, "toAccountRootEdited": toAccountRootEdited, "fromBalanceOrignal", "toBalanceOrignal", "fromAccountRootOrignal", "toAccountRootOrignal"}
	f, err := os.Create("../data/csv/gethInternalTxs.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	data := [][]string{
		{"transactionHash", "blockNum", "miner", "fromAddr", "toAddr", "fromBalanceEdited", "toBalanceEdited", "fromAccountRootEdited", "toAccountRootEdited", "fromBalanceOrignal", "toBalanceOrignal", "fromAccountRootOrignal", "toAccountRootOrignal","gasPrice","position"},
	}
	w.WriteAll(data) //写入数据
	w.Flush()

	var mapTxs0 = make(map[string]map[string]string)
	// var mapTxs1 = make(map[string]map[string]string)
	// var mapTxs2 = make(map[string]map[string]string)
	// var mapTxs3 = make(map[string]map[string]string)

	load(&mapTxs0, "../../output_internaltxs/500wToV2/mapTxs5199999")
	// load(&mapTxs1, "../outputnew/600wTo/mapTxs180448")
	// load(&mapTxs2, "../outputnew/670wTo/mapTxs48621")
	// load(&mapTxs3, "../outputnew/800wTo/mapTxs6748")

	for k, v := range mapTxs0 {
		// fmt.Println(k, v)
		data := [][]string{
			{k, v["blockNum"], v["miner"], v["fromAddr"], v["toAddr"], v["fromBalanceEdited"], v["toBalanceEdited"], v["fromAccountRootEdited"], v["toAccountRootEdited"], v["fromBalanceOrignal"], v["toBalanceOrignal"], v["fromAccountRootOrignal"], v["toAccountRootOrignal"],v["gasPrice"],v["position"]},
		}
		w.WriteAll(data) //写入数据
		w.Flush()
	}

	// for k, v := range mapTxs1 {
	// 	// fmt.Println(k, v)
	// 	data := [][]string{
	// 		{k, v["blockNum"], v["miner"], v["fromAddrEdited"], v["toAddrEdited"], v["fromBalanceEdited"], v["toBalanceEdited"], v["fromAccountRootEdited"], v["toAccountRootEdited"], v["fromBalanceOrignal"], v["toBalanceOrignal"], v["fromAccountRootOrignal"], v["toAccountRootOrignal"]},
	// 	}
	// 	w.WriteAll(data) //写入数据
	// 	w.Flush()
	// }

	// for k, v := range mapTxs2 {
	// 	// fmt.Println(k, v)
	// 	data := [][]string{
	// 		{k, v["blockNum"], v["miner"], v["fromAddrEdited"], v["toAddrEdited"], v["fromBalanceEdited"], v["toBalanceEdited"], v["fromAccountRootEdited"], v["toAccountRootEdited"], v["fromBalanceOrignal"], v["toBalanceOrignal"], v["fromAccountRootOrignal"], v["toAccountRootOrignal"]},
	// 	}
	// 	w.WriteAll(data) //写入数据
	// 	w.Flush()
	// }

	// for k, v := range mapTxs3 {
	// 	// fmt.Println(k, v)
	// 	data := [][]string{
	// 		{k, v["blockNum"], v["miner"], v["fromAddrEdited"], v["toAddrEdited"], v["fromBalanceEdited"], v["toBalanceEdited"], v["fromAccountRootEdited"], v["toAccountRootEdited"], v["fromBalanceOrignal"], v["toBalanceOrignal"], v["fromAccountRootOrignal"], v["toAccountRootOrignal"]},
	// 	}
	// 	w.WriteAll(data) //写入数据
	// 	w.Flush()
	// }

}

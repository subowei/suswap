from asyncore import read
from distutils.command.build_scripts import first_line_re
from re import L
import pandas as pd
import datetime
import numpy as np
import csv    #加载csv包便于读取csv文件
import zipfile
import os
import json

minBlockNum=13317502
maxBlockNum=13493995

#····························································
# 第一步骤：获取待研究交易哈希
txs_addrNum_Map={}
def getTransactionFromROG():
    print("getTransactionFromROG")
    df = pd.read_csv('../data/csv/erc20/txs_addrNum_Map.csv')
    for index, row in df.iterrows():
        transactionHash=row["transactionHash"]
        numOfAddr=row["numOfAddr"]
        txs_addrNum_Map[transactionHash]=numOfAddr

# 第二步骤：获取交易对应的内部to地址
res_txs_map={}
def getInterTransactionFromXblock():
    print("getInterTransactionFromXblock")
    fileDir = "/mnt/4t2/bowei/sbw/xblock/";

    files = [
        "13250000to13499999_InternalTransaction",
    ];

    for file in files:
        theZIP = zipfile.ZipFile(fileDir+file+".zip", 'r');
        theCSV = theZIP.open(file+".csv");

        head = theCSV.readline();
        oneLine = theCSV.readline().decode("utf-8").strip();
        while (oneLine!=""):
            oneArray = oneLine.split(",")
            transactionHash=oneArray[2]
            toAddr=oneArray[5]
            try:
                numOfAddr=txs_addrNum_Map[transactionHash]
                if transactionHash in res_txs_map.keys():
                    value=res_txs_map[transactionHash]
                    value["toAddrs"]=value["toAddrs"]+";"+toAddr
                else:
                    value={}
                    value["toAddrs"]=toAddr
                    value["numOfAddr"]=numOfAddr
                res_txs_map[transactionHash]=value
            except:
                pass
            oneLine = theCSV.readline().decode("utf-8").strip();

            
# 第三步骤：获取交易对应的外部to地址，交易手续费
# ['blockNumber', 'timestamp', 'transactionHash', 'from', 'to', 'toCreate', 
# 'fromIsContract', 'toIsContract', 'value', 'gasLimit', 'gasPrice', 'gasUsed', 
# 'callingFunction', 'isError', 'eip2718type', 'baseFeePerGas', 'maxFeePerGas', 'maxPriorityFeePerGas']
def getNormalTransactionFromXblock():
    print("getNormalTransactionFromXblock")
    fileDir = "/mnt/4t2/bowei/sbw/xblock/";

    files = [
        "13250000to13499999_BlockTransaction",
    ];

    for file in files:
        theZIP = zipfile.ZipFile(fileDir+file+".zip", 'r');
        theCSV = theZIP.open(file+".csv");

        head = theCSV.readline().decode("utf-8").strip();
        oneLine = theCSV.readline().decode("utf-8").strip();
        # title=head.split(",")
        # print("title",title)
        # return
        while (oneLine!=""):
            oneArray = oneLine.split(",")

            blockNumber=oneArray[0]
            transactionHash=oneArray[2]
            toAddr=oneArray[4]
            gasPrice=int(oneArray[10])
            gasUsed=int(oneArray[11])
            txFees=gasPrice*gasUsed

            try:       
                numOfAddr=txs_addrNum_Map[transactionHash]

                if transactionHash in res_txs_map.keys():
                    value=res_txs_map[transactionHash]
                    value["toAddrs"]=value["toAddrs"]+";"+toAddr
                    value["numOfAddr"]=numOfAddr
                    value["txFees"]=str(txFees)
                else:
                    value={}
                    value["toAddrs"]=toAddr
                    value["numOfAddr"]=numOfAddr
                    value["txFees"]=str(txFees)
                res_txs_map[transactionHash]=value
            except:
                pass
            oneLine = theCSV.readline().decode("utf-8").strip();

# 第四步骤：输出res_txs_map
def outPutMap0():
    print("outPutMap0")
    global outputCsvName0
    f = open(outputCsvName0,'w')
    writer = csv.writer(f)
    writer.writerow(["transactionHash","txFees","numOfAddr","toAddrs"])

    for key,value in res_txs_map.items():
        # value["txFees"]可能不存在，由于两个文件的差异导致
        row=[key,value["txFees"],value["numOfAddr"],value["toAddrs"]]
        writer.writerow(row)

# output: txsDataPlusXblock.csv
# csv type: transactionHash,txFees,numOfAddr,toAddrs 
def getTxData():
    getTransactionFromROG()
    getInterTransactionFromXblock()
    getNormalTransactionFromXblock()
    outPutMap0()

#····························································
temp_toAddr_map={}
res_toAddr_map={}
# input: txsDataPlusXblock.csv
# output: addrDataPlusXblock.csv
# csv type: toAddr,numOfTxs,totalTxFees,numOfAddrs,
def getAddrData():
    csvPath="../data/csv/erc20/txsDataPlusXblock.csv"
    df = pd.read_csv(csvPath)
    for index, row in df.iterrows():
        transactionHash=row["transactionHash"]
        txFees=row["txFees"]
        numOfAddr=row["numOfAddr"]
        toAddrs=row["toAddrs"]

        addrsList=toAddrs.split(";")
        for addr in addrsList:            
            if addr in temp_toAddr_map.keys():
                value=temp_toAddr_map[addr]
                value[transactionHash]={"txFees":txFees,"numOfAddr":numOfAddr}
                temp_toAddr_map[addr]=value
            else:
                temp_toAddr_map[addr]={transactionHash:{"txFees":txFees,"numOfAddr":numOfAddr}}

    for key,value in temp_toAddr_map.items():
        tempMap={"numOfTxs":0,"totalTxFees":0,"numOfAddrs":0}
        for vk,vv in value.items():
            tempMap["numOfTxs"]+=1
            tempMap["totalTxFees"]+=int(vv["txFees"])
            tempMap["numOfAddrs"]+=int(vv["numOfAddr"])

        res_toAddr_map[key]=tempMap

    global outputCsvName1
    f = open(outputCsvName1,'w')
    writer = csv.writer(f)
    writer.writerow(["toAddr","numOfTxs","totalTxFees","numOfAddrs"])

    for key,value in res_toAddr_map.items():
        row=[key,str(value["numOfTxs"]),str(value["totalTxFees"]),str(value["numOfAddrs"])]
        writer.writerow(row)


#····························································
temp_toAddr_Map={}
def getToAddrMap():
    df = pd.read_csv('../data/csv/erc20/addrDataPlusXblock.csv')
    for index, row in df.iterrows():
        temp_toAddr_Map[row["toAddr"]]={"numOfTxs":row["numOfTxs"],"totalTxFees":row["totalTxFees"],"numOfAddrs":row["numOfAddrs"]}

def getDefi_csv():
    fileDir = "/mnt/4t2/bowei/sbw/xblock/peilin_defi/"
    filesName=["BalancerV1_PoolInfo.csv","ShibaSwap_PairInfo.csv","SushiSwap_PairInfo.csv","UniswapV1_ExchangeInfo.csv","UniswapV2_PairInfo.csv","UniswapV3_PoolInfo.csv"]

    # BalancerV1_PoolInfo: poolAddress,creator
    # ShibaSwap_PairInfo: pairAddress,tokenAddress0,tokenAddress1
    # SushiSwap_PairInfo: pairAddress,tokenAddress0,tokenAddress1
    # UniswapV1_ExchangeInfo: exchangeAddress,tokenAddress
    # UniswapV2_PairInfo: pairAddress,tokenAddress0,tokenAddress1
    # UniswapV3_PoolInfo: poolAddress,tokenAddress0,tokenAddress1
    for fileName in filesName:
        # print("fileName",fileName)
        defiName=fileName.split('_')[0]
        with open(fileDir+fileName) as f:
            reader=csv.reader(f)
            header_row=next(reader)
            for row in reader:
                comparedAddr=row[0]
                try:
                    value=temp_toAddr_Map[comparedAddr]
                    value["defiName"]=defiName
                    value["defiType"]=header_row[0]
                    temp_toAddr_Map[comparedAddr]=value
                except:
                    pass


def getDefi_curveJson():
    filePath="/mnt/4t2/bowei/sbw/xblock/peilin_defi/Curve_PoolInfo.json"
    defiName="Curve_PoolInfo"
    with open(filePath) as f:
        jsonFile=json.load(f)
        for key,value in jsonFile.items():
            try:
                value=temp_toAddr_Map[key]
                value["defiName"]=defiName+"_"+value["name"]
                value["defiType"]=value["type"]
                temp_toAddr_Map[key]=value
            except:
                pass

def getDefi_makerDAOJson():
    filePath="/mnt/4t2/bowei/sbw/xblock/peilin_defi/MakerDAO_CollateralInfo.json"
    defiName="makerDAO"
    with open(filePath) as f:
        jsonFile=json.load(f)
        for key,value in jsonFile.items():
            try:
                value=temp_toAddr_Map[key]
                value["defiName"]=defiName
                value["defiType"]="none"
                temp_toAddr_Map[key]=value
            except:
                pass
            
def getTokenInfo():
    csv.field_size_limit(500 * 1024 * 1024)
    with open("/mnt/4t2/bowei/sbw/xblock/ERC20TokenInfo.csv") as f:
        reader = csv.reader(_.replace('\x00', '') for _ in f)
        header_row=next(reader)
        for row in reader:
            address=row[0]
            name=row[1]
            symbol=row[2]
            try:
                value=temp_toAddr_Map[address]
                value["symbol"]=symbol
            except:
                pass
            
def getContractInfo():
    dir="/mnt/4t2/bowei/sbw/xblock/VerifiedContract/"
    for fileName in os.listdir(dir):
        filePath=dir+fileName
        with open(filePath) as f:
            jsonFile=json.load(f)
            for key,value in jsonFile.items():
                address=key
                contractName=value["ContractName"]
                try:
                    value=temp_toAddr_Map[address]
                    value["name"]=contractName
                except:
                    pass
            
        

def outPutMap1():
    global outputCsvName2
    f = open(outputCsvName2,'w')
    writer = csv.writer(f)
    writer.writerow(["toAddr","numOfTxs","totalTxFees","numOfAddrs","defiName","defiType","name","symbol"])

    for key,value in temp_toAddr_Map.items():
        defiName="none"
        defiType="none"
        name="none"
        symbol="none"
        
        if "defiName" in value.keys():
            defiName=value["defiName"]
        if "defiType" in value.keys():
            defiType=value["defiType"]
        if "name" in value.keys():
            name=value["name"]
        if "symbol" in value.keys():
            symbol=value["symbol"]
            
        row=[key,value["numOfTxs"],value["totalTxFees"],value["numOfAddrs"],defiName,defiType,name,symbol]
        writer.writerow(row)

# input: addrDataPlusXblock.csv
# output: addrDataPlusXblockWithDefi.csv
def plusDefi():
    getToAddrMap()
    getDefi_csv()
    getDefi_curveJson()
    getDefi_makerDAOJson()
    getTokenInfo()
    getContractInfo()
    outPutMap1()
    
    
#····························································        
def getAddrData(targetAddrs):
    csvPath="../data/csv/erc20/txsDataPlusXblock.csv"
    df = pd.read_csv(csvPath)
    for index, row in df.iterrows():
        transactionHash=row["transactionHash"]
        toAddrs=row["toAddrs"]

        addrsList=toAddrs.split(";")
        for addr in addrsList:            
            if addr in targetAddrs:
                if addr in temp_toAddr_map.keys():
                    value=temp_toAddr_map[addr]
                    value[transactionHash]=1
                    temp_toAddr_map[addr]=value
                else:
                    value={transactionHash:1}
                    temp_toAddr_map[addr]=value
                    

def cropCsv(targetAddr,dir):
    filePath=dir+targetAddr+".csv"
    fNew = open(filePath,'w')
    writer = csv.writer(fNew)
    with open("../data/csv/erc20/addrBalanceDiff_allTokens_syndromicChange_new.csv") as fOld:
        reader = csv.reader(fOld)
        header_row=next(reader)
        writer.writerow(header_row)
        for row in reader:
            transactionHash=row[3]
            if transactionHash in temp_toAddr_map[targetAddr].keys():
                writer.writerow(row)
                
            
def findTxsRelatedSpecialAddr():
    targetAddrs=["0x383518188c0c6d7730d91b2c03a03c837814a899","0xc4ad29ba4b3c580e6d59105fff484999997675ff","0x4e3fbd56cd56c3e72c1403e103b45db9da5b9d2b","0x6c3f90f043a72fa612cbac8115ee7e52bde6e490"]
    dir="../data/addr_tx_v1/"
    getAddrData(targetAddrs)
    for targetAddr in targetAddrs:
        cropCsv(targetAddr,dir)
            
#····························································

outputCsvName0="../data/csv/erc20/txsDataPlusXblock.csv"
outputCsvName1="../data/csv/erc20/addrDataPlusXblock.csv"
outputCsvName2="../data/csv/erc20/addrDataPlusXblockWithDefi.csv"

def main():
    print("getTxData")
    getTxData()

    print("getAddrData")
    getAddrData()

    print("plusDefi")
    plusDefi()
    
    
    findTxsRelatedSpecialAddr()

if __name__ == '__main__':
    main()
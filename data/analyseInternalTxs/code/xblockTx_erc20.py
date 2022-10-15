from asyncore import read
from distutils.command.build_scripts import first_line_re
from re import L, T
import pandas as pd
import datetime
import numpy as np
import csv    #加载csv包便于读取csv文件
import zipfile
import os
import json
import pickle

minBlockNum=13317502
maxBlockNum=13493995

#····························································
# 第一步骤：获取待研究交易哈希
def getTransactionFromROG(txMap):
    print("getTransactionFromROG")        
    with open("../data/csv/erc20/addrBalanceDiff_allTokens_syndromicChange_minerIncome.csv") as f:
        reader = csv.reader(f)
        header_row=next(reader)
        for row in reader:
            address=row[0]
            transactionHash=row[3].lower()
            effectiveFeeOriginal=row[7]
            value2miner=row[8]
            mineIncome=row[9]
            gasUsedOriginal=row[10]
            
            if gasUsedOriginal=="21000":
                continue
            
            if transactionHash!=txMap.keys():
                valueMap={"effectiveFeeOriginal":effectiveFeeOriginal,"value2miner":value2miner,"mineIncome":mineIncome}
                valueMap["address"]=[address]
                txMap[transactionHash]=valueMap
            else:
                valueMap=txMap[transactionHash]
                valueMap["address"].append(address)
                txMap[transactionHash]=valueMap
    
    for key,value in txMap.items():
        value["numOfAddr"]=len(set(value["address"]))
        value["address"]=None
        
    global outputCsvName3
    f = open(outputCsvName3,'w')
    writer = csv.writer(f)
    writer.writerow(["transactionHash","numOfAddr","effectiveFeeOriginal","value2miner","mineIncome"])

    for key,value in txMap.items():
        row=[key,value["numOfAddr"],value["effectiveFeeOriginal"],value["value2miner"],value["mineIncome"]]
        writer.writerow(row)
        
    with open("/mnt/4t3/geth/reorder_geth-dev3/data/analyseInternalTxs/data/dict/erc20/txMap0.pkl", "wb") as tf:
        pickle.dump(txMap,tf)

# 第二步骤：获取交易对应的内部to地址
def getInterTransactionFromXblock(txMap):
    print("getInterTransactionFromXblock")
    fileDir = "/mnt/4t2/bowei/sbw/xblock/";

    files = [
		"13000000to13249999_InternalTransaction",
		"13250000to13499999_InternalTransaction",
		"13500000to13749999_InternalTransaction",
		"13750000to13999999_InternalTransaction",
    ];
    for file in files:
        print("file",file)
        i=0
        theZIP = zipfile.ZipFile(fileDir+file+".zip", 'r');
        theCSV = theZIP.open(file+".csv");

        head = theCSV.readline();
        oneLine = theCSV.readline().decode("utf-8").strip();
        while (oneLine!=""):
            if i%1000000==0:
                print(i)
            i+=1
            oneArray = oneLine.split(",")
            transactionHash=oneArray[2]
            toAddr=oneArray[5]
            try:
                value=txMap[transactionHash]
                if "toAddrs" not  in value.keys():
                    value["toAddrs"]={}
                if "toAddrs_internal" not  in value.keys():
                    value["toAddrs_internal"]={}
                value["toAddrs"][toAddr]=1
                value["toAddrs_internal"][toAddr]=1
                txMap[transactionHash]=value
            except:
                pass
            oneLine = theCSV.readline().decode("utf-8").strip();
        theZIP.close()
        
    with open("/mnt/4t3/geth/reorder_geth-dev3/data/analyseInternalTxs/data/dict/erc20/txMap_internal.pkl", "wb") as tf:
        pickle.dump(txMap,tf)
            
# 第三步骤：获取交易对应的外部to地址，交易手续费
# ['blockNumber', 'timestamp', 'transactionHash', 'from', 'to', 'toCreate', 
# 'fromIsContract', 'toIsContract', 'value', 'gasLimit', 'gasPrice', 'gasUsed', 
# 'callingFunction', 'isError', 'eip2718type', 'baseFeePerGas', 'maxFeePerGas', 'maxPriorityFeePerGas']
def getNormalTransactionFromXblock(txMap):
    print("getNormalTransactionFromXblock")
    fileDir = "/mnt/4t2/bowei/sbw/xblock/";

    files = [
		"13000000to13249999_BlockTransaction",
		"13250000to13499999_BlockTransaction",
		"13500000to13749999_BlockTransaction",
		"13750000to13999999_BlockTransaction",
    ];

    for file in files:
        print("file",file)
        theZIP = zipfile.ZipFile(fileDir+file+".zip", 'r');
        theCSV = theZIP.open(file+".csv");

        head = theCSV.readline().decode("utf-8").strip();
        oneLine = theCSV.readline().decode("utf-8").strip();
        # title=head.split(",")
        # print("title",title)
        # return
        i=0
        while (oneLine!=""):
            if i%1000000==0:
                print(i)
            i+=1
            oneArray = oneLine.split(",")

            blockNumber=oneArray[0]
            transactionHash=oneArray[2]
            toAddr=oneArray[4]
            gasPrice=int(oneArray[10])
            gasUsed=int(oneArray[11])
            txFees=gasPrice*gasUsed

            try:       
                value=txMap[transactionHash]
                value["txFees"]=str(txFees)   
                 
                if "toAddrs" not  in value.keys():
                    value["toAddrs"]={}
                if "toAddrs_normal" not  in value.keys():
                    value["toAddrs_normal"]={}
                value["toAddrs"][toAddr]=1
                value["toAddrs_normal"][toAddr]=1
                txMap[transactionHash]=value
            except:
                pass
            oneLine = theCSV.readline().decode("utf-8").strip();
        theZIP.close()
        
    with open("/mnt/4t3/geth/reorder_geth-dev3/data/analyseInternalTxs/data/dict/erc20/txMap_normal.pkl", "wb") as tf:
        pickle.dump(txMap,tf)
        
def concateMapValue(tempMap):
    tempStr=""
    for key in tempMap.keys():
        tempStr+=(key+";")
        
    return tempStr
        
    
# 第四步骤：输出res_txs_map
def outPutMap0(txMap):
    print("outPutMap0")
    global outputCsvName0
    f = open(outputCsvName0,'w')
    writer = csv.writer(f)
    writer.writerow(["transactionHash","numOfAddr","txFees","effectiveFeeOriginal","value2miner","mineIncome","toAddrs","toAddrs_internal","toAddrs_normal"])

    txFees="none"
    try:
        txFees=value["txFees"]
    except:
        pass
        
    for key,value in txMap.items():
        # value["txFees"]可能不存在，由于两个文件的差异导致
        row=[key,value["numOfAddr"],txFees,value["effectiveFeeOriginal"],value["value2miner"],value["mineIncome"],concateMapValue(value["toAddrs"]),concateMapValue(value["toAddrs_internal"]),concateMapValue(value["toAddrs_normal"])]
        writer.writerow(row)

# output: txsDataPlusXblock.csv
# csv type: transactionHash,txFees,numOfAddr,toAddrs 
def getTxData():
    txMap={}
    getTransactionFromROG(txMap)
    # getInterTransactionFromXblock(txMap)
    getNormalTransactionFromXblock(txMap)
    # outPutMap0(txMap)

#····························································
# input: txsDataPlusXblock.csv
# output: addrDataPlusXblock.csv
# csv type: toAddr,numOfTxs,totalTxFees,numOfAddrs,
def getAddrData():
    res_toAddr_map={}
    temp_toAddr_map={}
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
def getToAddrMap(temp_toAddr_Map):
    df = pd.read_csv('../data/csv/erc20/addrDataPlusXblock.csv')
    for index, row in df.iterrows():
        temp_toAddr_Map[row["toAddr"]]={"numOfTxs":row["numOfTxs"],"totalTxFees":row["totalTxFees"],"numOfAddrs":row["numOfAddrs"]}

def getDefi_csv(temp_toAddr_Map):
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


def getDefi_curveJson(temp_toAddr_Map):
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

def getDefi_makerDAOJson(temp_toAddr_Map):
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
            
def getTokenInfo(temp_toAddr_Map):
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
            
def getContractInfo(temp_toAddr_Map):
    dir="/mnt/4t2/bowei/sbw/xblock/VerifiedContract/"
    for fileName in os.listdir(dir):
        filePath=dir+fileName
        with open(filePath) as f:
            jsonFile=json.load(f)
            for key,value in jsonFile.items():
                address=key
                contractName=value["ContractName"]
                try:
                    valueMap=temp_toAddr_Map[address]
                    valueMap["contractName"]=contractName
                except:
                    pass

def getDappName(temp_toAddr_Map):
    path="/mnt/4t2/bowei/sbw/xblock/EthList.txt"
    with open(path, "r", encoding='utf-8') as file:
        for line in file:
            line = line.strip('\n')
            arr=line.split(",")
            dappName=arr[0]
            addrs=arr[1].split("-")

            for addr in addrs:
                try:
                    valueMap=temp_toAddr_Map[addr]
                    valueMap["dappName"]=dappName
                except:
                    pass

def outPutMap1(temp_toAddr_Map):
    global outputCsvName2
    f = open(outputCsvName2,'w')
    writer = csv.writer(f)
    writer.writerow(["toAddr","numOfTxs","totalTxFees","numOfAddrs","defiName","defiType","contractName","symbol","dappName"])

    for key,value in temp_toAddr_Map.items():
        defiName="none"
        defiType="none"
        contractName="none"
        symbol="none"
        dappName="none"
        
        if "defiName" in value.keys():
            defiName=value["defiName"]
        if "defiType" in value.keys():
            defiType=value["defiType"]
        if "contractName" in value.keys():
            contractName=value["contractName"]
        if "symbol" in value.keys():
            symbol=value["symbol"]
        if "dappName" in value.keys():
            dappName=value["dappName"]
            
        row=[key,value["numOfTxs"],value["totalTxFees"],value["numOfAddrs"],defiName,defiType,contractName,symbol,dappName]
        writer.writerow(row)

# input: addrDataPlusXblock.csv
# output: addrDataPlusXblockWithDefi.csv
def plusDefi():
    temp_toAddr_Map={}
    getToAddrMap(temp_toAddr_Map)
    getDefi_csv(temp_toAddr_Map)
    getDefi_curveJson(temp_toAddr_Map)
    getDefi_makerDAOJson(temp_toAddr_Map)
    getTokenInfo(temp_toAddr_Map)
    getContractInfo(temp_toAddr_Map)
    getDappName(temp_toAddr_Map)
    outPutMap1(temp_toAddr_Map)
    
    
#····························································        
def getSpecialAddrData(targetAddrs,temp_toAddr_map):
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
                    

def cropCsv(targetAddr,dir,temp_toAddr_map):
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
    temp_toAddr_map={}
    getSpecialAddrData(targetAddrs,temp_toAddr_map)
    for targetAddr in targetAddrs:
        cropCsv(targetAddr,dir,temp_toAddr_map)
            
    
#····························································
def getAddrToDappName(addr_dappName_map,dappNameList):            
    with open("../data/csv/erc20/addrDataPlusXblockWithDefi.csv") as fOld:
        reader = csv.reader(fOld)
        header_row=next(reader)
        for row in reader:
            toAddr=row[0]
            dappName=row[8]
            if dappName in dappNameList:
                addr_dappName_map[toAddr]=dappName

def getDappNameToTxMap(addr_dappName_map,dappName_tx_map,flashbotsTxMap):
    with open("../data/csv/erc20/txsDataPlusXblock_beifen.csv") as fOld:
        csv.field_size_limit(500 * 1024 * 1024)
        reader = csv.reader(fOld)
        header_row=next(reader)
        index=0
        for row in reader:
            transactionHash=row[0]
            toAddrs=row[3]
            addrsList=toAddrs.split(";")
            for addr in addrsList:
                try:
                    _=flashbotsTxMap[transactionHash]
                    dappName=addr_dappName_map[addr]
                    if dappName in dappName_tx_map.keys():
                        value=dappName_tx_map[dappName]
                        value[transactionHash]=1
                    else:
                        value={transactionHash:1}
                    dappName_tx_map[dappName]=value
                except:
                    pass
            
            index+=1
            if index%100000==0:
                print(index)

def cropCsv1(targetDappName,dir,dappName_tx_map):
    filePath=dir+targetDappName+".csv"
    fNew = open(filePath,'w')
    writer = csv.writer(fNew)
    with open("../data/csv/erc20/addrBalanceDiff_allTokens_syndromicChange_minerIncome.csv") as fOld:
        reader = csv.reader(fOld)
        header_row=next(reader)
        writer.writerow(header_row)
        for row in reader:
            transactionHash=row[3]
            try:
                valueMap=dappName_tx_map[targetDappName]
                if transactionHash in valueMap.keys():
                    writer.writerow(row)
            except:
                pass
    
    
def findTxsRelatedSpecialDappName():
    flashbotsTxMap={}
    file=open("/mnt/4t2/bowei/sbw/xblock/all_blocks").read()
    all_blocks=json.loads(file)
    for block in all_blocks:
        transactions=block["transactions"]
        for transaction in transactions:
            transaction_hash=transaction["transaction_hash"]
            flashbotsTxMap[transaction_hash]=1
            
    addr_dappName_map={}
    dappName_tx_map={}
    # dappNameList=["HEX","UDEX-Mobile-Decentralized-Exchange","Synthetix","Ampleforth","Liquity-Protocol","Fuse","Somnium-Space","Melon","The-Sandbox","Curve","Fulcrum","OpenSea"]
    dappNameList=["OpenSea"]


    print("getAddrToDappName")
    getAddrToDappName(addr_dappName_map,dappNameList)
    print("getDappNameToTxMap")
    getDappNameToTxMap(addr_dappName_map,dappName_tx_map,flashbotsTxMap)
    
    print("output")
    for dappName in dappNameList:
        cropCsv1(dappName,"../data/csv/erc20/specialDappName/",dappName_tx_map)
#····························································
def add_minerIncome_detailedTable():
    tx_value2miner_map={}
    csvPath="/mnt/4t2/bowei/sbw/xblock/transfer2miner.csv"
    df = pd.read_csv(csvPath)
    for index, row in df.iterrows():
        if index%10000==0:
            print(index)
        transactionHash=row["transactionHash"]
        value2miner=row["value2miner"]
        tx_value2miner_map[transactionHash]=value2miner
    
    print("len(tx_value2miner_map) ",len(tx_value2miner_map))
    # return
        
    fNew = open("../data/csv/erc20/addrBalanceDiff_allTokens_syndromicChange_minerIncome.csv",'w')
    writer = csv.writer(fNew)
    with open("../data/csv/erc20/addrBalanceDiff_allTokens_syndromicChange.csv") as fOld:
        reader = csv.reader(fOld)
        header_row=next(reader)
        header_row.insert(8,"value2miner")
        header_row.insert(9,"mineIncome")
        writer.writerow(header_row)
        i=0
        for row in reader:
            if i%10000==0:
                print(i)
            i+=1
            transactionHash=row[3]
            effectiveFeeOriginal=row[7]
            try:
                value2miner=tx_value2miner_map[transactionHash]
                row.insert(8,value2miner)
                mineIncome_int=int(value2miner)+int(effectiveFeeOriginal)
                row.insert(9,str(mineIncome_int))
                writer.writerow(row)
            except:
                row.insert(8,"0")
                row.insert(9,effectiveFeeOriginal)
                writer.writerow(row)

#····························································
def get_onlyFlashbot_detailedTable():
    flashbotsTxMap={}
    file=open("/mnt/4t2/bowei/sbw/xblock/all_blocks").read()
    all_blocks=json.loads(file)
    for block in all_blocks:
        transactions=block["transactions"]
        for transaction in transactions:
            transaction_hash=transaction["transaction_hash"]
            flashbotsTxMap[transaction_hash]=1
            
            
    fNew = open("../data/csv/erc20/addrBalanceDiff_allTokens_syndromicChange_minerIncome_onlyFlashbots.csv",'w')
    writer = csv.writer(fNew)
    with open("../data/csv/erc20/addrBalanceDiff_allTokens_syndromicChange_minerIncome.csv") as fOld:
        reader = csv.reader(fOld)
        header_row=next(reader)
        writer.writerow(header_row)
        i=0
        for row in reader:
            if i%10000==0:
                print(i)
            i+=1
            transactionHash=row[3]
            try:
                _=flashbotsTxMap[transactionHash]
                writer.writerow(row)
            except:
                pass
#····························································
def statistics():
    txMap={}
    tx_value2miner_map={}
    tx_jiaocha_map={}

    with open("../data/csv/erc20/addrBalanceDiff_allTokens_syndromicChange_minerIncome.csv") as f:
        reader = csv.reader(f)
        header_row=next(reader)
        print(header_row)
        
        for row in reader:
            transactionHash=row[3]
            txMap[transactionHash]=1
    
    
    csvPath="/mnt/4t2/bowei/sbw/xblock/transfer2miner.csv"
    df = pd.read_csv(csvPath)
    for index, row in df.iterrows():
        transactionHash=row["transactionHash"]
        tx_value2miner_map[transactionHash]=1
        
        
    for key,value in txMap.items():
        try:
            _=tx_value2miner_map[key]
            tx_jiaocha_map[key]=1
        except:
            pass
    

    print("len(txMap)", len(txMap))
    print("len(tx_value2miner_map) ",len(tx_value2miner_map))
    print("len(tx_jiaocha_map)", len(tx_jiaocha_map))

#····························································
# 1. 出一个交易表，每一行就是详细表的每个交易哈希，有：
# a. 矿工手续费
# b. 矿工转账所得
# c. positionOrigin
# d. positionEdited
# e. 矿工收入即a+b
# f. 账户净增加的ETH总和（即抵消手续费后存在同向变化为正的Ether，WETH)
# g. 账户净增加的USD 总和（即USDT USDC DAI）
# h. 账户净减少的ETH总和
# i. 账户净减少的USD总和


def sumOfList(arr):
    sum=0
    for item in arr:
        sum+=int(item)
    return str(sum)

def getTxsTable():
    nameToTokenAddr={"weth":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","usdt":"0xdac17f958d2ee523a2206206994597c13d831ec7","usdc":"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48","dai":"0x6b175474e89094c44da98b954eedeac495271d0f"}
    txMap={}
    with open("../data/csv/erc20/addrBalanceDiff_allTokens_syndromicChange_minerIncome.csv") as fOld:
        #address,addrType,blockNum,transactionHash,positionOriginal,positionEdited,
        #totalFeeDiff,effectiveFeeOriginal,value2miner,mineIncome
        reader = csv.reader(fOld)
        header_row=next(reader)
        i=0
        for row in reader:
            if i%10000==0:
                print(i)
            i+=1
            blockNum=row[2]
            transactionHash=row[3]
            positionOriginal=row[4]
            positionEdited=row[5]
            effectiveFeeOriginal=row[7]
            value2miner=row[8]
            mineIncome=row[9]
            gasUsedOriginal=row[10]
            if gasUsedOriginal=="21000":
                continue
            
            try:
                txMap_value=txMap[transactionHash]
                index=11
                while(index<len(row)):
                    tokenAddr=row[index]
                    tokenBalanceDiff=row[index+1]
                    if tokenAddr in txMap_value.keys():
                        tempValue=txMap_value[tokenAddr]
                        if int(tokenBalanceDiff)>0:
                            tempValue["pos"]+=(";"+tokenBalanceDiff)
                        else:
                            tempValue["neg"]+=(";"+tokenBalanceDiff)
                    else:
                        tempValue={"pos":"0","neg":"0"}
                        if int(tokenBalanceDiff)>0:
                            tempValue["pos"]=tokenBalanceDiff
                        else:
                            tempValue["neg"]=tokenBalanceDiff

                    txMap_value[tokenAddr]=tempValue
                    index+=2
            except:
                txMap_value={"blockNum":blockNum,"positionOriginal":positionOriginal,"positionEdited":positionEdited,"effectiveFeeOriginal":effectiveFeeOriginal,"value2miner":value2miner,"mineIncome":mineIncome}
                index=11
                while(index<len(row)):
                    tokenAddr=row[index]
                    tokenBalanceDiff=row[index+1]
                    if tokenAddr in txMap_value.keys():
                        tempValue=txMap_value[tokenAddr]
                        if int(tokenBalanceDiff)>0:
                            tempValue["pos"]+=(";"+tokenBalanceDiff)
                        else:
                            tempValue["neg"]+=(";"+tokenBalanceDiff)
                    else:
                        tempValue={"pos":"0","neg":"0"}
                        if int(tokenBalanceDiff)>0:
                            tempValue["pos"]=tokenBalanceDiff
                        else:
                            tempValue["neg"]=tokenBalanceDiff

                    txMap_value[tokenAddr]=tempValue
                    index+=2

            txMap[transactionHash]=txMap_value

    # 计算余额和
    for key,value in txMap.items():
        for tokenAddr in value.keys():
            if tokenAddr =="blockNum" or tokenAddr=="positionOriginal" or tokenAddr=="positionEdited" or tokenAddr=="effectiveFeeOriginal" or tokenAddr=="value2miner" or tokenAddr=="mineIncome":
                continue
            
            sumPos=sumOfList(value[tokenAddr]["pos"].split(";"))
            sumNeg=sumOfList(value[tokenAddr]["neg"].split(";"))

            value[tokenAddr]={"pos":sumPos,"neg":sumNeg}
            txMap[key]=value

    
    # 单独对usdt，usdc，dai求和
    for key,value in txMap.items():
        sumDollarPos=0.0
        sumDollarNeg=0.0
        for tokenAddr in value.keys():
            if tokenAddr =="blockNum" or tokenAddr=="positionOriginal" or tokenAddr=="positionEdited" or tokenAddr=="effectiveFeeOriginal" or tokenAddr=="value2miner" or tokenAddr=="mineIncome":
                continue
            
            if tokenAddr==nameToTokenAddr["usdt"] or tokenAddr==nameToTokenAddr["usdc"]:
                tempInt=float(value[tokenAddr]["pos"])/pow(10,6)
                sumDollarPos+=tempInt
                tempInt=float(value[tokenAddr]["neg"])/pow(10,6)
                sumDollarNeg+=tempInt
                
            if tokenAddr==nameToTokenAddr["dai"]:
                tempInt=float(value[tokenAddr]["pos"])/pow(10,18)
                sumDollarPos+=tempInt
                tempInt=float(value[tokenAddr]["neg"])/pow(10,18)
                sumDollarNeg+=tempInt
        value["dollar"]={"pos":str(sumDollarPos),"neg":str(sumDollarNeg)}
        txMap[key]=value

    # 输出
    print("output")
    f = open("../data/csv/erc20/txsDetail.csv",'w')
    writer = csv.writer(f)
    header_row=["transactionHash","blockNum","effectiveFeeOriginal","value2miner","mineIncome","positionOriginal","positionEdited","etherPos","etherNeg","wethPos","wethNeg","dollarPos","dollarNeg"]
    writer.writerow(header_row)
    
    for key,value in txMap.items():
        effectiveFeeOriginal=value["effectiveFeeOriginal"]
        value2miner=value["value2miner"]
        mineIncome=value["mineIncome"]
        positionOriginal=value["positionOriginal"]
        positionEdited=value["positionEdited"]
        blockNum=value["blockNum"]
        
        etherPos="0"
        etherNeg="0"
        wethPos="0"
        wethNeg="0"
        dollarPos="0"
        dollarNeg="0"
        
        
        if "ether" in value.keys():
            etherPos=value["ether"]["pos"]
            etherNeg=value["ether"]["neg"]
            
        if nameToTokenAddr["weth"] in value.keys():
            wethPos=value[nameToTokenAddr["weth"]]["pos"]
            wethNeg=value[nameToTokenAddr["weth"]]["neg"]
            
        if "dollar" in value.keys():
            dollarPos=value["dollar"]["pos"]
            dollarNeg=value["dollar"]["neg"]

        row=[key,blockNum,effectiveFeeOriginal,value2miner,mineIncome,positionOriginal,positionEdited,etherPos,etherNeg,wethPos,wethNeg,dollarPos,dollarNeg]
        writer.writerow(row)
        
        
def addFlashbots():
    txMap={}
    file=open("/mnt/4t2/bowei/sbw/xblock/all_blocks").read()
    all_blocks=json.loads(file)
    for block in all_blocks:
        transactions=block["transactions"]
        for transaction in transactions:
            transaction_hash=transaction["transaction_hash"]
            txMap[transaction_hash]=1
            
            
    fNew = open("../data/csv/erc20/txsDetail_flashbots.csv",'w')
    writer = csv.writer(fNew)
    with open("../data/csv/erc20/txsDetail.csv") as fOld:
        reader = csv.reader(fOld)
        header_row=next(reader)
        header_row.insert(1,"flashbots")
        writer.writerow(header_row)
        i=0
        for row in reader:
            if i%10000==0:
                print(i)
            i+=1
            transactionHash=row[0]
            try:
                value=txMap[transactionHash]
                row.insert(1,"true")
            except:
                row.insert(1,"false")
                
            writer.writerow(row)
#····························································
def getBlocksTable():
    blockMap={}
    # transactionHash,flashbots,blockNum,effectiveFeeOriginal,value2miner,mineIncome,
    # positionOriginal,positionEdited,etherPos,etherNeg,wethPos,wethNeg,dollarPos,dollarNeg
    with open("../data/csv/erc20/txsDetail_flashbots.csv") as f:
        reader = csv.reader(f)
        header_row=next(reader)
        i=0
        for row in reader:
            blockNum=row[2]
            flashbots=row[1]
            
            try:
                blockMap_value=blockMap[blockNum]
                blockMap_value["txCount"]+=1
                blockMap_value["effectiveFeeOriginal"]+=float(row[3])
                blockMap_value["value2miner"]+=float(row[4])
                blockMap_value["mineIncome"]+=float(row[5])
                blockMap_value["etherPos"]+=float(row[8])
                blockMap_value["etherNeg"]+=float(row[9])
                blockMap_value["wethPos"]+=float(row[10])
                blockMap_value["wethNeg"]+=float(row[11])
                blockMap_value["dollarPos"]+=float(row[12])
                blockMap_value["dollarNeg"]+=float(row[13])
            except:
                blockMap_value={}
                blockMap_value["txCount"]=1
                blockMap_value["effectiveFeeOriginal"]=float(row[3])
                blockMap_value["value2miner"]=float(row[4])
                blockMap_value["mineIncome"]=float(row[5])
                blockMap_value["etherPos"]=float(row[8])
                blockMap_value["etherNeg"]=float(row[9])
                blockMap_value["wethPos"]=float(row[10])
                blockMap_value["wethNeg"]=float(row[11])
                blockMap_value["dollarPos"]=float(row[12])
                blockMap_value["dollarNeg"]=float(row[13])
            blockMap[blockNum]=blockMap_value
    
    # output
    print("output")
    f = open("../data/csv/erc20/blockData_onlyFlashbots.csv",'w')
    writer = csv.writer(f)
    writer.writerow(["blockNum","txCount","effectiveFeeOriginal","value2miner","mineIncome","etherPos","etherNeg","wethPos","wethNeg","dollarPos","dollarNeg"])
    for key,value in blockMap.items():
        writer.writerow([key,str(value["txCount"]),str(value["effectiveFeeOriginal"]),str(value["value2miner"]),str(value["mineIncome"]),str(value["etherPos"]),str(value["etherNeg"]),str(value["wethPos"]),str(value["wethNeg"]),str(value["dollarPos"]),str(value["dollarNeg"])])


def getBlocksTable_OnlyFlashbots():
    blockMap={}
    # transactionHash,flashbots,blockNum,effectiveFeeOriginal,value2miner,mineIncome,
    # positionOriginal,positionEdited,etherPos,etherNeg,wethPos,wethNeg,dollarPos,dollarNeg
    with open("../data/csv/erc20/txsDetail_flashbots.csv") as f:
        reader = csv.reader(f)
        header_row=next(reader)
        i=0
        for row in reader:
            blockNum=row[2]
            flashbots=row[1]
            if flashbots=="false":
                continue
            
            try:
                blockMap_value=blockMap[blockNum]
                blockMap_value["txCount"]+=1
                blockMap_value["effectiveFeeOriginal"]+=float(row[3])
                blockMap_value["value2miner"]+=float(row[4])
                blockMap_value["mineIncome"]+=float(row[5])
                blockMap_value["etherPos"]+=float(row[8])
                blockMap_value["etherNeg"]+=float(row[9])
                blockMap_value["wethPos"]+=float(row[10])
                blockMap_value["wethNeg"]+=float(row[11])
                blockMap_value["dollarPos"]+=float(row[12])
                blockMap_value["dollarNeg"]+=float(row[13])
            except:
                blockMap_value={}
                blockMap_value["txCount"]=1
                blockMap_value["effectiveFeeOriginal"]=float(row[3])
                blockMap_value["value2miner"]=float(row[4])
                blockMap_value["mineIncome"]=float(row[5])
                blockMap_value["etherPos"]=float(row[8])
                blockMap_value["etherNeg"]=float(row[9])
                blockMap_value["wethPos"]=float(row[10])
                blockMap_value["wethNeg"]=float(row[11])
                blockMap_value["dollarPos"]=float(row[12])
                blockMap_value["dollarNeg"]=float(row[13])
            blockMap[blockNum]=blockMap_value
    
    # output
    print("output")
    f = open("../data/csv/erc20/blockData_onlyFlashbots.csv",'w')
    writer = csv.writer(f)
    writer.writerow(["blockNum","txCount","effectiveFeeOriginal","value2miner","mineIncome","etherPos","etherNeg","wethPos","wethNeg","dollarPos","dollarNeg"])
    for key,value in blockMap.items():
        writer.writerow([key,str(value["txCount"]),str(value["effectiveFeeOriginal"]),str(value["value2miner"]),str(value["mineIncome"]),str(value["etherPos"]),str(value["etherNeg"]),str(value["wethPos"]),str(value["wethNeg"]),str(value["dollarPos"]),str(value["dollarNeg"])])

#····························································
def getBlocksTable_fromAddrsTable_onlyFlashbots():    
    blockMap={}
    # with open("../data/csv/erc20/addrsTable_onlyFlashbots.csv") as fOld:
    with open("../data/csv/erc20/opensea/addrsTable_onlyFlashbots.csv") as fOld:
        reader = csv.reader(fOld)
        header_row=next(reader)
        i=0
        for row in reader:
            if i%10000==0:
                print(i)
            i+=1
            
            blockNum=row[2]
            etherPos=row[3]
            etherNeg=row[4]
            wethPos=row[5]
            wethNeg=row[6]
            dollarPos=row[7]
            dollarNeg=row[8]
            
            try:
                blockMap_value=blockMap[blockNum]
                blockMap_value["addrCount"]+=1
                blockMap_value["etherPos"]+=float(etherPos)
                blockMap_value["etherNeg"]+=float(etherNeg)
                blockMap_value["wethPos"]+=float(wethPos)
                blockMap_value["wethNeg"]+=float(wethNeg)
                blockMap_value["dollarPos"]+=float(dollarPos)
                blockMap_value["dollarNeg"]+=float(dollarNeg)
            except:
                blockMap_value={}
                blockMap_value["addrCount"]=1
                blockMap_value["etherPos"]=float(etherPos)
                blockMap_value["etherNeg"]=float(etherNeg)
                blockMap_value["wethPos"]=float(wethPos)
                blockMap_value["wethNeg"]=float(wethNeg)
                blockMap_value["dollarPos"]=float(dollarPos)
                blockMap_value["dollarNeg"]=float(dollarNeg)
                
            blockMap[blockNum]=blockMap_value

    # 输出
    print("output")
    f = open("../data/csv/erc20/opensea/blockData_fromAddrsTable_onlyFlashbots.csv",'w')
    writer = csv.writer(f)
    header_row=["blockNum","addrCount","etherPos","etherNeg","wethPos","wethNeg","dollarPos","dollarNeg"]
    writer.writerow(header_row)
    
    for key,value in blockMap.items():
        blockNum=key
        addrCount=str(value["addrCount"])
        etherPos=str(value["etherPos"])
        etherNeg=str(value["etherNeg"])
        wethPos=str(value["wethPos"])
        wethNeg=str(value["wethNeg"])
        dollarPos=str(value["dollarPos"])
        dollarNeg=str(value["dollarNeg"])
        
        row=[key,addrCount,etherPos,etherNeg,wethPos,wethNeg,dollarPos,dollarNeg]
        writer.writerow(row)
#····························································
def getAddrsTable_onlyFlashbots():
    txMap={}
    file=open("/mnt/4t2/bowei/sbw/xblock/all_blocks").read()
    all_blocks=json.loads(file)
    for block in all_blocks:
        transactions=block["transactions"]
        for transaction in transactions:
            transaction_hash=transaction["transaction_hash"]
            txMap[transaction_hash]=1
            
    nameToTokenAddr={"weth":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","usdt":"0xdac17f958d2ee523a2206206994597c13d831ec7","usdc":"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48","dai":"0x6b175474e89094c44da98b954eedeac495271d0f"}
    addrMap={}
    with open("../data/csv/erc20/addrBalanceDiff_allTokens_syndromicChange_minerIncome.csv") as fOld:
    # with open("/mnt/4t3/geth/reorder_geth-dev3/data/analyseInternalTxs/data/csv/erc20/opensea/OpenSea.csv") as fOld:
        #address,addrType,blockNum,transactionHash,positionOriginal,positionEdited,
        #totalFeeDiff,effectiveFeeOriginal,value2miner,mineIncome
        reader = csv.reader(fOld)
        header_row=next(reader)
        i=0
        for row in reader:
            if i%10000==0:
                print(i)
            i+=1
            address=row[0]
            addrType=row[1]
            blockNum=row[2]
            transactionHash=row[3]
            gasUsedOriginal=row[10]
            if gasUsedOriginal=="21000":
                continue
            
            if transactionHash not in txMap.keys():
                continue
            
            newkey=address+"_"+blockNum
            try:
                addrMap_value=addrMap[newkey]
                index=11
                while(index<len(row)):
                    tokenAddr=row[index]
                    tokenBalanceDiff=row[index+1]
                    if tokenAddr in addrMap_value.keys():
                        tempValue=addrMap_value[tokenAddr]
                        if int(tokenBalanceDiff)>0:
                            tempValue["pos"]+=(";"+tokenBalanceDiff)
                        else:
                            tempValue["neg"]+=(";"+tokenBalanceDiff)
                    else:
                        tempValue={"pos":"0","neg":"0"}
                        if int(tokenBalanceDiff)>0:
                            tempValue["pos"]=tokenBalanceDiff
                        else:
                            tempValue["neg"]=tokenBalanceDiff

                    addrMap_value[tokenAddr]=tempValue
                    index+=2
            except:
                addrMap_value={"address":address,"addrType":addrType,"blockNum":blockNum}
                index=11
                while(index<len(row)):
                    tokenAddr=row[index]
                    tokenBalanceDiff=row[index+1]
                    if tokenAddr in addrMap_value.keys():
                        tempValue=addrMap_value[tokenAddr]
                        if int(tokenBalanceDiff)>0:
                            tempValue["pos"]+=(";"+tokenBalanceDiff)
                        else:
                            tempValue["neg"]+=(";"+tokenBalanceDiff)
                    else:
                        tempValue={"pos":"0","neg":"0"}
                        if int(tokenBalanceDiff)>0:
                            tempValue["pos"]=tokenBalanceDiff
                        else:
                            tempValue["neg"]=tokenBalanceDiff

                    addrMap_value[tokenAddr]=tempValue
                    index+=2

            addrMap[newkey]=addrMap_value

    # 计算余额和
    for key,value in addrMap.items():
        for tokenAddr in value.keys():
            if tokenAddr =="address" or tokenAddr=="addrType" or tokenAddr=="blockNum":
                continue
            
            sumPos=sumOfList(value[tokenAddr]["pos"].split(";"))
            sumNeg=sumOfList(value[tokenAddr]["neg"].split(";"))

            value[tokenAddr]={"pos":sumPos,"neg":sumNeg}
            addrMap[key]=value

    
    # 单独对usdt，usdc，dai求和
    for key,value in addrMap.items():
        sumDollarPos=0.0
        sumDollarNeg=0.0
        for tokenAddr in value.keys():
            if tokenAddr =="address" or tokenAddr=="addrType" or tokenAddr=="blockNum":
                continue
            
            if tokenAddr==nameToTokenAddr["usdt"] or tokenAddr==nameToTokenAddr["usdc"]:
                tempInt=float(value[tokenAddr]["pos"])/pow(10,6)
                sumDollarPos+=tempInt
                tempInt=float(value[tokenAddr]["neg"])/pow(10,6)
                sumDollarNeg+=tempInt
                
            if tokenAddr==nameToTokenAddr["dai"]:
                tempInt=float(value[tokenAddr]["pos"])/pow(10,18)
                sumDollarPos+=tempInt
                tempInt=float(value[tokenAddr]["neg"])/pow(10,18)
                sumDollarNeg+=tempInt
        value["dollar"]={"pos":str(sumDollarPos),"neg":str(sumDollarNeg)}
        addrMap[key]=value

    # 输出
    print("output")
    f = open("../data/csv/erc20/addrsTable_onlyFlashbots.csv",'w')
    writer = csv.writer(f)
    header_row=["address","addrType","blockNum","etherPos","etherNeg","wethPos","wethNeg","dollarPos","dollarNeg"]
    writer.writerow(header_row)
    
    for key,value in addrMap.items():
        address=value["address"]
        addrType=value["addrType"]
        blockNum=value["blockNum"]
        
        etherPos="0"
        etherNeg="0"
        wethPos="0"
        wethNeg="0"
        dollarPos="0"
        dollarNeg="0"
        
        
        if "ether" in value.keys():
            etherPos=value["ether"]["pos"]
            etherNeg=value["ether"]["neg"]
            
        if nameToTokenAddr["weth"] in value.keys():
            wethPos=value[nameToTokenAddr["weth"]]["pos"]
            wethNeg=value[nameToTokenAddr["weth"]]["neg"]
            
        if "dollar" in value.keys():
            dollarPos=value["dollar"]["pos"]
            dollarNeg=value["dollar"]["neg"]

        row=[address,addrType,blockNum,etherPos,etherNeg,wethPos,wethNeg,dollarPos,dollarNeg]
        writer.writerow(row)
#····························································
# % MEV-Bots即外部交易的to
# % 横向柱状图：y轴是Top10的地址，x轴为矿工手续费所得+矿工内部转账所得，

# % Contract即内外部交易的to
# % 横向柱状图：图a,y轴是Top10的地址，x轴为矿工手续费所得+矿工内部转账所得；
# %           图b,y轴地址相同，x轴为交易数量

# % DApps即to对应的DApp作进一步归集（不能生硬的把to的交易加起来，因为不同的to可能会对应同一笔交易）
# % 横向柱状图：y轴改为DApp，这里估计Top5就够
# # 得到to地址对应交易的手续费和，私发的总和

#····························································

outputCsvName3="../data/csv/erc20/txs_addrCount.csv"
outputCsvName0="../data/csv/erc20/txsDataPlusXblock.csv"
outputCsvName1="../data/csv/erc20/addrDataPlusXblock.csv"
outputCsvName2="../data/csv/erc20/addrDataPlusXblockWithDefi.csv"

def main():
    # print("getTxData")
    # getTxData()

    # print("getAddrData")
    # getAddrData()

    # print("plusDefi")
    # plusDefi()
    
    # findTxsRelatedSpecialDappName()
    
    #·······················································
    # add_minerIncome_detailedTable()
    get_onlyFlashbot_detailedTable()
    #·······················································
    
    # statistics()
    # getTxsTable()
    # addFlashbots()
    # getBlocksTable()
    # getBlocksTable_OnlyFlashbots()
    
    # getAddrsTable_onlyFlashbots()
    getBlocksTable_fromAddrsTable_onlyFlashbots()
    

if __name__ == '__main__':
    main()
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
    with open("../data/csv/erc20/getUniV3.csv") as fOld:
        reader = csv.reader(fOld)
        # header_row=next(reader)
        i=0
        for row in reader:
            if i%10000==0:
                print(i)
            i+=1
            # print("len(row)",len(row))
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
    f = open("../data/csv/erc20/uniswap/txsDetail_flashbots.csv",'w')
    writer = csv.writer(f)
    header_row=["transactionHash","flashbots","blockNum","effectiveFeeOriginal","value2miner","mineIncome","positionOriginal","positionEdited","etherPos","etherNeg","wethPos","wethNeg","dollarPos","dollarNeg"]
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

        row=[key,"true",blockNum,effectiveFeeOriginal,value2miner,mineIncome,positionOriginal,positionEdited,etherPos,etherNeg,wethPos,wethNeg,dollarPos,dollarNeg]
        writer.writerow(row)
        
#····························································
def getBlocksTable():
    blockMap={}
    # transactionHash,flashbots,blockNum,effectiveFeeOriginal,value2miner,mineIncome,
    # positionOriginal,positionEdited,etherPos,etherNeg,wethPos,wethNeg,dollarPos,dollarNeg
    with open("../data/csv/erc20/uniswap/txsDetail_flashbots.csv") as f:
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
    f = open("../data/csv/erc20/uniswap/blockData_onlyFlashbots.csv",'w')
    writer = csv.writer(f)
    writer.writerow(["blockNum","txCount","effectiveFeeOriginal","value2miner","mineIncome","etherPos","etherNeg","wethPos","wethNeg","dollarPos","dollarNeg"])
    for key,value in blockMap.items():
        writer.writerow([key,str(value["txCount"]),str(value["effectiveFeeOriginal"]),str(value["value2miner"]),str(value["mineIncome"]),str(value["etherPos"]),str(value["etherNeg"]),str(value["wethPos"]),str(value["wethNeg"]),str(value["dollarPos"]),str(value["dollarNeg"])])

#····························································
def getBlocksTable_fromAddrsTable_onlyFlashbots():    
    blockMap={}
    # with open("../data/csv/erc20/addrsTable_onlyFlashbots.csv") as fOld:
    with open("../data/csv/erc20/uniswap/addrsTable_onlyFlashbots.csv") as fOld:
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
    f = open("../data/csv/erc20/uniswap/blockData_fromAddrsTable_onlyFlashbots.csv",'w')
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
    nameToTokenAddr={"weth":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","usdt":"0xdac17f958d2ee523a2206206994597c13d831ec7","usdc":"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48","dai":"0x6b175474e89094c44da98b954eedeac495271d0f"}
    addrMap={}
    with open("../data/csv/erc20/getUniV3.csv") as fOld:
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
    f = open("../data/csv/erc20/uniswap/addrsTable_onlyFlashbots.csv",'w')
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
    
    # add_minerIncome_detailedTable()
    # statistics()
    
    
    getTxsTable()
    getBlocksTable()
    # getBlocksTable_OnlyFlashbots()
    
    getAddrsTable_onlyFlashbots()
    getBlocksTable_fromAddrsTable_onlyFlashbots()
    

if __name__ == '__main__':
    main()
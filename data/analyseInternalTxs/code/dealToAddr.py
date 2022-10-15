
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

def output(tempMap,name):
    file="/mnt/4t3/geth/reorder_geth-dev3/data/analyseInternalTxs/data/csv/erc20/to/"+name+".csv"
    f = open(file,'w')
    writer = csv.writer(f)
    writer.writerow(["address","mineIncome"])
        
    for key,value in tempMap.items():
        # value["txFees"]可能不存在，由于两个文件的差异导致
        row=[key,value["mineIncome"]]
        writer.writerow(row)


def getTopToaddr(flashbotsTxMap):
    txsMap={}
    with open("/mnt/4t3/geth/reorder_geth-dev3/data/analyseInternalTxs/data/dict/erc20/txMap_normal_complete.pkl", "rb") as tf:
        txsMap=pickle.load(tf)
    toInterMap={}
    toNormalMap={}
    toMap={}
    for key,value in txsMap.items():
        effectiveFeeOriginal=float(value["effectiveFeeOriginal"])
        value2miner=float(value["value2miner"])
        mineIncome=float(value["mineIncome"])
        try:
            _=flashbotsTxMap[key]
            if "toAddrs_internal" in value.keys():
                for toAddr_internal in value["toAddrs_internal"].keys():
                    if toAddr_internal in toInterMap.keys():
                        toInterMap[toAddr_internal]["effectiveFeeOriginal"]+=effectiveFeeOriginal
                        toInterMap[toAddr_internal]["value2miner"]+=value2miner
                        toInterMap[toAddr_internal]["mineIncome"]+=mineIncome
                    else:
                        toInterMap[toAddr_internal]={"effectiveFeeOriginal":effectiveFeeOriginal,"value2miner":value2miner,"mineIncome":mineIncome}
            
            if "toAddrs_normal" in value.keys():
                for toAddr_normal in value["toAddrs_normal"].keys():
                    if toAddr_normal in toNormalMap.keys():
                        toNormalMap[toAddr_normal]["effectiveFeeOriginal"]+=effectiveFeeOriginal
                        toNormalMap[toAddr_normal]["value2miner"]+=value2miner
                        toNormalMap[toAddr_normal]["mineIncome"]+=mineIncome
                    else:
                        toNormalMap[toAddr_normal]={"effectiveFeeOriginal":effectiveFeeOriginal,"value2miner":value2miner,"mineIncome":mineIncome}
            
            if "toAddrs_internal" in value.keys():
                for toAddr in value["toAddrs_internal"].keys():
                    if toAddr in toMap.keys():
                        toMap[toAddr]["effectiveFeeOriginal"]+=effectiveFeeOriginal
                        toMap[toAddr]["value2miner"]+=value2miner
                        toMap[toAddr]["mineIncome"]+=mineIncome
                    else:
                        toMap[toAddr]={"effectiveFeeOriginal":effectiveFeeOriginal,"value2miner":value2miner,"mineIncome":mineIncome}
        except:
            pass
        
    print("output")
    output(toInterMap,"toInterMap")
    output(toNormalMap,"toNormalMap")
    output(toMap,"toMap")


def main():
    flashbotsTxMap={}
    file=open("/mnt/4t2/bowei/sbw/xblock/all_blocks").read()
    all_blocks=json.loads(file)
    for block in all_blocks:
        transactions=block["transactions"]
        for transaction in transactions:
            transaction_hash=transaction["transaction_hash"]
            flashbotsTxMap[transaction_hash]=1
            
    print("getTopToaddr")
    getTopToaddr(flashbotsTxMap)
    
if __name__ == '__main__':
    main()
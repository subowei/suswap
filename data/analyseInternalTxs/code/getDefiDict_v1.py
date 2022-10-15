# coding: utf-8

import zipfile
import json
import time
import os
import sys
import pickle 
import pandas as pd
import numpy as np
import json
from pandas import Series, DataFrame

keyStr="Transaction.csv"

def get_allfilesPath(fileDir):  # 获取训练文件的文件名
    all_file = []
    for fileName in os.listdir(fileDir):  #listdir返回文件中所有目录
        if keyStr not in fileName:
            continue
        
        filePath=fileDir+fileName

        all_file.append(filePath)
        
    return all_file


def getTxsFromCsv(path):
    txsMap={}
    df = pd.read_csv(path)
    for index, row in df.iterrows():
        transactionHash=row["transactionHash"]
        txsMap[transactionHash]=1
    return txsMap


def main():

    fileDir = "/mnt/4t2/bowei/sbw/xblock/peilin_defi/"
    filesPath=get_allfilesPath(fileDir)
    print("len:",len(filesPath))

    txsMap=getTxsFromCsv("../data/csv/addrBalanceDiff_specialTokens_relatedTxs.csv")

    defiMap={}

    a=0
    for filePath in filesPath:
        fileName=filePath.split('/')[-1].split('.')[0]
        defiName=fileName.split('_')[0]

        df = pd.read_csv(filePath)
        print(a)
        a+=1
        for index, row in df.iterrows():
            print("row",row)
            transactionHash=row["transactionHash"]
            try:
                _=txsMap[transactionHash]
                tempType="None"
                if "type" in row.index:
                    tempType=row["type"]
                try:
                    valueOfDict=defiMap[transactionHash]
                    valueOfDict[defiName]=tempType
                except:
                    valueOfDict={}
                    valueOfDict[defiName]=tempType
                    
                defiMap[transactionHash]=valueOfDict
            except:
                continue

        with open('../data/dict/defiPartMap'+str(a)+'.data','wb') as f:
            pickle.dump(defiMap,f)  



if __name__ == '__main__':
    main()
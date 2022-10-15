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
import csv


# keyStr="addrBalanceDiff_specialTokens_sortedBy"
keyStr="addrBalanceDiff_allTokens_sortedByETHER"

def get_allfilesPath(fileDir):  # 获取训练文件的文件名
    all_file = []
    for fileName in os.listdir(fileDir):  #listdir返回文件中所有目录
        if keyStr not in fileName:
            continue
        
        filePath=fileDir+fileName

        all_file.append(filePath)
        
    return all_file


def main():
    # defi信息
    defiMap={}
    with open('../data/dict/defiPartMap13.data','rb') as f:
        defiMap=pickle.load(f)  

    # 原始数据
    fileDir = "/home/dell/sbw/reorder_geth/data/analyseInternalTxs/data/csv/"
    filesPath=get_allfilesPath(fileDir)
    print("len:",len(filesPath))

    a=0
    for filePath in filesPath:
        print(filePath)
        print(a)
        a+=1

        df = pd.read_csv(filePath)

        newFilePath=filePath.split(".")[0]+"_peilinDefi.csv"
        f = open(newFilePath,'w')
        writer = csv.writer(f)
        tempColumns=list(df.columns)
        tempColumns.extend(["defiName","defiType"])
        writer.writerow(tempColumns)

        for index, row in df.iterrows():
            transactionHash=row["transactionHash"]
            try:
                defiMap_value=defiMap[transactionHash]
                for k,v in defiMap_value.items():
                    row=row.append( pd.Series([k,v],index=["defiType","transactionType"]))
                    # print(row)
                    # return
            except:
                pass

            writer.writerow(row)

        


if __name__ == '__main__':
    main()
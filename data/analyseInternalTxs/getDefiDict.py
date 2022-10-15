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


fileDir = "/mnt/g/geth_data/xblock/defi/";

files = [
    "AaveV1_Transaction.csv",
    "BalancerV1_Transaction.csv",
    "Compound_Transaction.csv",
    "Curve_Transaction.csv",
    "MakerDAO_Transaction.csv",
    "UniswapV1_Transaction.csv"
];

defiDict={}

a=0
for file in files:
    defiName=file.split('_')[0]
    df = pd.read_csv(fileDir+defiName+"_eth/"+file)
    print(a)
    a+=1
    for index, row in df.iterrows():
        transactionHash=row["transactionHash"]
        tempType="None"
        if "type" in row.index:
            tempType=row["type"]
        if defiDict.get(transactionHash)!=None:
            valueOfDict=defiDict.get(transactionHash)
            valueOfDict.extend([defiName,tempType])
        else:
            valueOfDict=[defiName,tempType]
            
        defiDict[transactionHash]=valueOfDict
        
with open('../data/dict/defiDict.data','wb') as f:
    pickle.dump(defiDict,f)  
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

oriDf = pd.read_csv('../data/csv/gethInternalTxs.csv')

tempColumns=list(oriDf.columns)
tempColumns.extend(["defiType","transactionType"])

defiDictV1={}
with open('../data/dict/defiDict.data','rb') as f:
    defiDictV1=pickle.load(f)
        
f = open('../data/csv/gethInternalTxsWithDefi.csv','w')
writer = csv.writer(f)
writer.writerow(tempColumns)
a=0
for index, row in oriDf.iterrows():
    if a%10000==0:
        print(a)
    a+=1
    transactionHash=row["transactionHash"]
    

    if defiDictV1.get(transactionHash)!=None:
        dictValue=defiDictV1[transactionHash]

        defiTypeString=""
        transactionTypeString=""

        for i in range(len(dictValue)):
            if i%2==0:   
                defiTypeString+=dictValue[i]
                if i!=len(dictValue)-2:
                    defiTypeString+=','
            else:
                transactionTypeString+=dictValue[i]
                if i!=len(dictValue)-1:
                    transactionTypeString+=','
            
        row=row.append( pd.Series([defiTypeString,transactionTypeString], index=["defiType","transactionType"]))
        writer.writerow(row)
    else:
        row=row.append( pd.Series(["None","None"], index=["defiType","transactionType"]) )
        writer.writerow(row)
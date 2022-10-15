import pandas as pd
import numpy as np
import csv
import zipfile
import pickle


txsDf=pd.read_csv("../data/csv/txs.csv")

# txsList=[]
txsMap={}
for index, row in txsDf.iterrows():
    # txsList.append(row["transaction"])
    txsMap[row["transaction"]]=1

fileDir = "/mnt1/sbw/xblock/";

files = [
    "11000000to11999999_InternalTransaction",
    "12000000to12999999_InternalTransaction",
    "11000000to11999999_BlockTransaction",
    "12000000to12999999_BlockTransaction"
];

addrMap={}
for file in files:
    if file.find("InternalTransaction")!=-1:
        theZIP = zipfile.ZipFile(fileDir+"internalTransaction/"+file+".zip", 'r');
    else:
        theZIP = zipfile.ZipFile(fileDir+"normalTransaction/"+file+".zip", 'r');

#     theZIP = zipfile.ZipFile(fileDir+file+".zip", 'r');
    theCSV = theZIP.open(file+".csv");

    head = theCSV.readline();

    oneLine = theCSV.readline().decode("utf-8").strip();
    
    isInternalTransaction=True
    if file.find("InternalTransaction")==-1:
        isInternalTransaction=False
    
    a=0
    while (oneLine!=""):
        if a%100000==0:
            print(a)
        a+=1
        oneArray = oneLine.split(",");
        transaction=oneArray[2]
        
        try:
            value=txsMap[transaction]

            if isInternalTransaction==False:
                to=oneArray[4]
            else:
                to=oneArray[5]
                
            if to in addrMap.keys():
                addrMap[to]+=1
            else:
                addrMap[to]=1
                
            oneLine = theCSV.readline().decode("utf-8").strip()
        except:
            oneLine = theCSV.readline().decode("utf-8").strip()
            continue
        
    with open("../data/dict/toAddrsAndCount_"+file+".pkl", "wb") as tf:
        pickle.dump(addrMap,tf)       
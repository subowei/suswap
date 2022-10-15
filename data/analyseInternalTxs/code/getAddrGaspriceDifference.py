import pandas as pd
import numpy as np
import pickle


def getDictFromDf(df):
    tempDict={}
    for index, row in df.iterrows():
        fromAddr=row["fromAddr"]
        toAddr=row["toAddr"]
        
        tempDict[fromAddr]=[row,"fromAddr"]
        tempDict[toAddr]=[row,"toAddr"]
    return tempDict
    
def getGaspriceDifference(oriDict,editedDict):
    df=pd.DataFrame(columns=('address','balanceDifference','blockNum','editedTransaction','oriTransaction',"serviceChargeDiff"))
    for (key,value) in oriDict.items():
        oriRow=value[0]
        oriSymbol=value[1]
        oriPrice=float(0)
        
        editedRow=editedDict[key][0]
        editedSymbol=editedDict[key][1]
        editedPrice=float(0)
        
        if oriSymbol=="fromAddr":
            oriPrice=float(oriRow["fromBalanceOriginal"])
        else:
            oriPrice=float(oriRow["toBalanceOriginal"])
            
        if editedSymbol=="fromAddr":
            editedPrice=float(editedRow["fromBalanceEdited"])
        else:
            editedPrice=float(editedRow["toBalanceEdited"])
        
        balanceDifference=oriPrice-editedPrice
        serviceChargeDiff=oriRow["serviceChargeOriginal"]-oriRow["serviceChargeEdited"]
        
        df=df.append([{'address':key,'balanceDifference':balanceDifference,'blockNum':oriRow["blockNum"],'editedTransaction':editedRow["transactionHash"],'oriTransaction':oriRow["transactionHash"],'serviceChargeDiff':serviceChargeDiff}],ignore_index = True)
        
    return df
    
if __name__ == '__main__':
    oriOrderDf=pd.read_csv('../data/csv/gethInternalTxs_oriOrder_v4.csv')
    editedOrderDF=pd.read_csv('../data/csv/gethInternalTxs_editedOrder_v4.csv')

    lastBlockNum=0
    currentBlockNum=0
    fromIndex=0
    endIndex=0

    resultDf=pd.DataFrame(columns=('address','balanceDifference','blockNum','editedTransaction','oriTransaction','serviceChargeDiff'))

    # 临时存储同一个区块中的交易数据
    oriDict={}
    editedDict={}
    a=0
    for index, row in oriOrderDf.iterrows():
        if a%1000==0:
            print(a)
        a+=1
        
        #收集同一个区块的交易的索引号
        currentBlockNum=row["blockNum"]
        endIndex=index
        
        # print('fromIndex',fromIndex)
        # print('endIndex',endIndex)
        
        # 如果当前区块号与上一不相同，说明应该对前一个区块中的进行处理
        if(currentBlockNum!=lastBlockNum):
            # 如果区块中的交易个数小于等于1,则不予处理。因为只有一条交易那么余额差肯定为0
            if endIndex-fromIndex>1:

                oriDict=getDictFromDf(oriOrderDf[fromIndex:endIndex])
                editedDict=getDictFromDf(editedOrderDF[fromIndex:endIndex])
                
                resultDf = resultDf.append( getGaspriceDifference(oriDict,editedDict) )
            
            fromIndex=index
            
        lastBlockNum=currentBlockNum
        
    resultDf.to_csv('../data/csv/addrBalanceDifference.csv',index=0)
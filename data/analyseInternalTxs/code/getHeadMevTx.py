
inputF = open("/mnt/4t3/geth/reorder_geth-dev3/data/analyseInternalTxs/data/csv/erc20/addrBalanceDiff_allTokens_syndromicChange.csv")

outputF = open("/mnt/4t3/geth/reorder_geth-dev3/data/analyseInternalTxs/data/csv/erc20_underSpecialPosition/getHeadMevTx.csv", "w")


head = inputF.readline().strip()
outputF.write(head+"\n")
print(head)
# address,addrType,blockNum,transactionHash,positionOriginal,positionEdited,totalFeeDiff
line = inputF.readline().strip()
while line != "":
    arr = line.strip().split(",")
    positionOriginal = int(arr[4])
    if positionOriginal < 5:
        # print(positionOriginal)
        outputF.write(line+"\n")


    line = inputF.readline().strip()

inputF.close()
outputF.close()
getDefiDict.py
输入：沛哥的defi文件
输出：defiDict.data     对应的dict文件


dictToCsv.go:
输入：geth原始结果，dict
输出：gethInternalTxs.csv

dictToCsv_sort.go:
输入：geth原始结果，dict
输出：
gethInternalTxs_oriOrder.csv：按照交易位置从小到大排序，也就是原始的交易顺序
gethInternalTxs_editedOrder.csv：按照gasprice从大到小排序

对上面结果统计交易地址的余额差
处理文件：getAddrGaspriceDifference.py
gethInternalTxs_oriOrder.csv，gethInternalTxs_editedOrder.csv
输出：addrBalanceDifference.csv

统计每个区块的收益
处理文件：getBlockProfit.ipynb
输入：addrBalanceDifference.csv


pushDefiForGethRes.py:
作用：将gethInternalTxs.csv加上沛哥defi数据
输入：gethInternalTxs.csv
输出：gethInternalTxsWithDefi.csv


// defi···················································
getDefiDict_v1.py
输入：沛哥的defi文件(更新版本),addrBalanceDiff_specialTokens_relatedTxs.csv
原因：如果统计完整的数据太大，所以只统计给点交易哈希的数据
输出：defiPartMap.data     对应的dict文件,eg:defiPartMap13
dict格式：{transactionHash:{defiName:functionName}}


// eth····················································
getProfitFromAddrMap.go
作用：
1. 处理geth输出的addrmap，得到每个区块的收益（应该为0），用于验证重排序是否正确
2. 将结果按照balanceDifference排序,func getAddrData(),func storeAddrDataAfterSorting()
3. 获得得到键值对{交易hash：1}，该交易hash相关的地址存在余额差，func getTransactionWithBalanceDifference()
输入：mapAddr
输出：blockProfitFromAddrMap.csv,addrBalanceDiff_sortedByDiff.csv,txs.csv

getProfitFromAddrMap.ipynb
作用：对getProfitFromAddrMap的输出进行展示
输入：blockProfitFromAddrMap.csv

getNumOfAddrWithBDFromXblock.ipynb
getNumOfAddrWithBDFromXblock.py
作用：使用txs.csv和xblock，交叉，得到该交易相关的所有地址，并进行统计
输入：txs.csv，xblock内部和外部数据
输出：toAddrsAndCount.csv


getProfitFromAddrMap_v1.go
函数：getData0
作用：得到txHash, 外部和内部交易的to地址, 交易手续费
输出：data0.csv


// erc20····················································
getProfitFromAddrMapERC20.go
作用：处理geth输出的addrmapERC20，得到每个区块的收益
输入：mapAddr
输出：blockProfitFromAddrMap.csv, addrMapERC20Res.csv

getAddrBalanceDiff():某个地址在某个区块中涉及所有代币的余额差
getAddrBalanceDiff_specialTokens():
    作用：某个地址在某个区块中几种代币的余额差
    输出：
        1.某个地址在特定token上面的余额差   addrBalanceDiff_specialTokens_sortedByETHER.csv
        2.上述涉及的所有交易哈希    addrBalanceDiff_specialTokens_relatedTxs.csv


analyseERC20.ipynb
作用：对上述结果进行统计


getProfitFromAddrMapERC20WithDefi.py
作用：将上述的信息加上沛哥的defi信息
输入：addrBalanceDiff_specialTokens_sortedByETHER.csv,  defiPartMap.data
输出：addrBalanceDiff_specialTokens_sortedByETHER_peilinDefi.csv

getProfitFromAddrMapERC20WithDefi_irregular.py
作用：同上，但是有些csv是不规则的，比如一个地址的多个代币发生变化，此时的csv无法通过pd来读取
输入：addrBalanceDiff_allTokens_sortedByETHER.csv,  defiPartMap.data
····················································

getProfitFromAddrMapERC721.go
作用：处理geth输出的addrmapERC721，统计nft抢跑的现象，统计nft套利现象
输入：mapAddr
输出：addrMapERC721Res.csv, arbitrageERC721Txs.csv


nftArbitrage.go
作用：使用完整的event记录统计套利，更准确
输入：mapEvent
输出：

// 交易····················································
getTransactions.go
作用：统计具体的交易
输入：mapTxs的内容

analyseTxs.ipynb

findTransactionsWithUnreasonableOrder.go
作用：找出指定区块中，顺序不合理的交易


// 混合分析1····················································2022.9.28
文件：getProfitFromAddrMapERC20.go
函数：getAddrBalanceDiff_specialTokens_afterSorting_all()
输入：output_addrInfo_erc20
输出：addrBalanceDiff_allTokens_syndromicChange.csv

文件：xblockTx_erc20.py
使用：getTxData(), getAddrData(), plusDefi()
输入：xblock内外部交易数据，addrBalanceDiff_allTokens_syndromicChange.csv
输出：txHash, 外部和内部交易的to地址, 交易手续费，余额发生变化的地址总数。addrDataPlusXblockWithDefi.csv

// 混合分析2····················································
文件：xblockTx_erc20.py
使用：findTxsRelatedSpecialAddr()
输入：txsDataPlusXblock.csv, addrBalanceDiff_allTokens_syndromicChange.csv
输出：指定地址所涉及的交易哈希      /data/addrs_txs/?.txt
作用：在addrBalanceDiff_allTokens_syndromicChange筛选指定地址涉及的交易的条目输出


// 混合分析1····················································2022.10.2
文件：getProfitFromAddrMapERC721_v2.go
函数：getTxs0
输入：output_addrInfo_erc20
输出：txs_addrNum_Map.csv（交易哈希，nft归属权发生变化的地址）

文件：xblockTx_erc721.py
使用：按照顺序执行main中三个函数
输入：xblock内外部交易数据，txs_addrNum_Map.csv
输出：txHash, 外部和内部交易的to地址, 交易手续费，余额发生变化的地址总数



函数，事件选择器：
erc20:
网络版：
    event transfer: 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
    func balanceOf: 0x70a08231
官方版本：
    event transfer: 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
    func balanceOf: 0x70a08231

erc721:
官方版本：
    event transfer: 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
    func ownerOf: 0x6352211e


余额:
地址：

defi:
WETH USDT USDC WBTC DA
套利交易
    先找哪些人赚钱多，和少

    三明治。具体信息：{地址：赚了多少钱，什么类型}
    清算：多个人抢着买
    


nft：
地址：改变tokenid的数量
    套利
    矿工抢发交易

    
scp:
scp -r    root@172.18.197.97:/home/dell/sbw/gethdir2/  ./1000beifen

scp -r    root@192.168.0.100:/home/dell/sbw/gethv7/  ./


scp -r    root@172.18.197.98:/mnt/4t2/bowei/sbw/gethdir  ./gethdir_light

./build/bin/geth --http --http.api web3,eth,net,personal --http.port 9876 --syncmode "light"  --datadir /mnt/4t2/bowei/sbw/gethdir/

./build/bin/geth attach http://localhost:9876
eth.syncing

//完整同步
./build/bin/geth --http --http.api web3,eth,net,personal --http.port 9876  --datadir /mnt/4t2/bowei/sbw/gethdir/


wget http://172.18.197.86:3389/down/QLJFXHagyA5s

wget  http://172.18.197.86:3389/down/U12KlIF0XfNc?fname=/UniswapV2_eth.zip
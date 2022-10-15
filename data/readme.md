

导出：
./build/bin/geth --datadir "/home/dell/sbw/gethdir2" export /home/dell/sbw/exported/0To1000wBySbw.rlp 0 10000000


geth文件夹被锁：

ps -aux | grep geth

killall -HUP geth

然后把它kill -9


//``````````````````````````````````````````````````````````````````````````````````````````````````````

https://gitee.com/subway2020/reorder_geth.git
git remote add origin https://github.com/subowei/reorder_geth.git
git push -u origin dev2


git remote add origin2 https://gitee.com/subway2020/reorder_geth.git
git push -u origin2 master


·······································
geth设置dir，直接设置不行，所以使用软连接
"/home/sbw/.ethereum"
ln -s "/mnt/a/geth/gethdir/geth"  "/home/sbw/.ethereum/geth" 

ln -s "/mnt/a/geth/gethdir1/geth"  "/home/sbw/.ethereum/geth" 

./build/bin/geth --http --http.port 9876  --http.api web3,eth,net,personal


./build/bin/geth attach http://localhost:9876

./build/bin/geth attach http://localhost:30303

eth.syncing

·······································
geth同步：

./build/bin/geth --http  --http.api web3,eth,net,personal --http.port 9876 --syncmode=full --snapshot=false --datadir=/mnt/a/geth/gethdir2


./build/bin/geth --http  --http.api web3,eth,net,personal --http.port 30303  --snapshot=false 

./build/bin/geth --http  --http.api web3,eth,net,personal --http.port 9876 --syncmode "full"

./build/bin/geth --http  --http.api web3,eth,net,personal --http.port 8545 --syncmode "snap" --cache 14436

./build/bin/geth --syncmode full

·······································
export:

./build/bin/geth --datadir "/home/sbw/.ethereum"  export /mnt/a/geth/exportDir/1400w.block 0 1400000

·······································
geth离线模式：

./build/bin/geth --http  --http.api web3,eth,net,personal --http.port 9876 --nodiscover  

·······································
eth命令：

eth.syncing
eth.getTransaction("0x5a4bf6970980a9381e6d6c78d96ab278035bbff58c383ffe96a0a2bbc7c02a4b")
eth.getBlock("14725000",true)  //加上true会显示交易细节

eth.getBlock("14732315",true) 

eth.getBlock("latest").number

0x838025a84c1b64b1d82f0102f6029fbc6c59c8d3
0xe2abfc3f7c1e8be80bc2af509ace5f7a58660be7




杀死进程
ps aux
kill -9 


·······································
scp -r    root@172.18.197.98:/mnt/4t2/bowei/sbw/gethv3/import/11-12Mnew.blockchain  ./import/

scp -r  ./1100w/*  root@172.18.197.98:/mnt/4t1/sbw/1200w/

192.168.0.101

192.168.1.116
scp -r   192.168.1.116:/mnt/a/geth/exportDir/0_1200w.block  ./




./build/bin/geth import /mnt/4t2/bowei/sbw/import/13-14M.blockchain  --datadir "/mnt/sbw/gethData/1300w/" --cache=14436  --snapshot=false


./build/bin/geth import /mnt/4t2/bowei/sbw/import/12-13Mnew.blockchain  --datadir "/mnt/4t1/sbw/1200w" --cache=14436  --snapshot=false

/mnt/4t1/sbw/1200w
·······································
备份数据


查看磁盘剩余空间
df -h
查看文件大小
du * --max-depth=0 -h


scp -r  ./1100w/*  root@172.18.197.98:/mnt/4t/sbw/1300w/


scp -r  dell@172.18.197.97:/home/dell/sbw/import/13-14M.blockchain  /mnt/4t2/bowei/sbw/import/

scp -r  dell@172.18.197.97:/mnt/4t/sbw/1300w/* ./

scp -r  dell@172.18.197.97:/home/dell/sbw/geth/output/* ./output_97




output存储位置：
/mnt/4t3/geth/output

data存储位置:
/mnt/4t3/geth/data/csv
/mnt/4t3/geth/data/dict


xblock

exported

import

 find . -maxdepth 1 ! -name ./import -exec rm -f {} ;

rm -f !(import|xblock)



vi /etc/resolv.conf



cp -r output_97/output_tx/1300To/* output/output_tx/1300To/
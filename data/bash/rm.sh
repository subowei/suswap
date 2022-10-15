shopt -s extglob
filename=`ls -t /home/dell/sbw/reorder_geth/data/output/output_addrInfo/1300ToV1/|head -n1|awk '{print $0}'`
cd /home/dell/sbw/reorder_geth/data/output/output_addrInfo/1300ToV1/
rm -rf !($filename)

filename=`ls -t /home/dell/sbw/reorder_geth/data/output/output_addrInfo_erc20/1300ToV1/|head -n1|awk '{print $0}'`
cd /home/dell/sbw/reorder_geth/data/output/output_addrInfo_erc20/1300ToV1/
rm -rf !($filename)

filename=`ls -t /home/dell/sbw/reorder_geth/data/output/output_addrInfo_erc721/1300ToV1/|head -n1|awk '{print $0}'`
cd /home/dell/sbw/reorder_geth/data/output/output_addrInfo_erc721/1300ToV1/
rm -rf !($filename)

filename=`ls -t /home/dell/sbw/reorder_geth/data/output/output_event_erc721/1300ToV1/|head -n1|awk '{print $0}'`
cd /home/dell/sbw/reorder_geth/data/output/output_event_erc721/1300ToV1/
rm -rf !($filename)

filename=`ls -t /home/dell/sbw/reorder_geth/data/output/output_tx/1300ToV1/|head -n1|awk '{print $0}'`
cd /home/dell/sbw/reorder_geth/data/output/output_tx/1300ToV1/
rm -rf !($filename)

# another

# filename=`ls -t /home/dell/sbw/gethv7/data/output_addrInfo/1000To/|head -n1|awk '{print $0}'`
# cd /home/dell/sbw/gethv7/data/output_addrInfo/1000To/
# rm -rf !($filename)

# filename=`ls -t /home/dell/sbw/gethv7/data/output_addrInfo_erc20/1000To/|head -n1|awk '{print $0}'`
# cd /home/dell/sbw/gethv7/data/output_addrInfo_erc20/1000To/
# rm -rf !($filename)

# filename=`ls -t /home/dell/sbw/gethv7/data/output_addrInfo_erc721/1000To/|head -n1|awk '{print $0}'`
# cd /home/dell/sbw/gethv7/data/output_addrInfo_erc721/1000To/
# rm -rf !($filename)

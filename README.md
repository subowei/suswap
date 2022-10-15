## Building the source

Building `geth` requires both a Go (version 1.16 or later) and a C compiler. You can install
them using your favourite package manager. Once the dependencies are installed, run

```shell
make geth
```

## Document Description

The document we modified is as follows:
1. core/blockchain.go
The core of reorder-geth. In function of insertChain, it will execute the transaction twice and use sorted transactions for the second time. Then compare the results after two executions and store the changed account balance, token balance, etc. 

2. core/state_processor.go
Include reordering algorithm which accords to effectiveGasTipValue and nonce

3. data/analyse/code
Inside is the code for processing the geth output
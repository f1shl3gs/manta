# Raft

## Notice
BoltDB's performance is not good enough as epected even if we only call sync every 100ms or when
we apply every 10000 entries.  

```text
Total writes    50000
Time:           16.794864529s
Writers:        32
Throughtput:    12.793364 MB/s
QPS:            2977.100525
DB Size:        283.562500 MB
```

## TODO
- Read is snapshot read so `dirty read` will never happen, `WriteTx` is `ReadTx` with `write cache`,
when WriteTx commit it propsal the txn to raft, when raft apply this txn request, the data we read might changed.
if conflict is used, client should implement retry.

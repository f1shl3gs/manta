# Raft

## TODO
- Read is snapshot read so `dirty read` will never happen, `WriteTx` is `ReadTx` with `write cache`,
when WriteTx commit it propsal the txn to raft, when raft apply this txn request, the data we read might changed.
if conflict is used, client should implement retry.

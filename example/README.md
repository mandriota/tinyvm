# Examples

1. Sum Of Array
```
INT GET
POPC
PUSH 0

SumArray:
  INT GET
  ADD
  PUSH SumArray
  LOOP

INT PUT
```
Binary code you can find in test_sumArr.bin

2. 2 Powers
```
INT GET
POPC
PUSH 1

Squares:
  DUP
  DUP
  INT PUT
  ADD
  PUSH Squares
  LOOP
```
Binary code you can find in test_2Pow.bin

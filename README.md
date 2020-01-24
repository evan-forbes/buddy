# Buddy

__This project is a work in progress and I'm adding commands as I need them.__ Unless of course you want to add any, which you totally should!

Interface focused go bindings for ethereum contracts. Allows for easy mock testing for smart contracts in go (!), amongst other quality of life improvements for those of us who prefer to use go to interact with smart contracts.

### Cool Stuff

While generating go bindings for smart contracts is nothing new, these bindings are focused on simpler more solidity like usage.

```go
package main

import (
    "github.com/your-username/coin" // import your generated bindings
)

type Erc20 interface {
    Approve(opts *bind.TransactOpts, usr common.Address, wad *big.Int) (*types.Transaction, error)
    Transfer(opts *bind.TransactOpts, dst common.Address, wad *big.Int) (*types.Transaction, error)
}

func Erc20Procedure(erc Erc20) {
    // do stuff with erc20s
}

func main() { 
    c, err := coin.NewCoin(common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f"))
    if err != nil {
        log.Fatal(err)
    }

    Erc20Procedure(c)
}

```
Doing something similar in the standard abigen in go-ethereum would require a bit more work and a lot less flexibility. 

### Dependencies

[go-ethereum](https://github.com/ethereum/go-ethereum)
### Code Generation using abigen sub command
Use the abigen command to generate go bindings from an abi file. If you include a binary of the compiled contract, then a deployment method is also generated.
```
buddy abigen "path/to/abi/and/maybe/a/bin" -p packageName
```
or in the working directory with a specified abi or bin file
```
buddy abigen --pkg=packageName --abi=contract.abi --bin=contract.bin
```
Multiple contracts can be generated into the same package if so desired.


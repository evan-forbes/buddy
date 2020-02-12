# Buddy

__This project is a work in progress and I'm adding commands as I need them.__

Interface focused go bindings for ethereum contracts. Allows for easy mock testing for smart contracts in go (!), amongst other quality of life improvements for those of us who prefer to use go to interact with smart contracts.

### Dependencies

[go-ethereum](https://github.com/ethereum/go-ethereum)

## Usage

### Smart contract binding generation using abigen sub command
Use the abigen command to generate go bindings from an abi file. If you include a binary of the compiled contract, then a deployment method is also generated. This command uses much of the same code as the abigen cli command included in go-ethereum
```
buddy abigen "path/to/abi/and/maybe/a/bin" -p packageName
```
or in the working directory with a specified abi or bin file
```
buddy abigen --pkg=packageName --abi=contract.abi --bin=contract.bin
```
Multiple contracts can be generated into the same package if so desired.

### Cool Stuff

While generating go bindings for smart contracts is nothing new, these bindings allow one to write go interfaces for generated code.

```go
package main

import (
    "github.com/your-username/coin" // import your generated bindings
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
)

type Erc20 interface {
    Approve(opts *bind.TransactOpts, usr common.Address, wad *big.Int) (*types.Transaction, error)
    Transfer(opts *bind.TransactOpts, dst common.Address, wad *big.Int) (*types.Transaction, error)
}

func Erc20Procedure(erc Erc20) {
    // do stuff with erc20s
}

func main() { 
    c, err := uniswap.NewUniswapExchange(common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f"))
    if err != nil {
        log.Fatal(err)
    }

    Erc20Procedure(c)

    // while still retaining the ability to call contract specific methods without a seperate instance
    c.ContractSpecificMethodNotInTheInterface(opts, new(big.Int).SetFromString("12345678910111213141516"))
}

``` 
Doing something similar with the official abigen command would require twice the compiling, code generating, and keeping track of two contract instances. 


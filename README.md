# Buddy

Interface focused go bindings for ethereum contracts. Allows for easy mock testing for smart contracts in go (!), amongst other quality of life improvements for those of us who prefer to use go to interact with smart contracts. Works by slightly tweaking the templates of go-etheruem along with organizing each contract into its own package.

## Dependencies

go-ethereum (probably something)

## Usage

### cli 

buddy abigen "path/to/abi/and/or/bin"
buddy solc "path/to/source"


# ivy-deploy

## Setup:
1. Download the latest version of [Go](https://go.dev/)
2. Copy the `sample.yaml` file to a new file named `config.yaml`
3. Fill in the network, your operator ID, and operator key. On the testnet and previewnet, this can be found by going to https://portal.hedera.com/register

## Hedera Crypto Service Interactions
### Transfer to an Account

Transfer HBAR to an account:
```
go run ivy.go --config config.yaml transfer --account 0.0.34806395 --hbar 10 --memo "Hello World"
```

Transfer TINYBAR to an account:
```
go run ivy.go --config config.yaml transfer --account 0.0.34806395 --tinybar 400000 --memo "Hello World"
```

Note: HBAR or TINYBAR must be present with a nonzero integer. If both are present, HBAR amount will be used.

### Get the Balance of an Account
```
go run ivy.go --config config.yaml balance --account 0.0.34818022
```

### Create Accounts
Create a set of accounts to test with using the following command
```
go run ivy.go --config config.yaml create-accounts --count 5 --out test.json
```

The `out` file will contain aarray of accounts with Account IDs and private keys

## Hedera Smart Contract Service Interactions
### Deploy a smart contract using a truffle file (must not have parameters in constructor):
```
go run ivy.go --config config.yaml deploy --truffle ~/workspace/truffle-project/build/contracts/contract.json --out deployment.json
```
The deployment.json will then contain the contract ID, fees paid during deployment, and the solidity address for the contract

### Get a Smart Contract's Info
```
go run ivy.go --config config.yaml contract-info --contract 0.0.34806395 --out contract-info.json
````

Where `out` is the json file that the smart contract info will be written to


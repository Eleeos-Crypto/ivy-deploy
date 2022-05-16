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

## Hedera Smart Contract Service Interactions
### Deploy a smart contract using a truffle file (must not have parameters in constructor):
```
go run ivy.go --config config.yaml deploy --truffle ~/workspace/truffle-project/build/contracts/contract.json --out deployment.json
```
The deployment.json will then contain the contract ID, fees paid during deployment, and the solidity address for the contract


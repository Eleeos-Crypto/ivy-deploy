# ivy-deploy

## Setup:
1. Download the latest version of [Go](https://go.dev/)
2. Copy the `sample.yaml` file to a new file named `config.yaml`
3. Fill in the network, your operator ID, and operator key. On the testnet and previewnet, this can be found by going to https://portal.hedera.com/register

## Hedera Smart Contract Service Interactions
### Deploy a smart contract using a truffle file (must not have parameters in constructor):
```
go run ivy.go --config config.yaml deploy --truffle ~/workspace/truffle-project/build/contracts/contract.json --out deployment.json
```
The deployment.json will then contain the contract ID, fees paid during deployment, and the solidity address for the contract


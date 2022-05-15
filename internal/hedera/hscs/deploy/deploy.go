package deploy

import (
	"fmt"

	"github.com/eleeos-crypto/ivy-deploy/internal/hedera/fs/chunked"
	"github.com/eleeos-crypto/ivy-deploy/internal/notify"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

type DeploymentDetail struct {
	ContractID         string `json:"contractId"`
	ContractSolidityID string `json:"contractSolidityID"`
	DeploymentFee      string `json:"deploymentFee"`
	DeploymentGas      string `json:"deploymentGas"`
}

func Deploy(client *hedera.Client, maxFee int64, gasAmount uint64, bytecode []byte) (*DeploymentDetail, error) {

	notify.PrettyInfo("Client Created")

	notify.PrettyInfo(fmt.Sprintf("Bytecode length: %d", len(bytecode)))

	byteCodeFileID, err := chunked.UploadLargeFile(client, maxFee, bytecode)
	if err != nil {
		notify.PrettyError("Error uploading chunked bytecode")
		return nil, err
	}

	notify.PrettyInfo(fmt.Sprintf("Contract File ID: %d", byteCodeFileID.File))

	// Instantiate the contract instance
	contractTransactionResponse, err := hedera.NewContractCreateTransaction().
		// Failing to set this to a sufficient amount will result in "INSUFFICIENT_GAS" status
		SetGas(gasAmount).
		// The file ID we got from the record of the file created previously
		SetBytecodeFileID(*byteCodeFileID).
		// No constructor parameters
		SetConstructorParameters(hedera.NewContractFunctionParameters()).
		// Setting an admin key allows you to delete the contract in the future
		SetAdminKey(client.GetOperatorPublicKey()).
		Execute(client)

	if err != nil {
		notify.PrettyError("Error creating smart contract")
		return nil, err
	}

	// get the record for the contract we created
	contractRecord, err := contractTransactionResponse.GetRecord(client)
	if err != nil {
		notify.PrettyError("Error retrieving smart contract create record")
		return nil, err
	}

	contractCreateResult, err := contractRecord.GetContractCreateResult()
	if err != nil {
		notify.PrettyError("Error retrieving smart contract create result")
		return nil, err
	}

	// get the contract ID from the record
	newContractID := *contractRecord.Receipt.ContractID

	notify.PrettyInfo(fmt.Sprintf("Contract create gas used: %v\n", contractCreateResult.GasUsed))
	notify.PrettyInfo(fmt.Sprintf("Contract create transaction fee: %v\n", contractRecord.TransactionFee))
	notify.PrettyInfo(fmt.Sprintf("Contract: %v\n", newContractID))

	dtl := DeploymentDetail{
		ContractID:         newContractID.String(),
		ContractSolidityID: newContractID.ToSolidityAddress(),
		DeploymentFee:      contractRecord.TransactionFee.String(),
		DeploymentGas:      fmt.Sprintf("%v", contractCreateResult.GasUsed),
	}

	return &dtl, nil
}

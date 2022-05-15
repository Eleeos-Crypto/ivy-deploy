package chunked

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
)

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

func UploadLargeFile(client *hedera.Client, maxFee int64, content []byte) (*hedera.FileID, error) {

	chunks := split(content, 2048)

	// Create a new file
	newFileResponse, err := hedera.NewFileCreateTransaction().
		SetKeys(client.GetOperatorPublicKey()).
		SetContents(chunks[0]).
		SetMaxTransactionFee(hedera.HbarFromTinybar(maxFee)).
		Execute(client)
	if err != nil {
		return nil, err
	}

	// Get receipt to make sure the transaction worked
	receipt, err := newFileResponse.GetReceipt(client)
	if err != nil {
		return nil, err
	}

	// Retrieve file ID from the receipt
	fileID := *receipt.FileID

	for _, chunk := range chunks[1:] {
		// File append
		fileResponse, err := hedera.NewFileAppendTransaction().
			SetNodeAccountIDs([]hedera.AccountID{newFileResponse.NodeID}).
			SetFileID(fileID).
			SetContents(chunk).
			SetMaxTransactionFee(hedera.HbarFromTinybar(maxFee)).
			Execute(client)
		if err != nil {
			return nil, err
		}

		// Checking if transaction went through
		receipt, err = fileResponse.GetReceipt(client)
		if err != nil {
			return nil, err
		}
	}

	return &fileID, nil
}

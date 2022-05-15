package hedera

import "github.com/hashgraph/hedera-sdk-go/v2"

func SetupClient(network, operatorId, operatorKey string) (*hedera.Client, error) {
	var client *hedera.Client
	var err error

	client, err = hedera.ClientForName(network)
	if err != nil {
		return nil, err
	}

	if operatorId != "" && operatorKey != "" {
		operatorAccountID, err := hedera.AccountIDFromString(operatorId)
		if err != nil {
			return nil, err
		}

		operatorKey, err := hedera.PrivateKeyFromString(operatorKey)
		if err != nil {
			return nil, err
		}

		client.SetOperator(operatorAccountID, operatorKey)
	}

	return client, err
}

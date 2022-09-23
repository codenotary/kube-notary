package verify

import (
	"errors"
	"os"

	"github.com/codenotary/vcn-enterprise/pkg/api"
	"github.com/codenotary/vcn-enterprise/pkg/meta"
	"github.com/sirupsen/logrus"
)

// @TODO: Recovered code from vcn-enterprise repository, does it make sense to be here? CLARIFY
var ErrNoLcApiKeyEnv = errors.New(`no API key configured. Please set the environment variable VCN_LC_API_KEY=<API-KEY> or use --lc-api-key flag on each request before running any commands`)

// PublicCNLCVerify allow connection and verification on CNLC ledger with a single call.
// LcLedger parameter is used when a cross-ledger key is provided in order to specify the ledger on which future operations will be directed. Empty string is accepted
// signerID parameter is used to filter result on a specific signer ID. If empty value is provided is used the current logged signerID value.
func PublicCNLCVerify(hash, lcLedger, signerID, lcHost, lcPort, lcCert string, lcSkipTlsVerify, lcNoTls bool) (a *api.LcArtifact, err error) {
	logrus.WithFields(logrus.Fields{
		"hash": hash,
	}).Trace("LcVerify")

	apiKey := os.Getenv(meta.VcnLcApiKey)
	if apiKey == "" {
		logrus.Trace("Lc api key provided (environment)")
		return nil, ErrNoLcApiKeyEnv
	}

	client, err := api.NewLcClient(apiKey, lcLedger, lcHost, lcPort, lcCert, lcSkipTlsVerify, lcNoTls, nil)
	if err != nil {
		return nil, err
	}

	lcUser := &api.LcUser{Client: client}

	err = lcUser.Client.Connect()
	if err != nil {
		return nil, err
	}

	if hash != "" {
		a, _, err = lcUser.LoadArtifact(
			hash,
			signerID,
			"",
			0,
			map[string][]string{meta.VcnLCCmdHeaderName: {meta.VcnLCVerifyCmdHeaderValue}})
		if err != nil {
			return nil, err
		}
	}

	return a, nil
}

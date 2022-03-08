package verify

import "github.com/codenotary/vcn-enterprise/pkg/meta"

type Verification struct {
	Level   meta.Level  `json:"level"`
	Status  meta.Status `json:"status"`
	Date    string      `json:"timestamp"`
	Trusted bool        `json:"trusted"`
}

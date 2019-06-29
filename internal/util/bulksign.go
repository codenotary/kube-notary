package util

import (
	"text/template"
)

const script = `#!/usr/bin/env bash
        
set -euo pipefail

# Ensure authentication
vcn login

# Read the passphrase
echo -n Key passphrase: 
read -s passphrase
echo
echo

export KEYSTORE_PASSWORD=$passphrase

# To use a different signing key, please uncomment the following line (require vcn v0.5.1 or above)
#export VCN_KEY=<key>

# Bulk sign
echo Signing...
echo

{{ range .Results -}}
vcn s --hash {{ .Hash }} --name "{{ if not .Containers -}}
		sha256:{{ .Hash }}
	{{- else -}}
		{{ with index .Containers 0 }}{{ .Image }}{{ end }}
	{{- end}}"
{{ end }}
`

var BulkSigningScriptTemplate = template.Must(template.New("script").Parse(script))

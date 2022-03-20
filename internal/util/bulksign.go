package util

import (
	"text/template"
)

const script = `#!/usr/bin/env bash
        
set -euo pipefail

# Read the login variables
echo -n CNC API key: 
read -s apikey
echo

echo -n CNC API Host:
read -s cnchost
echo

echo

export VCN_LC_API_KEY=$apikey
export VCN_LC_HOST=$cnchost

# Ensure authentication
vcn login


# Bulk sign
echo Signing...
echo

{{ range .Results -}}
{{ if .Hash -}}
vcn n --hash {{ .Hash }} --name "{{ if not .Containers -}}
		sha256:{{ .Hash }}
	{{- else -}}
		{{ with index .Containers 0 }}{{ .Image }}{{ end }}
	{{- end}}"
{{- end}}
{{ end }}
`

var BulkSigningScriptTemplate = template.Must(template.New("script").Parse(script))

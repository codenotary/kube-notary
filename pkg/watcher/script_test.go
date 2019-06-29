package watcher

import (
	"bytes"
	"testing"
)

var results = []Result{
	Result{
		Hash: "ed72e25b3d2033bf74f2b110a4ddc283ec3f404e2db611caf0d608ef8c3314f4",
		Containers: []ContainerInfo{
			ContainerInfo{
				Namespace:   "default",
				Pod:         "flailing-donkey-prometheus-alertmanager-556cb88cf4-5c45x",
				Container:   "prometheus-alertmanager",
				ContainerID: "containerd://c209aa6d772e55acc77726b6793f245db646b5e698c55cb1fca5bb6c19d585e3",
				Image:       "docker.io/prom/alertmanager:v0.15.3",
				ImageID:     "docker.io/prom/alertmanager@sha256:196af0317d3449c1300aa26ff0366f68c67d04581d2c9f8609cbb227424e226c",
			},
		},
	},
	Result{
		Hash: "7a344aad0fdbe8fd3ebd3ace7268d59946408503db1fe7c171bdb016a51729b7",
	},
}

func TestBulkSigningScript(t *testing.T) {
	buf := &bytes.Buffer{}
	err := bulkSigningScript(buf, results)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buf.String())
}

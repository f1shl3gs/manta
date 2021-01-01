package membership

import (
	"time"

	"github.com/hashicorp/memberlist"
)

type Metadata struct {
	Service string `json:"svc"`
}

type Cluster struct {
	ml *memberlist.Memberlist
}

func (cluster *Cluster) Unregister() {

}

func (cluster *Cluster) Stop() error {
	timeout := 5 * time.Second
	return cluster.ml.Leave(timeout)
}

func (cluster *Cluster) Register(m Metadata) {

}

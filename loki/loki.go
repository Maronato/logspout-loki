package loki

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gliderlabs/logspout/router"
	lokiclient "github.com/livepeer/loki-client/client"
	"github.com/livepeer/loki-client/model"
)

const glogTimeFormat = "20060102 15:04:05.999999"

var (
	levels       = []string{"info", "warning", "error", "fatal"}
	errShortLine = errors.New("Too short line")
	year         string
)

func init() {
	year = strconv.FormatInt(int64(time.Now().Year()), 10)
	router.AdapterFactories.Register(NewLokiAdapter, "loki")
}

// LokiAdapter is an adapter that streams logs to Loki.
type LokiAdapter struct {
	route  *router.Route
	client *lokiclient.Client
}

func logger(v ...interface{}) {
	fmt.Println(v...)
}

// NewLokiAdapter creates a LokiAdapter.
func NewLokiAdapter(route *router.Route) (router.LogAdapter, error) {
	baseLabels := model.LabelSet{}
	lokiURL := "http://" + route.Address + "/api/prom/push"
	fmt.Printf("Using Loki url: %s\n", lokiURL)
	client, err := lokiclient.NewWithDefaults(lokiURL, baseLabels, logger)
	if err != nil {
		return nil, err
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go waitExit(client, c)

	return &LokiAdapter{
		route:  route,
		client: client,
	}, nil
}

// Stream implements the router.LogAdapter interface.
func (a *LokiAdapter) Stream(logstream chan *router.Message) {
	defer a.client.Stop()

	for m := range logstream {

		container := m.Container.Name[1:]
		service := m.Container.Config.Labels["com.docker.swarm.service.name"]
		node := m.Container.Config.Labels["com.docker.swarm.node.id"]
		labels := model.LabelSet{"container": container, "service": service, "node": node}
		line := strings.TrimSpace(m.Data)
		a.client.Handle(labels, m.Time, line)
	}
}

func waitExit(client *lokiclient.Client, c chan os.Signal) {
	<-c
	client.Stop()
}

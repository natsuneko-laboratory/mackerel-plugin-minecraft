package main

import (
	"flag"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/sch8ill/mclib"
)

// Mackerel Agent
type MinecraftPlugin struct {
	Prefix string
	Server string
}

func (mc MinecraftPlugin) FetchMetrics() (map[string]float64, error) {
	client, err := mclib.NewClient(mc.GetServerAddress())
	if err != nil {
		return nil, err
	}

	res, err := client.StatusPing()
	if err != nil {
		return nil, err
	}

	return map[string]float64{
		"max_players":    float64(res.Players.Max),
		"online_players": float64(res.Players.Online),
		"latency":        float64(res.Latency),
	}, nil
}

func (mc MinecraftPlugin) GraphDefinition() map[string]mp.Graphs {
	prefix := mc.GetPrefix()

	return map[string]mp.Graphs{
		prefix: {
			Label: "Minecraft Server Status",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "max_players", Label: "Limit"},
				{Name: "online_players", Label: "Current Players"},
				{Name: "latency", Label: "Latency"},
			},
		},
	}
}

func (mc MinecraftPlugin) GetPrefix() string {
	if mc.Prefix == "" {
		mc.Prefix = "minecraft"
	}

	return mc.Prefix
}

func (mc MinecraftPlugin) GetServerAddress() string {
	if mc.Server == "" {
		mc.Server = "localhost:25565"
	}

	return mc.Server
}

func main() {
	prefix := flag.String("metric-key-prefix", "minecraft", "Metric key prefix")
	server := flag.String("server", "localhost:25565", "Server address")
	tempfile := flag.String("tempfile", "", "Temp file name")

	flag.Parse()

	mc := MinecraftPlugin{
		Prefix: *prefix,
		Server: *server,
	}

	helper := mp.NewMackerelPlugin(mc)
	helper.Tempfile = *tempfile

	helper.Run()
}

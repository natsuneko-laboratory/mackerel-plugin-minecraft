package main

import (
	"flag"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/sch8ill/mclib"
	"os"
	filepath2 "path/filepath"
	"strings"
)

// Mackerel Agent
type MinecraftPlugin struct {
	Prefix   string
	Server   string
	SaveData string
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

	var overworldSize int64 = 0
	var netherSize int64 = 0
	var theEndSize int64 = 0

	if mc.SaveData != "" {
		path, _ := filepath2.Abs(mc.SaveData)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, err
		}

		err = filepath2.Walk(path, func(data string, info os.FileInfo, err error) error {
			if filepath2.Ext(data) == ".dat" {
				rel, err := filepath2.Rel(path, data)
				if err != nil {
					return err
				}

				// Dimension: -1 is Nether
				if strings.HasPrefix(rel, "DIM-1") {
					netherSize += info.Size()
					return nil
				}

				// Dimension: 1 is The End
				if strings.HasPrefix(rel, "DIM1") {
					theEndSize += info.Size()
					return nil
				}

				// Dimension: 0 is Overworld
				if strings.HasPrefix(rel, "data") {
					overworldSize += info.Size()
					return nil
				}
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return map[string]float64{
		"max_players":    float64(res.Players.Max),
		"online_players": float64(res.Players.Online),
		"latency":        float64(res.Latency),
		"overworldSize":  float64(overworldSize) / float64(1024),
		"netherSize":     float64(netherSize) / float64(1024),
		"the_end_size":   float64(theEndSize) / float64(1024),
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
				{Name: "overworld_size", Label: "Overworld DataSize"},
				{Name: "nether_size", Label: "Nether DataSize"},
				{Name: "the_end_size", Label: "TheEnd DataSize"},
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
	savedata := flag.String("savedata", "", "Minecraft saved data dir")
	tempfile := flag.String("tempfile", "", "Temp file name")

	flag.Parse()

	mc := MinecraftPlugin{
		Prefix:   *prefix,
		Server:   *server,
		SaveData: *savedata,
	}

	helper := mp.NewMackerelPlugin(mc)
	helper.Tempfile = *tempfile

	helper.Run()
}

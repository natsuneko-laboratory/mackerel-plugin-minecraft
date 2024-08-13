package main

import (
	"flag"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/natsuneko-laboratory/mackerel-plugin-minecraft/minecraft"
)

// MinecraftPlugin Mackerel Agent
type MinecraftPlugin struct {
	Prefix   string
	Server   string
	SaveData string
}

func (mc MinecraftPlugin) FetchMetrics() (map[string]float64, error) {
	status, err := minecraft.GetServerStatus(mc.GetServerAddress())
	if err != nil {
		return nil, err
	}

	data, err := minecraft.GetSaveDataStats(mc.SaveData)
	if err != nil {
		return nil, err
	}

	return map[string]float64{
		"max_players":           float64(status.MaxPlayers),
		"online_players":        float64(status.OnlinePlayers),
		"latency":               float64(status.Latency),
		"overworld.size":        float64(data.Overworld.SizeInBytes),
		"overworld.chunk.x_max": float64(data.Overworld.ChunkX.P * 32),
		"overworld.chunk.x_min": float64(data.Overworld.ChunkX.N * 32),
		"overworld.chunk.z_max": float64(data.Overworld.ChunkZ.P * 32),
		"overworld.chunk.z_min": float64(data.Overworld.ChunkZ.N * 32),
		"nether.size":           float64(data.Nether.SizeInBytes),
		"nether.chunk.x_max":    float64(data.Nether.ChunkX.P * 32),
		"nether.chunk.x_min":    float64(data.Nether.ChunkX.N * 32),
		"nether.chunk.z_max":    float64(data.Nether.ChunkZ.P * 32),
		"nether.chunk.z_min":    float64(data.Nether.ChunkZ.N * 32),
		"the_end.size":          float64(data.EndWorld.SizeInBytes),
		"the_end.chunk.x_max":   float64(data.EndWorld.ChunkX.P * 32),
		"the_end.chunk.x_min":   float64(data.EndWorld.ChunkX.N * 32),
		"the_end.chunk.z_max":   float64(data.EndWorld.ChunkZ.P * 32),
		"the_end.chunk.z_min":   float64(data.EndWorld.ChunkZ.N * 32),
		"world_total_size":      float64(0),
	}, nil
}

func (mc MinecraftPlugin) GraphDefinition() map[string]mp.Graphs {
	prefix := mc.GetPrefix()

	players := prefix + ".player"
	data := prefix + ".data"

	return map[string]mp.Graphs{
		players: {
			Label: "Minecraft Server Status",
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "max", Label: "Limit"},
				{Name: "online", Label: "Current Players"},
				{Name: "latency", Label: "Latency"},
			},
		},
		data: {
			Label: "Minecraft Server DataSize",
			Unit:  mp.UnitBytes,
			Metrics: []mp.Metrics{
				{Name: "overworld.size", Label: "Overworld DataSize"},
				{Name: "overworld.chunk.x_max", Label: "Overworld Chunk X Max"},
				{Name: "overworld.chunk.x_min", Label: "Overworld Chunk X Min"},
				{Name: "overworld.chunk.z_max", Label: "Overworld Chunk Z Max"},
				{Name: "overworld.chunk.z_min", Label: "Overworld Chunk Z Min"},
				{Name: "nether.size", Label: "Nether DataSize"},
				{Name: "nether.chunk.x_max", Label: "Nether Chunk X Max"},
				{Name: "nether.chunk.x_min", Label: "Nether Chunk X Min"},
				{Name: "nether.chunk.z_max", Label: "Nether Chunk Z Max"},
				{Name: "nether.chunk.z_min", Label: "Nether Chunk Z Min"},
				{Name: "the_end.size", Label: "TheEnd DataSize"},
				{Name: "the_end.chunk.x_max", Label: "TheEnd Chunk X Max"},
				{Name: "the_end.chunk.x_min", Label: "TheEnd Chunk X Min"},
				{Name: "the_end.chunk.z_max", Label: "TheEnd Chunk Z Max"},
				{Name: "the_end.chunk.z_min", Label: "TheEnd Chunk Z Min"},
				{Name: "world_total_size", Label: "World Total DataSize"},
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

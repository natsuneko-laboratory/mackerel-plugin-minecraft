package minecraft

import (
	"github.com/sch8ill/mclib"
)

type ServerStatus struct {
	Latency       int
	MaxPlayers    int
	OnlinePlayers int
}

func GetServerStatus(address string) (*ServerStatus, error) {
	client, err := mclib.NewClient(address)
	if err != nil {
		return nil, err
	}

	res, err := client.StatusPing()
	if err != nil {
		return nil, err
	}

	return &ServerStatus{
		Latency:       res.Latency,
		MaxPlayers:    res.Players.Max,
		OnlinePlayers: res.Players.Online,
	}, nil
}

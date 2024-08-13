package minecraft

import (
	"os"
	filepath2 "path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Chunk struct {
	P int64
	N int64
}

type DimensionData struct {
	SizeInBytes int64
	ChunkX      Chunk
	ChunkZ      Chunk
}

type SaveData struct {
	Overworld *DimensionData
	Nether    *DimensionData
	EndWorld  *DimensionData
}

func GetSaveDataStats(dir string) (*SaveData, error) {
	path, _ := filepath2.Abs(dir)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	data := &SaveData{
		Overworld: GetDimensionData(dir, ""),
		Nether:    GetDimensionData(dir, "DIM-1"),
		EndWorld:  GetDimensionData(dir, "DIM1"),
	}

	return data, nil
}

func GetDimensionData(dir string, dimension string) *DimensionData {
	var size int64 = 0

	path := filepath2.Join(dir, dimension)
	chunkX := Chunk{P: 0, N: 0}
	chunkZ := Chunk{P: 0, N: 0}
	pattern := regexp.MustCompile(`r\.(-?\d+)\.(-?\d+)\.mca$`)

	filepath2.Walk(path, func(data string, info os.FileInfo, err error) error {
		rel, err := filepath2.Rel(path, data)

		// ChunkData
		if filepath2.Ext(data) == ".mca" {
			if strings.HasPrefix(rel, "region") {
				size += info.Size()

				matches := pattern.FindAllStringSubmatch(rel, -1)[0]
				x, _ := strconv.ParseInt(matches[1], 10, 64)
				z, _ := strconv.ParseInt(matches[2], 10, 64)

				chunkX.P = max(chunkX.P, x)
				chunkX.N = min(chunkX.N, x)

				chunkZ.P = max(chunkZ.P, z)
				chunkZ.N = min(chunkZ.N, z)
			}
		}

		return nil
	})

	return &DimensionData{
		SizeInBytes: size,
		ChunkX:      chunkX,
		ChunkZ:      chunkZ,
	}
}

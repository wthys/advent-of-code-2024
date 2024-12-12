package day9

import (
	"fmt"
	"strings"
	"slices"

	"github.com/wthys/advent-of-code-2024/solver"
	"github.com/wthys/advent-of-code-2024/util"
	S "github.com/wthys/advent-of-code-2024/collections/set"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "9"
}

func (s solution) Part1(input []string, opts solver.Options) (string, error) {
	dm, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	sectors := dm.Sectors()
	opts.Debugf("START : %v\n", sectors)

	for true {
		firstGapIndex := -1
		for idx, sector := range slices.All(sectors) {
			if sector.id < 0 {
				firstGapIndex = idx
				break
			}
		}

		lastSectorIndex := len(sectors)
		for idx, sector := range slices.Backward(sectors) {
			if sector.id >= 0 {
				lastSectorIndex = idx
				break
			}
		}

		opts.Debugf("DURING: %v\n", sectors)

		if firstGapIndex > lastSectorIndex {
			break
		}

		// lets exchange
		left, right := sectors[firstGapIndex].Exchange(sectors[lastSectorIndex])
		sectors = Sectors(slices.Concat(
			sectors[:firstGapIndex],
			left,
			sectors[firstGapIndex+1:lastSectorIndex],
			right,
			sectors[lastSectorIndex+1:],
		)).Compact()
	}

	opts.Debugf("AFTER : %v\n", sectors)

	return solver.Solved(sectors.Checksum())
}

func (s solution) Part2(input []string, opts solver.Options) (string, error) {
	dm, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	sectors := dm.Sectors()
	seen := S.New[int]()

	for true {
		lastFileIndex := len(sectors)
		for idx, sector := range slices.Backward(sectors) {
			if seen.Has(sector.id) {
				continue
			}
			if sector.id >= 0 {
				lastFileIndex = idx
				break
			}
		}
		
		if lastFileIndex >= len(sectors) {
			break
		}

		fileToMove := sectors[lastFileIndex]
		seen.Add(fileToMove.id)
		opts.Debugf("COMPACTING %v : %v\n", fileToMove.id, sectors)

		firstSuitableGapIndex := -1
		for idx, sector := range slices.All(sectors) {
			if sector.id >= 0 {
				continue
			}
			if sector.length >= fileToMove.length {
				firstSuitableGapIndex = idx
				break
			}
		}

		if firstSuitableGapIndex < 0 || firstSuitableGapIndex >= lastFileIndex {
			continue
		}
		suitableGap := sectors[firstSuitableGapIndex]

		left, right := suitableGap.Exchange(fileToMove)
		sectors = Sectors(slices.Concat(
			sectors[:firstSuitableGapIndex],
			left,
			sectors[firstSuitableGapIndex+1:lastFileIndex],
			right,
			sectors[lastFileIndex+1:],
		)).Compact()
	}

	return solver.Solved(sectors.Checksum())
}

type (
	DiskMap []int
	Sectors []Sector
	Sector struct {
		length int
		id int
	}
)

func sector(id, length int) Sector {
	return Sector{length, id}
}

func (dm DiskMap) Sectors() Sectors {
	sectors := Sectors{}
	id := 0
	for idx, n := range dm {
		if idx % 2 == 0 {
			sectors = append(sectors, Sector{n, id})
			id += 1
		} else {
			sectors = append(sectors, Sector{n, -1})
		}
	}
	return sectors
}

func (sectors Sectors) DiskMap() DiskMap {
	dm := DiskMap{}
	lastId := -1
	for _, sector := range sectors {
		if lastId >= 0 && sector.id >= 0 {
			dm = append(dm, 0)
		}
		dm = append(dm, sector.length)
		lastId = sector.id
	}
	return dm
}

func (sectors Sectors) String() string {
	str := ""
	alfa := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ"
	for _, sector := range sectors {
		if sector.id >= 0 {
			str += strings.Repeat(string(alfa[sector.id % len(alfa)]), sector.length)
		} else {
			str += strings.Repeat(".", sector.length)
		}
	}
	return str
}

func (sectors Sectors) Compact() Sectors {
	if len(sectors) <= 1 {
		return sectors
	}

	compacted := Sectors{}
	newSector := sectors[0]
	lastId := newSector.id
	for _, sector := range sectors[1:] {
		if sector.id == lastId {
			newSector.length += sector.length
		} else {
			compacted = append(compacted, newSector)
			newSector = sector
			lastId = newSector.id
		}
	}

	return append(compacted, newSector)
}

func (sectors Sectors) Checksum() int {
	index := 0
	total := 0
	for _, sector := range sectors {
		n := sector.length
		if sector.id >= 0 {
			idx := index
			sectorSum := 0
			for idx < index + sector.length {
				sectorSum += sector.id * idx
				idx += 1
			}
			total += sectorSum
		}
		index += n
	}
	return total
}

func (left Sector) Exchange(right Sector) ([]Sector, []Sector) {
	if left.id >= 0 || right.id < 0 {
		return []Sector{left}, []Sector{right}
	}

	if left.length == right.length {
		return []Sector{
			sector(right.id, right.length),
		}, []Sector{
			sector(left.id, left.length),
		}
	}

	if left.length > right.length {
		return []Sector{
			sector(right.id, right.length),
			sector(left.id, left.length - right.length),
		}, []Sector{
			sector(-1, right.length),
		}
	}

	return []Sector{
		sector(right.id, left.length),
	}, []Sector{
		sector(right.id, right.length - left.length),
		sector(-1, left.length),
	}
}

func parseInput(input []string) (DiskMap, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("no disk map found")
	}

	values, _ := util.StringsToInts(util.ExtractRegex("[0-9]", input[0]))
	if len(values) == 0 {
		return nil, fmt.Errorf("no disk map found")
	}
	return DiskMap(values), nil
}
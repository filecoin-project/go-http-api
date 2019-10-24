package types

import "encoding/json"

// Protocol contains protocol-related settings for a Filecoin node
type Protocol struct {
	AutosealInterval uint64   `json:"autosealInterval,omitempty"`
	SectorSizes      []uint64 `json:"sectorSizes,omitempty"`
}

// BitswapStats contains Bitswap related stats for a Filecoin node
type BitswapStats struct{}

// Node contains the Node level Filecoin information
type Node struct {
	Kind         string       `json:"kind"`
	ID           string       `json:"id,omitempty"`
	Addresses    []string     `json:"addresses,omitempty"`
	Version      string       `json:"version,omitempty"`
	Commit       string       `json:"commit,omitempty"`
	Protocol     Protocol     `json:"protocol,omitempty"`
	BitswapStats BitswapStats `json:"bitswapStats,omitempty"`
}

// MarshalJSON marshals a Node struct
func (n Node) MarshalJSON() ([]byte, error) {
	type alias Node
	out := alias(n)
	out.Kind = "node"
	return json.Marshal(out)
}

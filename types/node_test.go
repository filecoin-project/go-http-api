package types_test

import (
	"testing"

	"github.com/filecoin-project/go-http-api/test"
	"github.com/filecoin-project/go-http-api/test/fixtures"
	. "github.com/filecoin-project/go-http-api/types"
)

func TestNode_MarshalJSON(t *testing.T) {
	t.Run("full struct is serialized correctly and includes Kind", func(t *testing.T) {
		n := Node{
			ID:           "abcd123",
			Addresses:    []string{fixtures.TestAddress0, fixtures.TestAddress1},
			Version:      "1",
			Commit:       "somehash",
			Protocol:     Protocol{AutosealInterval: 1, SectorSizes: []uint64{1, 2}},
			BitswapStats: BitswapStats{},
		}
		expected := `{"kind":"node","id":"abcd123","addresses":["t2gmpzificaunkf47tzkt377a6yllmcfj3g3qbyti","t12cvsox5neub6y4vupgsogbrfaljiot4eaenkkyy"],"version":"1","commit":"somehash","protocol":{"autosealInterval":1,"sectorSizes":[1,2]},"bitswapStats":{}}`
		test.AssertMarshaledEquals(t, n, expected)
	})
	t.Run("empty struct is serialized correctly and includes Kind", func(t *testing.T) {
		n := Node{}
		expected := `{"kind":"node","protocol":{},"bitswapStats":{}}`
		test.AssertMarshaledEquals(t, n, expected)
	})
}

package v1

import (
	"github.com/filecoin-project/go-http-api/types"
	logging "github.com/ipfs/go-log"
)

// For needed package-level items
var log = logging.Logger("rest-api-handlers")

// Callbacks is a struct for callbacks configurable for the given API endpoint,
// shown by the 'path' tag
// To add a new endpoint:
//   * Write a new handler to use a new V1Callback, with tests
//   * Add a new callback name/signature to Callbacks
//   * Add a case to SetupV1Handlers that uses the callback
type Callbacks struct {
	GetActorByID   func(string) (*types.Actor, error)
	GetActors      func() ([]*types.Actor, error)
	GetBlockByID   func(string) (*types.Block, error)
	CreateMessage  func(*types.Message) (*types.Message, error)
	GetMessageByID func(string) (*types.Message, error)
	GetNode        func() (*types.Node, error)
}

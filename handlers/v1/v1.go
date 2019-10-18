package v1

import (
	"github.com/filecoin-project/go-http-api/types"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
	"math/big"
	"net/http"
	"reflect"
)

// For needed package-level items
var log = logging.Logger("rest-api-handlers")

// Callbacks is a struct for callbacks to be used for an API endpoint.
// This is a struct rather than an interface so implementers do not have to write their own
// version of every callback.  Missing callbacks will cause a DefaultHandler to be used.
// To add a new endpoint:
//   * Write a new handler to use a new V1Callback, with tests
//   * Add a new callback name/signature to Callbacks
//   * Add a case to BuildHandlers that uses the callback
type Callbacks struct {
	GetActorByID        func(string) (*types.Actor, error)
	GetActors           func() ([]*types.Actor, error)
	GetBlockByID        func(string) (*types.Block, error)
	CreateMessage       func(*types.Message) (*types.Message, error)
	CreateSignedMessage func(*types.SignedMessage) (*types.SignedMessage, error)
	GetMessageByID      func(string) (*types.Message, error)
	GetNode             func() (*types.Node, error)
	SendSignedMessage   func(*types.SignedMessage) (*types.SignedMessage, error)
	WaitForMessage      func(cid *cid.Cid, limitBH *big.Int) (bH *big.Int, err error)
}

// BuildHandlers takes a V1Callback struct and iterates over all
// functions, creating a handler with a callback for each supported endpoint.
func (cb *Callbacks) BuildHandlers() *map[string]http.Handler {
	cb1t := reflect.TypeOf(*cb)
	cb1v := reflect.ValueOf(*cb)

	numCallbacks := cb1t.NumField()
	handlers := make(map[string]http.Handler, numCallbacks)

	for i := 0; i < numCallbacks; i++ {
		fieldName := cb1t.Field(i).Name

		fieldValue := cb1v.Field(i)
		if fieldValue.IsNil() {
			handlers[fieldName] = &DefaultHandler{}
		} else {
			switch fieldName {
			case "GetBlockByID":
				handlers[fieldName] = &BlockHandler{Callback: cb.GetBlockByID}
			case "GetActors":
				handlers[fieldName] = &ActorsHandler{Callback: cb.GetActors}
			case "GetActorByID":
				handlers[fieldName] = &ActorHandler{Callback: cb.GetActorByID}
			case "CreateMessage":
				handlers[fieldName] = &CreateMessageHandler{Callback: cb.CreateMessage}
			case "GetMessageByID":
				handlers[fieldName] = &MessageHandler{Callback: cb.GetMessageByID}
			case "GetNode":
				handlers[fieldName] = &NodeHandler{Callback: cb.GetNode}
			default:
				log.Errorf("skipping unknown handler: %s.", fieldName)
			}
		}
	}
	return &handlers
}

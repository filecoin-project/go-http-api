package v1

import (
	"math/big"
	"net/http"
	"reflect"

	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"

	"github.com/filecoin-project/go-http-api/types"
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
	// GetActorByID retrieves an Actor by its ID
	GetActorByID func(string) (*types.Actor, error)

	// GetActorNonce is specifically for retrieving the actor nonce in preparation for sending a signed message.
	GetActorNonce func(string) (*big.Int, error)

	// GetActors retrieves known information about all Actors of the node
	GetActors func() ([]*types.Actor, error)

	// GetBlockByID retrieves the BlockHeader, Messages and Receipts for the Block
	GetBlockByID func(string) (*types.Block, error)

	// CreateMessage creates and sends an unsigned Message
	CreateMessage func(*types.Message) (*types.SignedMessage, error)

	// CreateSignedMessage creates and sends a SignedMessage using the node's default
	// account
	CreateSignedMessage func(*types.SignedMessage) (*types.SignedMessage, error)

	// GetMessageByID fetches a Message by its CID
	GetMessageByID func(string) (*types.SignedMessage, error)

	// GetNode gets information about the node that implements this API
	GetNode func() (*types.Node, error)

	// SendSignedMessage posts an already signed message to the message pool.
	// Since actor Nonce is required to sign a message, the caller must first
	// know the actor nonce.  See GetActorNonce
	SendSignedMessage func(*types.SignedMessage) (*types.SignedMessage, error)

	// WaitForMessage waits for a message to appear on chain until the given block height.
	WaitForMessage func(*cid.Cid, *big.Int) (*types.SignedMessage, error)
}

// BuildHandlers takes a V1Callback struct and iterates over all
// functions, creating a handler with a callback for each supported endpoint.
func (cb *Callbacks) BuildHandlers() *map[string]http.Handler {
	cb1t := reflect.TypeOf(*cb)
	cb1v := reflect.ValueOf(*cb)

	numCallbacks := cb1t.NumField()
	handlers := make(map[string]http.Handler, numCallbacks)
	defH := &DefaultHandler{}
	for i := 0; i < numCallbacks; i++ {
		fieldName := cb1t.Field(i).Name

		fieldValue := cb1v.Field(i)
		if fieldValue.IsNil() {
			handlers[fieldName] = defH
		} else {
			switch fieldName {
			case "GetBlockByID":
				handlers[fieldName] = &BlockHandler{Callback: cb.GetBlockByID}
			case "GetActors":
				handlers[fieldName] = &ActorsHandler{Callback: cb.GetActors}
			case "GetActorByID":
				handlers[fieldName] = &ActorHandler{Callback: cb.GetActorByID}
			case "GetActorNonce":
				handlers[fieldName] = &ActorNonceHandler{Callback: cb.GetActorNonce}
			case "CreateMessage":
				handlers[fieldName] = &CreateMessageHandler{Callback: cb.CreateMessage}
			case "GetMessageByID":
				handlers[fieldName] = &MessageHandler{Callback: cb.GetMessageByID}
			case "GetNode":
				handlers[fieldName] = &NodeHandler{Callback: cb.GetNode}
			case "WaitForMessage":
				handlers[fieldName] = &WaitForMessageHandler{Callback: cb.WaitForMessage}
			default:
				log.Errorf("skipping unknown handler: %s.", fieldName)
			}
		}
	}
	return &handlers
}

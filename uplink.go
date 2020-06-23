// +build js

package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"storj.io/uplink"
)

// JsUplink holds the uplink functionalities exposed in the browser through
// Javascript.
type JsUplink struct {
	config uplink.Config
	access *uplink.Access
}

// NewJsUplink creates a new JsUplink instance returning a map with all the
// JsUplink methods to be exposed to the browser.
//
// It returns an error if there is an error mapping any of the methods to Js.
func NewJsUplink() (_ *JsUplink, ulkJsObj map[string]interface{}, _ error) {
	ulkn := &JsUplink{

		config: uplink.Config{
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				addressParts := strings.Split(address, ":")
				port, _ := strconv.Atoi(addressParts[1])
				return NewJsConn(addressParts[0], port)
			},
		},
	}

	ulkJsObj = make(map[string]interface{})

	jsFn, err := funcToJs(ulkn.LoadAccess)
	if err != nil {
		return nil, nil, err
	}
	ulkJsObj["loadAccess"] = jsFn

	return ulkn, ulkJsObj, nil
}

// LoadAccess loads a new access for being used in the next uplink calls.
// The previous access it's overridden by this new one.
func (ulkn *JsUplink) LoadAccess(satAddr, apiKey, passphrase string) string {
	access, err := ulkn.config.RequestAccessWithPassphrase(context.TODO(), satAddr, apiKey, passphrase)
	if err != nil {
		return fmt.Sprintf("could not request access grant: %+v", err)
	}

	ulkn.access = access
	return ""
}

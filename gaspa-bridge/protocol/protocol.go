// Package protocol contains protocol constants for Gaspa
package protocol

import (
	"net"

	"github.com/satori/uuid"
)

const (
	// DispatcherServiceRegister is member of enum DispatcherService
	DispatcherServiceRegister = 'r'
	// DispatcherServiceMeta is member of enum DispatcherService
	DispatcherServiceMeta = 'm'
	// DispatcherServiceMetaQuery is member of enum DispatcherService
	DispatcherServiceMetaQuery = 'q'
	// DispatcherServiceConnect is member of enum DispatcherService
	DispatcherServiceConnect = 'c'
	// DispatcherServiceJoin is member of enum DispatcherService
	DispatcherServiceJoin = 'j'
)

// RegisteredConnection is used in bridge meta
type RegisteredConnection struct {
	Name string       `json:"name"`
	UUID uuid.UUID    `json:"uuid"`
	Conn *net.TCPConn `json:"conn"`
}

// MetaResult of Meta query
type MetaResult struct {
	Status string                             `json:"status"`
	Result map[uuid.UUID]RegisteredConnection `json:"result"`
}

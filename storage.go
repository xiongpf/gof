/**
 * Copyright 2015 @ S1N1 Team.
 * name : storage.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */

package gof

// Storage
type Storage interface {
	// Get Value
	Get(key string, dst interface{}) error
	// Set Value
	Set(key string, v interface{}) error
	// Delete Storage
	Del(key string)
	// Auto Delete Set
	SetExpire(key string, v interface{}, seconds int64) error
}

// Package profile describes data models to interact with user profiles
// provided by Mixpanel Ingestion API.
package profile

// Set describes model of request to set values of user profile properties.
// If the profile does not exist, it creates it with these properties.
// If it does exist, it sets the properties to these values, overwriting existing values.
type Set struct {
	Token      string                 `json:"$token"`
	DistinctID string                 `json:"$distinct_id"`
	Set        map[string]interface{} `json:"$set"`
}

// SetOnce describes model of request to set values of user profile properties only once, without overwriting.
// If the profile does not exist, it creates it with these properties.
// If it does exist, it sets the properties to these values, overwriting existing values.
type SetOnce struct {
	Token      string                 `json:"$token"`
	DistinctID string                 `json:"$distinct_id"`
	SetOnce    map[string]interface{} `json:"$set_once"`
}

// NumberAdd describes model of request to increment or decrement numerical profile property value.
// If the property is not present on the profile, the value will be added to 0 (will set to specified in request).
type NumberAdd struct {
	Token      string                 `json:"$token"`
	DistinctID string                 `json:"$distinct_id"`
	Add        map[string]interface{} `json:"$add"`
}

// ListAppend describes model of request to append item from profile list property.
type ListAppend struct {
	Token      string                 `json:"$token"`
	DistinctID string                 `json:"$distinct_id"`
	Append     map[string]interface{} `json:"$append"`
}

// ListRemove describes model of request to remove item from profile list property.
type ListRemove struct {
	Token      string                 `json:"$token"`
	DistinctID string                 `json:"$distinct_id"`
	Remove     map[string]interface{} `json:"$remove"`
}

// Unset describes model of request to unset user profile property.
type Unset struct {
	Token      string   `json:"$token"`
	DistinctID string   `json:"$distinct_id"`
	Unset      []string `json:"$unset"`
}

// Mutator is internal interface to mark models as profile mutation actions.
type Mutator interface {
	isMutator()
}

func (x Set) isMutator()        {}
func (x SetOnce) isMutator()    {}
func (x NumberAdd) isMutator()  {}
func (x ListAppend) isMutator() {}
func (x ListRemove) isMutator() {}
func (x Unset) isMutator()      {}

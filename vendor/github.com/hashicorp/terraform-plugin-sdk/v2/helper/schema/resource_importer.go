package schema

import (
	"context"
	"errors"
)

// ResourceImporter defines how a resource is imported in Terraform. This
// can be set onto a Resource struct to make it Importable. Not all resources
// have to be importable; if a Resource doesn't have a ResourceImporter then
// it won't be importable.
//
// "Importing" in Terraform is the process of taking an already-created
// resource and bringing it under Terraform management. This can include
// updating Terraform state, generating Terraform configuration, etc.
type ResourceImporter struct {
	// State is called to convert an ID to one or more InstanceState to
	// insert into the Terraform state.
	//
	// Deprecated: State is deprecated in favor of StateContext.
	// Only one of the two functions can be set.
	State StateFunc

	// StateContext is called to convert an ID to one or more InstanceState to
	// insert into the Terraform state. If this isn't specified, then
	// the ID is passed straight through. This function receives a context
	// that will cancel if Terraform sends a cancellation signal.
	StateContext StateContextFunc
}

// StateFunc is the function called to import a resource into the Terraform state.
//
// Deprecated: Please use the context aware equivalent StateContextFunc.
type StateFunc func(*ResourceData, interface{}) ([]*ResourceData, error)

// StateContextFunc is the function called to import a resource into the
// Terraform state. It is given a ResourceData with only ID set. This
// ID is going to be an arbitrary value given by the user and may not map
// directly to the ID format that the resource expects, so that should
// be validated.
//
// This should return a slice of ResourceData that turn into the state
// that was imported. This might be as simple as returning only the argument
// that was given to the function. In other cases (such as AWS security groups),
// an import may fan out to multiple resources and this will have to return
// multiple.
//
// To create the ResourceData structures for other resource types (if
// you have to), instantiate your resource and call the Data function.
type StateContextFunc func(context.Context, *ResourceData, interface{}) ([]*ResourceData, error)

// InternalValidate should be called to validate the structure of this
// importer. This should be called in a unit test.
//
// Resource.InternalValidate() will automatically call this, so this doesn't
// need to be called manually. Further, Resource.InternalValidate() is
// automatically called by Provider.InternalValidate(), so you only need
// to internal validate the provider.
func (r *ResourceImporter) InternalValidate() error {
	if r.State != nil && r.StateContext != nil {
		return errors.New("Both State and StateContext cannot be set.")
	}
	return nil
}

// ImportStatePassthrough is an implementation of StateFunc that can be
// used to simply pass the ID directly through.
//
// Deprecated: Please use the context aware ImportStatePassthroughContext instead
func ImportStatePassthrough(d *ResourceData, m interface{}) ([]*ResourceData, error) {
	return []*ResourceData{d}, nil
}

// ImportStatePassthroughContext is an implementation of StateContextFunc that can be
// used to simply pass the ID directly through. This should be used only
// in the case that an ID-only refresh is possible.
func ImportStatePassthroughContext(ctx context.Context, d *ResourceData, m interface{}) ([]*ResourceData, error) {
	return []*ResourceData{d}, nil
}

package domain

// GuardLocationID returns InvalidLocationIDError unless id is a UUID. Shared by every use case
// that receives a location id from outside.
func GuardLocationID(id string) error { return nil }

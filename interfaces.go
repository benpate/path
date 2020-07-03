package path

// Getter interface allows other objects to make it easy to trace through their property trees, and get values from them.
type Getter interface {
	GetPath(Path) (interface{}, error)
}

// Setter interface allows other objects to make it easy to trace through their property trees, and set values into them.
type Setter interface {
	SetPath(Path, interface{}) error
}

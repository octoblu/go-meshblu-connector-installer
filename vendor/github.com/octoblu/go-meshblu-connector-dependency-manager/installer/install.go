package installer

// Install is a convenience method for constructing an installer client
// and calling client.Do
func Install(depType, tag string) error {
	client := New()
	return client.Do(depType, tag)
}

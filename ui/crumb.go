package ui

// Crumb is a screen the user is visiting. A list of Crumbs provides a path indicating
// where the user is in the screen hierarchy
type Crumb struct {
	URL  string
	Name string
}

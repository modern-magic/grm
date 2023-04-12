package action

type ViewOptions struct {
	All bool
}

type Action interface {
	View(option ViewOptions) int
	Join() int
	Drop() int
	Test() int
	Use() int
}
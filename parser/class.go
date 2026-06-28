package parser

type Class struct {
	Name string
}

func InitClass(name string) *Class {
	return &Class{
		Name: name,
	}
}

func (c *Class) String() string {
	return c.Name
}

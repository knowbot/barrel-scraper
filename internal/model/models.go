package model

type Named interface {
	GetName() string
}

type Category struct {
	Name          string
	SubCategories []SubCategory
}

func (c Category) GetName() string { return c.Name }

type SubCategory struct {
	Name string
	URL  string
}

func (sc SubCategory) GetName() string { return sc.Name }

type Region struct {
	Name      string
	Provinces []Province
}

func (r Region) GetName() string { return r.Name }

type Province struct {
	Name string
	Code string
}

func (p Province) GetName() string { return p.Name }

func BuildMap[T Named](items []T) map[string]T {
	m := make(map[string]T, len(items))
	for _, item := range items {
		m[item.GetName()] = item
	}
	return m
}

func BuildList[T Named](items []T) []string {
	l := make([]string, len(items))
	for i, item := range items {
		l[i] = item.GetName()
	}
	return l
}

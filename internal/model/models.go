package model

type Named interface {
	GetName() string
}

type Category struct {
	ID            int64
	Name          string
	SubCategories []SubCategory
}

func (ca Category) GetName() string { return ca.Name }

type SubCategory struct {
	ID   int64
	Name string
	URL  string
}

func (sc SubCategory) GetName() string { return sc.Name }

type Region struct {
	ID        int64
	Name      string
	Provinces []Province
}

func (re Region) GetName() string { return re.Name }

type Province struct {
	ID   int64
	Name string
	Code string
}

func (pr Province) GetName() string { return pr.Name }

type Company struct {
	ID            int64
	Sector        string
	Name          string
	StreetAddress string
	CAP           string
	City          string
	Province      string
	Phone         string
	Fax           string
	Website       string
}

func (co Company) GetName() string { return co.Name }

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

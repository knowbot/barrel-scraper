package main

type Category struct {
	Name          string
	SubCategories []SubCategory
}

type SubCategory struct {
	Name string
	URL  string
}

type Region struct {
	Name      string
	Provinces []Province
}

type Province struct {
	Name string
	Code string
}

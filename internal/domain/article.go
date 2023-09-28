package domain

type Article struct {
	Id        uint64
	Title     string
	Content   string
	ImageList []string
	Author    Author
}

type Author struct {
	Id   uint64
	Name string
}

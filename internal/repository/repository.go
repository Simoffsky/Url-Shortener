package repository

type LinksRepository interface {
	CreateLink(url, short string) (error)
	GetLink(short string) (string, error)
}

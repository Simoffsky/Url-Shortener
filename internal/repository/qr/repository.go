package qr

type QrRepository interface {
	CreateQRCode(link string, imgSize int) ([]byte, error)
	GetQRCode(link string) ([]byte, error)
	DeleteQRCode(link string) error
}


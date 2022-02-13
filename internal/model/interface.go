package model

type RequestDataService interface {
	Start()
	GetData(key int) (string, error)
}

type CacheDataRepository interface {
	GetData(key int) ([]byte, error)
	SetData(key int, data []byte)
}

package downloader

//Service is an interface used for downloading an resource from provided @URL at @DownloadLocation with @FileName
type Service interface {
	Download(conf *Config) error
}

//Config is a struct provides a download information to Download Service
type Config struct {
	URL              string
	FileName         string
	DownloadLocation string
	TransactionID    string
	CheckSum         string
	KeepOriginalName bool
	CheckSumType     string
	Header           map[string]string
}

package backend

//Log type
type Log struct {
	UUID              string `json:"uuid"`
	TimeStamp         int64  `json:"timestamp"`
	HumanReadableTime string `json:"humantime"`
	LogClass          string `json:"class"`
	Category          string `json:"category"`
	Body              string `json:"body"`
}

//LogStore : Log storage specification
type LogStore interface {
	Put(data []*Log) error
	Del(uuid string) error
	Get(uuid string) (*Log, error)
}

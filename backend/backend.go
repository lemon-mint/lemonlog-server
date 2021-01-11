package backend

//Log type
type Log struct {
	UUID              string
	TimeStamp         int64
	HumanReadableTime string
	LogClass          string
	Category          string
	Body              string
}

//LogStore : Log storage specification
type LogStore interface {
	Put(data []*Log) error
	Del(uuid string) error
	Get(uuid string) (*Log, error)
}

package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/lemon-mint/lemonlog-server/backend"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lemon-mint/godotenv"
)

var mq = make(map[string]*backend.Log)
var mqLock sync.Mutex
var counter int = 0
var uuidstart = time.Now().UnixNano()
var uuidstate0 = time.Now().UnixNano()
var uuidstate1 = uint32(time.Now().UnixNano() % 2147483647)
var uuidcounter = 0

func getQ(lim int) []*backend.Log {
	i := 0
	logs := make([]*backend.Log, 0, lim)
	mqLock.Lock()
	for key := range mq {
		if !(i < lim) {
			break
		}
		logs = append(logs, mq[key])
		delete(mq, key)
		i++
	}
	mqLock.Unlock()
	return logs
}

func addQ(addr *backend.Log) {
	mqLock.Lock()
	counter++
	mq[strconv.Itoa(counter)] = addr
	mqLock.Unlock()
}

func genUUID() string {
	uuidstate0 += time.Now().UnixNano()
	uuidstate1 += uint32(time.Now().UnixNano() % 2147483647)
	uuidcounter++
	return fmt.Sprintf("%x%x%x%x", uuidstart, uuidstate0, uuidstate1, uuidcounter)
}

func main() {
	godotenv.Load()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "static")

	e.GET("/log/append", appendLog)
	e.PUT("/log/append", appendLog)
	e.POST("/log/append", appendLog)

	if os.Getenv("PORT_FROM_ENV") != "" {
		e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
	} else {
		e.Logger.Fatal(e.Start(":16745"))
	}
}

type logRequest struct {
	Class    string `json:"class" form:"class" query:"class"`
	Category string `json:"category" form:"category" query:"category"`
	Data     string `json:"data" json:"data" json:"data"`
}

func appendLog(c echo.Context) error {
	logreq := new(logRequest)
	err := c.Bind(logreq)
	if err != nil {
		return c.JSONPretty(400, struct {
			Success bool       `json:"success"`
			Request logRequest `json:"request"`
		}{false, logRequest{}}, "  ")
	}
	t := time.Now().UTC()
	uuid := genUUID()
	addQ(&backend.Log{
		UUID:              uuid,
		TimeStamp:         t.Unix(),
		HumanReadableTime: t.String(),
		LogClass:          logreq.Class,
		Category:          logreq.Category,
		Body:              logreq.Data,
	})
	return c.JSONPretty(200, struct {
		Success bool       `json:"success"`
		Request logRequest `json:"request"`
		UUID    string     `json:"uuid"`
	}{true, *logreq, uuid}, "  ")
}

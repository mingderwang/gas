package middleware

import (
	"github.com/go-gas/gas"
	"github.com/go-gas/gas/logger"

	"net"
	"strconv"
	"time"
)

func LogMiddleware(next gas.GasHandler) gas.GasHandler {
	return func(c *gas.Context) error {
		// req := c.Request()
		// res := c.Response()
		l := logger.New("log/logs.txt")

		remoteAddr := c.RemoteAddr().String()
		if ip := string(c.Request.Header.Peek(gas.XRealIP)); ip != "" {
			remoteAddr = ip
		} else if ip = string(c.Request.Header.Peek(gas.XForwardedFor)); ip != "" {
			remoteAddr = ip
		} else {
			remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
		}

		start := time.Now()

		err := next(c)

		stop := time.Now()
		method := string(c.Method())
		path := string(c.Path())
		if path == "" {
			path = "/"
		}
		// size := c.Writer.Size()

		status := c.Response.StatusCode()//RespWriter.Status()

		// logger.Printf(format, remoteAddr, method, path, code, stop.Sub(start), size)

		logstr := "[" + start.Format("2006-01-02 15:04:05") + "][" + strconv.Itoa(status) + "][" + remoteAddr + "] " + method + " " + path + " Params: " + c.Request.PostArgs().String() + " ExecTime: " + stop.Sub(start).String()
		l.Info(logstr)

		return err
	}
}

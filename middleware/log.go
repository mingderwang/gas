package middleware

import (
	"github.com/gowebtw/gas"
	"github.com/gowebtw/gas/logger"

	"net"
	"strconv"
	"time"
)

func LogMiddleware(next gas.CHandler) gas.CHandler {
	return func(c *gas.Context) error {
		// req := c.Request()
		// res := c.Response()
		l := logger.New("log/logs.txt")

		remoteAddr := c.Req.RemoteAddr
		if ip := c.Req.Header.Get(gas.XRealIP); ip != "" {
			remoteAddr = ip
		} else if ip = c.Req.Header.Get(gas.XForwardedFor); ip != "" {
			remoteAddr = ip
		} else {
			remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
		}

		start := time.Now()

		err := next(c)

		stop := time.Now()
		method := c.Req.Method
		path := c.Req.URL.Path
		if path == "" {
			path = "/"
		}
		// size := c.Writer.Size()

		status := c.RespWriter.Status()

		// logger.Printf(format, remoteAddr, method, path, code, stop.Sub(start), size)

		logstr := "[" + start.Format("2006-01-02 15:04:05") + "][" + strconv.Itoa(status) + "][" + remoteAddr + "] " + method + " " + path + " Params: " + c.Req.Form.Encode() + " ExecTime: " + stop.Sub(start).String()
		l.Info(logstr)

		return err
	}
}

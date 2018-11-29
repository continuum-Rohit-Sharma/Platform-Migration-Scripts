package logging

import (
	glog "log"
)

const glogFlags int = glog.Ldate | glog.Ltime | glog.LUTC | glog.Lshortfile

type glogFactory struct{}

func (*glogFactory) New() goLogger {
	return glog.New(GetLoggerFactory().GetWriter().Get(), "", glogFlags)
}

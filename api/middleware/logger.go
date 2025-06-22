package middleware

import (
	"log"
	"log/syslog"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
  return func(c *gin.Context) {
	  c.Next()
	  file, err := syslog.New(syslog.LOG_SYSLOG, "GO-TODO")
	  if err != nil {
		  log.Fatalln("Unable to set logfile:", err.Error())
		}
		log.SetOutput(file)
		c.ClientIP()
		log.Printf("%s %s %s %d \"%s\"\n",
		c.ClientIP(),
		c.Request.Method,
		c.FullPath(),
		c.Writer.Status(),
		c.Request.UserAgent(),
		)
  }
}
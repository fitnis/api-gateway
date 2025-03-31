package proxy

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.Any("/appointments/*path", forward("http://appointment-service:8081"))
	r.Any("/orders/*path", forward("http://order-service:8083"))
	r.Any("/patients/*path", forward("http://patient-service:8082"))
	r.Any("/records/*path", forward("http://record-service:8084"))
	r.Any("/lab/*path", forward("http://lab-service:8085"))
	r.Any("/prescriptions/*path", forward("http://prescription-service:8086"))
	r.Any("/referrals/*path", forward("http://referral-service:8087"))
}

func forward(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := &http.Client{}
		req, err := http.NewRequest(c.Request.Method, target+c.Param("path"), c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "gateway error"})
			return
		}
		req.Header = c.Request.Header

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}

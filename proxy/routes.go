package proxy

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.Any("/api/appointments/*path", forward("http://appointment-service:8081/appointments"))
	r.Any("/api/orders/*path", forward("http://order-service:8083/orders"))
	r.Any("/api/patients/*path", forward("http://patient-service:8082/patients"))
	r.Any("/api/records/*path", forward("http://record-service:8084/records"))
	r.Any("/api/lab/*path", forward("http://lab-service:8085/lab"))
	r.Any("/api/prescriptions/*path", forward("http://prescription-service:8086/prescriptions"))
	r.Any("/api/referrals/*path", forward("http://referral-service:8087/referrals"))
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

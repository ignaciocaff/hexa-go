package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

type RequestDetail struct {
	Time   time.Time   `json:"time"`
	Method string      `json:"method"`
	Body   interface{} `json:"body,omitempty"`
	Params interface{} `json:"params,omitempty"`
}

type ResponseStats struct {
	StatusCode   int
	Count        int
	ErrorMessage string
	//Requests     []RequestDetail
}

type RouteStats struct {
	Path  string
	Stats []*ResponseStats
}

var routeStatsMap = make(map[string]*RouteStats)

func StatsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		statusCode := c.Writer.Status()
		errorMessage := ""

		if c.Writer.Status() >= 400 {
			errorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		}

		routePath := c.FullPath()

		if routeStats, ok := routeStatsMap[routePath]; ok {
			var existingStat *ResponseStats
			for _, stat := range routeStats.Stats {
				if stat.StatusCode == statusCode {
					existingStat = stat
					break
				}
			}

			if existingStat != nil {
				existingStat.Count++
				existingStat.ErrorMessage = errorMessage
			} else {
				routeStats.Stats = append(routeStats.Stats, &ResponseStats{
					StatusCode:   statusCode,
					Count:        1,
					ErrorMessage: errorMessage,
				})
			}
		} else {
			routeStatsMap[routePath] = &RouteStats{
				Path: routePath,
				Stats: []*ResponseStats{
					{
						StatusCode:   statusCode,
						Count:        1,
						ErrorMessage: errorMessage,
					},
				},
			}
		}
	}
}

func GetStatsMap() map[string]*RouteStats {
	return routeStatsMap
}

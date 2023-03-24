package service

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"github.com/goodguy-project/goodguy-crawl/handler"
	"github.com/goodguy-project/goodguy-crawl/proto"
	"github.com/goodguy-project/goodguy-crawl/util/jsonx"
)

func getRequestFromGin[T any](c *gin.Context) (T, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return lo.Empty[T](), err
	}
	data, err := jsonx.Unmarshal[T](body)
	if err != nil {
		return lo.Empty[T](), err
	}
	return data, nil
}

func RunHttpService() {
	r := gin.Default()
	r.POST("/GetContestRecord", func(c *gin.Context) {
		req, err := getRequestFromGin[*proto.GetContestRecordRequest](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		resp, err := handler.GetContestRecord(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	r.POST("/GetSubmitRecord", func(c *gin.Context) {
		req, err := getRequestFromGin[*proto.GetSubmitRecordRequest](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		resp, err := handler.GetSubmitRecord(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	r.POST("/GetRecentContest", func(c *gin.Context) {
		req, err := getRequestFromGin[*proto.GetRecentContestRequest](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		resp, err := handler.GetRecentContest(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	r.POST("/GetDailyQuestion", func(c *gin.Context) {
		req, err := getRequestFromGin[*proto.GetDailyQuestionRequest](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		resp, err := handler.GetDailyQuestion(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	err := r.Run(":9850")
	if err != nil {
		panic(err)
	}
}

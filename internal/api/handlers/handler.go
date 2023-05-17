package handlers

import (
	"compress/gzip"
	"fmt"
	"github.com/artemys/pprof-visualizer/internal/api/services"
	"github.com/artemys/pprof-visualizer/internal/pkg/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/proto"
	"io"
	"net/http"
)

func Healthcheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func NoRoute(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": "Not Found"})
}

func Index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}

func Visualize() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
			return
		}
		defer file.Close()

		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
			return
		}
		defer gzipReader.Close()

		// Lecture du contenu décompressé du fichier "pprof.pb"
		data, err := io.ReadAll(gzipReader)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": http.StatusBadRequest})
			return
		}

		var profile pprof.Profile
		if err := proto.Unmarshal(data, &profile); err != nil {
			fmt.Println("error reading file")
			return
		}

		c.HTML(http.StatusOK, "tree.html", services.Visualize(profile))
	}
}

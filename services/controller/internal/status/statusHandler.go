package status

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleStatus() {
	router := gin.Default()

	router.GET("/status", func(ctx *gin.Context) {
		file, err := os.ReadFile("../config/config.json")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError,
				gin.H{"error": "Error opening file"},
			)
			return
		}

		var data map[string]any
		if err := json.Unmarshal(file, &data); err != nil {
			ctx.JSON(http.StatusInternalServerError,
				gin.H{"error": "Error unparsing data"},
			)
			return
		}

		ctx.JSON(http.StatusOK, data)
	})

	router.Run(":8080")
}

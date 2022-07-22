package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/DhananjayPurohit/gin-lsat/ginlsat"
	"github.com/DhananjayPurohit/gin-lsat/ln"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func fileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func main() {
	router := gin.Default()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load .env file")
	}
	lnClient, err := ginlsat.InitLnClient(&ln.LNClientConfig{
		LNClientType: os.Getenv("LN_CLIENT_TYPE"),
		LNDConfig: ln.LNDoptions{
			Address:     os.Getenv("LND_ADDRESS"),
			MacaroonHex: os.Getenv("MACAROON_HEX"),
		},
		LNURLConfig: ln.LNURLoptions{
			Address: os.Getenv("LNURL_ADDRESS"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	lsatmiddleware, err := ginlsat.NewLsatMiddleware(&ginlsat.GinLsatMiddleware{
		Amount:   5,
		LNClient: lnClient,
	})
	if err != nil {
		log.Fatal(err)
	}

	router.Use(lsatmiddleware.Handler)

	router.GET("/:folder/:file", func(c *gin.Context) {
		lsatInfo := c.Value("LSAT").(*ginlsat.LsatInfo)
		folder := c.Param("folder")
		fileName := c.Param("file")
		if lsatInfo.Type == ginlsat.LSAT_TYPE_FREE {
			c.File(fmt.Sprintf("%s/%s", folder, fileName))
		} else if lsatInfo.Type == ginlsat.LSAT_TYPE_PAID {
			filePaidType := fileNameWithoutExt(fileName) + "-lsat" + filepath.Ext(fileName)
			if _, err := os.Stat(fmt.Sprintf("%s/%s", folder, filePaidType)); err == nil {
				c.File(fmt.Sprintf("%s/%s", folder, filePaidType))
			} else {
				c.File(fmt.Sprintf("%s/%s", folder, fileName))
			}
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"code":    http.StatusInternalServerError,
				"message": fmt.Sprint(lsatInfo.Error),
			})
		}
	})

	router.Run("localhost:8080")
}

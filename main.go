package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/getAlby/gin-lsat/ginlsat"
	"github.com/getAlby/gin-lsat/ln"
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

	router.GET("/:file", func(c *gin.Context) {
		lsatInfo := c.Value("LSAT").(*ginlsat.LsatInfo)
		fileName := c.Param("file")
		paidPath := fmt.Sprintf("assets/paid/%s", fileName)
		freePath := fmt.Sprintf("assets/free/%s", fileName)

		if lsatInfo.Type == ginlsat.LSAT_TYPE_FREE {
			c.File(freePath)
		} else if lsatInfo.Type == ginlsat.LSAT_TYPE_PAID {
			if _, err := os.Stat(paidPath); err == nil {
				c.File(paidPath)
			} else {
				c.File(freePath)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": fmt.Sprint(lsatInfo.Error),
			})
		}
	})

	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	router.Run(fmt.Sprintf(":%v", port))
}

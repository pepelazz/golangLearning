package uploadImage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const IMAGE_DIR = "./image"
const STAT_IMAGE_PATH = "/stat-img"

// генерим случаный uid для названия файла
func randomFilename() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return strings.ToLower(fmt.Sprintf("%v", ulid.MustNew(ulid.Timestamp(t), entropy)))
}

func httpError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"ok":      false,
		"message": message,
	})
	c.Abort()
}

func httpSuccess(c *gin.Context, res interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"result": res,
	})
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

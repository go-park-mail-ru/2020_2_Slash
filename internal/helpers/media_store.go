package helpers

import (
	"strconv"
	"strings"
)

func GetContentDirTitle(originalName string, cid uint64) string {
	title := strings.ReplaceAll(originalName, " ", "")
	title += "_" + strconv.FormatUint(cid, 10)
	title = strings.ToLower(title)
	return title
}

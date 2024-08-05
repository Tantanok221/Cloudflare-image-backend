package helper

import (
	"strconv"
)

func HandleParam(param string, defaultValue int) (int, error) {
	var limit int
	var err error
	if len(param) > 0 {
		limit, err = strconv.Atoi(param)
	} else {
		limit = defaultValue
	}
	return limit, err
}

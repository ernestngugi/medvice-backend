package webutils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ernestngugi/medvice-backend/internal/forms"
	"github.com/gin-gonic/gin"
)

func FilterFromContext(
	c *gin.Context,
) (*forms.Filter, error) {

	filter := &forms.Filter{}

	page, per, err := paginationFromContext(c)
	if err != nil {
		return filter, err
	}

	filter.Page = page
	filter.Per = per

	isValid := strings.TrimSpace(c.Query("valid"))
	if isValid != "" {
		isValid, err := strconv.ParseBool(isValid)
		if err != nil {
			return filter, fmt.Errorf("invalid valid argument %v", err)
		}

		filter.Valid = &isValid
	}

	return filter, nil
}

func paginationFromContext(
	c *gin.Context,
) (int, int, error) {

	page := 1
	per := 20

	var err error

	pageQuery := strings.TrimSpace(c.Query("page"))
	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			return page, per, fmt.Errorf("invalid page argument %v", err)
		}
	}

	perQuery := strings.TrimSpace(c.Query("per"))
	if perQuery != "" {
		per, err = strconv.Atoi(perQuery)
		if err != nil {
			return page, per, fmt.Errorf("invalid per argument %v", err)
		}
	}

	return page, per, nil
}

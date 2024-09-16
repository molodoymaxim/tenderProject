package limitandoffsetcheck

import (
	"fmt"
	"strconv"
)

func LimitAndOffsetCheck(limit, offset string) (limitInt int, offsetInt int, err error) {
	const op = "lib.api.LimitAndOffsetCheck"

	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			msgErr := fmt.Errorf("error in  %s. Error %s", op, err)
			return 0, 0, msgErr
		}
	} else {
		limitInt = 5
	}

	if offset != "" {
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			msgErr := fmt.Errorf("error in  %s. Error %s", op, err)
			return 0, 0, msgErr
		}
	} else {
		offsetInt = 0
	}

	if offsetInt < 0 || limitInt < 0 {
		msgErr := fmt.Errorf("limit or offset is negative", op)
		return 0, 0, msgErr
	}

	return limitInt, offsetInt, nil
}

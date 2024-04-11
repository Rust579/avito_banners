package errs

const defaultErr = 99

type Error struct {
	Key int    `json:"key"`
	Msg string `json:"message"`
}

var listErrors = map[int]string{
	99:  "unknown error",
	100: "error request read",
	101: "banner item can not be empty",
	102: "tag_id or feature_id must be empty",
	104: "unknown request",
	105: "tag_ids array can not be empty",
	106: "feature_id must be a positive integer",
	107: "tag_id must be a positive integer",
	108: "incorrect admin or user token",
	109: "token not found",
	110: "incorrect admin token",
	111: "banner already exists",
	112: "banner not found",
	113: "id must be a positive integer",
	114: "tag_id with feature_id must be a positive integer",
	115: "banners not found",
}

func GetErr(num int, str ...string) Error {
	err, ok := listErrors[num]
	if !ok {
		return Error{defaultErr, listErrors[defaultErr]}
	}
	if len(str) > 0 {
		err = err + " : " + str[0]
	}
	return Error{num, err}
}

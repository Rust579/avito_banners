package errs

const defaultErr = 99

type Error struct {
	Key int    `json:"key"`
	Msg string `json:"message"`
}

var listErrors = map[int]string{
	99:  "unknown error",
	100: "error request read",
	101: "text can not be empty",
	102: "title can not be empty",
	103: "url can not be empty",
	104: "unknown request",
	105: "tag_ids can not be empty",
	106: "feature_id can not be zero or negative",
	400: "Некорректные данные",
	401: "Пользователь не авторизован",
	403: "Пользователь не имеет доступа",
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
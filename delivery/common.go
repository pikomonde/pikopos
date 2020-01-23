package delivery

type responseAPI struct {
	Status      int         `json:"status"`
	ProcessTime int         `json:"process_time"`
	Data        interface{} `json:"data"`
}

const errorWrongJSONFormat = "Wrong JSON Format"

// config.JWTSecret

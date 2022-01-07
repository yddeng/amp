package exec

type Application struct {
	AppId    int   `json:"app_id"`
	Pid      int   `json:"pid"`
	CreateAt int64 `json:"create_at"`
}

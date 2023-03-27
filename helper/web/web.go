package web

type BaseResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type BaseMeta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DataUpdate struct {
	RowsAffected int64 `json:"rows_affected"`
}

type MetaPagination struct {
	BaseMeta
	Page      int `json:"page"`
	Size      int `json:"size"`
	TotalPage int `json:"totalPage"`
	TotalData int `json:"totalData"`
}
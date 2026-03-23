package rpc

type Meta struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func NewMeta(total, offset, limit int) *Meta {
	return &Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
}

type Response struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
	Meta    *Meta  `json:"meta,omitempty"`

	// Deprecated legacy fields
	Status int      `json:"status,omitempty"`
	Error  string   `json:"error,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

func NewResponse(isSuccess bool) *Response {
	return &Response{
		Success: isSuccess,
	}
}

func (r *Response) SetData(data any) *Response {
	r.Data = data

	return r
}

func (r *Response) SetMeta(total, page, limit int) *Response {
	r.Meta = NewMeta(total, page, limit)

	return r
}

func (r *Response) SetMessage(message string) *Response {
	r.Message = message

	return r
}

func (r *Response) SetError(err error) *Response {
	if err != nil {
		errMsg := err.Error()
		if r.Message == "" {
			r.Message = errMsg
		}
		r.Error = errMsg
	}

	return r
}

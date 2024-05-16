package entity

type Response struct {
	id     string
	status string
}

func NewResponse(id string, status string) *Response {
	return &Response{
		id:     id,
		status: status,
	}
}

func (res *Response) GetId() string {
	return res.id
}

func (res *Response) GetStatus() string {
	return res.status
}

func (res *Response) SetId(id string) {
	res.id = id
}

func (res *Response) SetStatus(status string) {
	res.status = status
}

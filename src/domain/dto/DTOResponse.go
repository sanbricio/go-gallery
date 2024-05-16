package dto

import entity "api-upload-photos/src/domain/entities"

type DTOResponse struct {
	Id     string
	Status string
}

func FromResponse(res *entity.Response) *DTOResponse {
	return &DTOResponse{
		Id:     res.GetId(),
		Status: res.GetStatus(),
	}
}

func (res *DTOResponse) AsResponseEntity() *entity.Response {
	return entity.NewResponse(res.Id, res.Status)
}

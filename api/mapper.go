package api

import "github.com/bipuldutta/blogzilla/domain"

/*
	This is where all the conversion between API->Domain objects and the Domain->API
*/

func convertCreateBlogRequestToDomain(userID int64, request *CreateBlogRequestV1) *domain.Blog {
	return &domain.Blog{
		UserID:  userID,
		Title:   request.Title,
		Content: request.Content,
		Tags:    request.Tags,
	}
}

func convertCreateUserRequestToDomain(request *CreateUserRequestV1) *domain.User {
	return &domain.User{
		Username:  request.Username,
		Password:  request.Password,
		FirstName: request.FirstName,
		LastName:  request.LastName,
	}
}

func convertUserDomainObjToAPI(dom *domain.User) *UserResponseV1 {
	return &UserResponseV1{
		ID:        dom.ID,
		Username:  dom.Username,
		FirstName: dom.FirstName,
		LastName:  dom.LastName,
		CreatedAt: dom.CreatedAt,
		UpdatedAt: dom.UpdatedAt,
	}
}

package response

type ResponseRemove struct {
	Id         int  `json:"id"`
	ProjectsId int  `json:"projectsId"`
	Removed    bool `json:"removed"`
}
type RequestCreate struct {
	Name string `json:"name" binding:"required"`
}

type RequestUpdate struct {
	RequestCreate
	Description string `json:"description,omitempty"`
}

package dto

type ObjectIDsDto struct {
	ObjectIDs []ObjectID `json:"object_ids"`
}

type ObjectID struct {
	ObjID        string `json:"obj_id"`
	FileCategory string `json:"file_category"`
}

package db

type Elector struct {
	ID                     int64   `bson:"id" json:"id" bindings:"required"`
	FullName               string  `bson:"fullname" json:"fullname"`
	Email                  string  `bson:"email" json:"email" bindings:"required"`
	IsDormitoryStudent     bool    `bson:"isDormitoryStudent" json:"isDormitoryStudent" bindings:"required"`
	DormitoryId            int64   `bson:"dormitoryId" json:"dormitoryId"`
	FacultyIds             []int64 `bson:"facultyIds" json:"facultyIds" bindings:"required"`
	IsCouncil              bool    `bson:"isCouncil" json:"isCouncil" bindings:"required"`
	CouncilOrganizationIds []int64 `bson:"councilOrganizationIds" json:"councilOrganizationIds"`
	IsPostGraduate         bool    `bson:"isPostGraduate" json:"isPostGraduate" bindings:"required"`
	IsNearForeign          bool    `bson:"isNearForeign" json:"isNearForeign" bindings:"required"`
	IsFarForeign           bool    `bson:"isFarForeign" json:"isFarForeign" bindings:"required"`
}

type NameTag struct {
	ID   int64  `bson:"id" json:"id" bindings:"required"`
	Name string `bson:"name" json:"name" bindings:"required"`
}

type Faculty struct {
	NameTag
}

type Dormitory struct {
	NameTag
}

type CouncilOrganization struct {
	NameTag
}

type Election struct {
	NameTag
	Priority                        int64   `bson:"priority" json:"priority" bindings:"required"`
	IsRunoff                        bool    `bson:"isRunoff" json:"isRunoff" bindings:"required"`
	Mandates                        int64   `bson:"mandates" json:"mandates" bindings:"required"`
	IsForNearForeign                bool    `bson:"isForNearForeign" json:"isForNearForeign" bindings:"required"`
	IsForFarForeign                 bool    `bson:"isForFarForeign" json:"isForFarForeign" bindings:"required"`
	AcceptedCouncilOrganizationsIds []int64 `bson:"acceptedCouncilOrganizationsIds" json:"acceptedCouncilOrganizationsIds"`
	AcceptedFacultiesIds            []int64 `bson:"acceptedFacultiesIds" json:"acceptedFacultiesIds"`
	AcceptedDormitoriesIds          []int64 `bson:"acceptedDormitoriesIds" json:"acceptedDormitoriesIds"`
	StartTime                       string  `bson:"startTime" json:"startTime" bindings:"required"`
	FinishTime                      string  `bson:"finishTime" json:"finishTime"`
	Status                          string  `bson:"status" json:"status" bindings:"required"`
}

type ElectionResults struct {
	// TODO
}

type CandidateRequest struct {
	Name        string `bson:"name" json:"name" bindings:"required"`
	PhotoUrl    string `bson:"photourl" json:"photourl" bindings:"required"`
	Description string `bson:"description" json:"description" bindings:"required"`
}

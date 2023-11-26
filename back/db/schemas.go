package db

type Elector struct {
	ID                     string   `bson:"id" json:"id" bindings:"required"`
	FullName               string   `bson:"fullname" json:"fullname"`
	Email                  string   `bson:"email" json:"email" bindings:"required"`
	IsDormitoryStudent     bool     `bson:"isDormitoryStudent" json:"isDormitoryStudent" bindings:"required"`
	DormitoryId            string   `bson:"dormitoryId" json:"dormitoryId"`
	FacultyIds             []string `bson:"facultyIds" json:"facultyIds" bindings:"required"`
	IsCouncil              bool     `bson:"isCouncil" json:"isCouncil" bindings:"required"`
	CouncilOrganizationIds []string `bson:"councilOrganizationIds" json:"councilOrganizationIds"`
	IsPostGraduate         bool     `bson:"isPostGraduate" json:"isPostGraduate" bindings:"required"`
	IsNearForeign          bool     `bson:"isNearForeign" json:"isNearForeign" bindings:"required"`
	IsFarForeign           bool     `bson:"isFarForeign" json:"isFarForeign" bindings:"required"`
}

type Faculty struct {
	ID   string `bson:"id" json:"id" bindings:"required"`
	Name string `bson:"name" json:"name" bindings:"required"`
}

type Dormitory struct {
	ID   string `bson:"id" json:"id" bindings:"required"`
	Name string `bson:"name" json:"name" bindings:"required"`
}

type CouncilOrganization struct {
	ID   string `bson:"id" json:"id" bindings:"required"`
	Name string `bson:"name" json:"name" bindings:"required"`
}

type Election struct {
	ID                              string  `bson:"id" json:"id" bindings:"required"`
	Name                            string  `bson:"name" json:"name" bindings:"required"`
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
	IsVoted                         bool    `json:"IsVoted"`
}

type ElectionResults struct {
	// TODO
}

type Candidate struct {
	ID                string `bson:"id" json:"id"`
	ElectionId        string `bson:"electionId" json:"electionId"`
	Email             string `bson:"email" json:"email"` // This field is filled only for admin.
	Name              string `bson:"name" json:"name" bindings:"required"`
	PhotoUrl          string `bson:"photourl" json:"photourl" bindings:"required"`
	Description       string `bson:"description" json:"description" bindings:"required"`
	Approved          bool   `bson:"approved" json:"approved"`
	WaitingForApprove bool   `bson:"waitingForApprove" json:"waitingForApprove"`
}

type ElectionId struct {
	ID string `bson:"id" json:"id" uri:"electionId" bindings:"required"`
}

type CandidateId struct {
	ID string `bson:"id" json:"id" uri:"candidateId" bindings:"required"`
}

type PublicKey struct {
	Key string `bson:"key" json:"key" bindings:"required"`
	ID  string `bson:"id" json:"id"`
}

type PrivateKey struct {
	Key string `bson:"key" json:"key" bindings:"required"`
	ID  string `bson:"id" json:"id"`
}

type EncryptedVote struct {
	Vote    string `bson:"vote" json:"vote" bindings:"required"`
	VoterID string `bson:"voterId"`
}

type DecryptedVote struct {
	Vote string `bson:"vote" json:"vote" bindings:"required"`
}

type UserId struct {
	Email string `bson:"email"`
	ID    string `bson:"id"`
}

const (
	Draft     = "draft"
	Created   = "created"
	Waiting   = "waiting"
	Started   = "started"
	Voted     = "voted" // Not used in Statuses, as this status is different for different users
	Finished  = "finished"
	Decrypted = "decrypted"
	Results   = "results"
)

var Statuses = []string{
	Draft,
	Created,
	Waiting,
	Started,
	Finished,
	Decrypted,
	Results,
}

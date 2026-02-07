package climbs

type ClimbModel struct {
	PK                  string `json:"-" dynamodbav:"PK"`
	SK                  string `json:"-" dynamodbav:"SK"`
	EntityType          string `json:"-" dynamodbav:"EntityType"`
	ID                  string `json:"id" dynamodbav:"ID"`
	UserId              string `json:"userId" dynamodbav:"UserId"`
	Grade               string `json:"grade" dynamodbav:"Grade"`
	PerceivedDifficulty string `json:"perceivedDifficulty" dynamodbav:"PerceivedDifficulty"`
	TypeOfClimb         string `json:"typeOfClimb" dynamodbav:"TypeOfClimb"`
	GradeSystem         string `json:"gradeSystem" dynamodbav:"GradeSystem"`
	CreatedAt           int64  `json:"createdAt" dynamodbav:"CreatedAt"`
	UpdatedAt           int64  `json:"updatedAt" dynamodbav:"UpdatedAt"`
}

type ClimbUpdateReq struct {
	Grade               string `json:"grade" dynamodbav:"Grade"`
	PerceivedDifficulty string `json:"perceivedDifficulty" dynamodbav:"PerceivedDifficulty"`
	TypeOfClimb         string `json:"typeOfClimb" dynamodbav:"TypeOfClimb"`
	GradeSystem         string `json:"gradeSystem" dynamodbav:"GradeSystem"`
}

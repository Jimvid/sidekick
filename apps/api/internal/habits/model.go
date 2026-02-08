package habits

type habitItem struct {
	UserId string `json:"-" dynamodbav:"UserId"` // Used as primary key
	ItemId string `json:"-" dynamodbav:"ItemId"` // used for sorting key
	HabitModel
}

type HabitModel struct {
	ID          string `json:"id" dynamodbav:"ID"`
	Name        string `json:"name" dynamodbav:"Name"`
	Description string `json:"description" dynamodbav:"Description"`
	Color       string `json:"color" dynamodbav:"Color"`
	CreatedAt   int64  `json:"createdAt" dynamodbav:"CreatedAt"`
	UpdatedAt   int64  `json:"updatedAt" dynamodbav:"UpdatedAt"`
}

type HabitReq struct {
	Name        string `json:"name" dynamodbav:"Name"`
	Description string `json:"description" dynamodbav:"Description"`
	Color       string `json:"color" dynamodbav:"Color"`
}

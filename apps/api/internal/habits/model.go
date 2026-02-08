package habits

type habitItem struct {
	UserId string `json:"-" dynamodbav:"userId"` // Used as primary key
	ItemId string `json:"-" dynamodbav:"itemId"` // used for sorting key
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

type habitLogItem struct {
	UserId string `json:"-" dynamodbav:"userId"`
	ItemId string `json:"-" dynamodbav:"itemId"`
	HabitLogModel
}

type HabitLogModel struct {
	ID        string `json:"id" dynamodbav:"ID"`
	HabitId   string `json:"habitId" dynamodbav:"HabitId"`
	Date      string `json:"date" dynamodbav:"Date"`
	Note      string `json:"note" dynamodbav:"Note"`
	CreatedAt int64  `json:"createdAt" dynamodbav:"CreatedAt"`
	UpdatedAt int64  `json:"updatedAt" dynamodbav:"UpdatedAt"`
}

type HabitLogReq struct {
	HabitId string `json:"habitId"`
	Date    string `json:"date"`
	Note    string `json:"note"`
}

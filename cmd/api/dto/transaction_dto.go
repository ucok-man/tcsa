package dto

type TransactionCreateDTO struct {
	UserId int `json:"user_id" validate:"required,min=1"`
	Amount int `json:"amount" validate:"required,min=0"`
}

type TransactionUpdateDTO struct {
	Amount *int    `json:"amount" validate:"omitempty,min=0"`
	Status *string `json:"status" validate:"omitempty,oneof=pending failed success"`
}

type TransactionGetAllDTO struct {
	Pagination struct {
		Page     int `validate:"omitempty,min=1,max=1000"`
		PageSize int `validate:"omitempty,min=1,max=100"`
		Limit    *int
		Offset   *int
	}
	Sort struct {
		RawValue    string `validate:"omitempty,oneof=amount status created_at -amount -status -created_at ''"`
		ColumnValue *string
		Direction   *string
	}
	Filter struct {
		Status string `validate:"omitempty,oneof=pending failed success ''"`
	}
}

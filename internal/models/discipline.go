package models

type Discipline struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	Credits        int    `json:"credits"`
	WorkloadHours  int    `json:"workload_hours"`
	Description    string `json:"description"`
	DepartmentID   int    `json:"department_id"`
	DepartmentName string `json:"department_name"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

package models

import "fmt"

type DisciplineOffer struct {
	ID             int    `json:"id"`
	DisciplineID   int    `json:"discipline_id"`
	SemeterID      int    `json:"semester_id"`
	TeacherID      int    `json:"teacher_id"`
	Schedule       string `json:"schedule"`
	ClassCode      string `json:"class_code"`
	DisciplineName string `json:"discipline_name,omitempty"`
	SemesterLabel  string `json:"semester_label,omitempty"`
	TeacherName    string `json:"teacher_name,omitempty"`
}

func (o *DisciplineOffer) SetSemesterLabel(year, period int) {
	o.SemesterLabel = fmt.Sprintf("%d.%d", year, period)
}

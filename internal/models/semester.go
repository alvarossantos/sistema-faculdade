package models

import "fmt"

type AcademicSemester struct {
	ID     int `json:"id"`
	Year   int `json:"year"`
	Period int `json:"period"`
}

func (s *AcademicSemester) String() string {
	return fmt.Sprintf("%d.%d", s.Year, s.Period)
}

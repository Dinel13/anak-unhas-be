package helper

import (
	"strings"
)

func IsDomainUnhas(email string) bool {
	domain := strings.Split(email, "@")
	return domain[len(domain)-1] == "student.unhas.ac.id"
}

package utils

import (
	"strconv"

	"github.com/bookpanda/mygraderlist-auth/src/constant/utils"
	"github.com/pkg/errors"
)

const CurrentYear = 65

func GetFacultyFromID(sid string) (*utils.Faculty, error) {
	if len(sid) != 10 {
		return nil, errors.New("Invalid faculty id")
	}

	result, ok := utils.Faculties[sid[8:10]]
	if !ok {
		return nil, errors.New("Invalid faculty id")
	}
	return &result, nil
}

func CalYearFromID(sid string) (string, error) {
	if len(sid) != 10 {
		return "", errors.New("Invalid student id")
	}

	yearIn, err := strconv.Atoi(sid[:2])
	if err != nil {
		return "", errors.New("Invalid student id")
	}

	studYear := CurrentYear - yearIn + 1
	if studYear <= 0 {
		return "", errors.New("Invalid student ID")
	}

	return strconv.Itoa(studYear), nil
}

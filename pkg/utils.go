package utils

import (
	"crypto/sha512"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"unicode/utf8"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const (
	FilmTitleBegin       = 1
	FilmTitleEnd         = 150
	FilmDescriptionBegin = 1
	FilmDescriptionEnd   = 1000
	FilmRatingBegin      = 0
	FilmRatingEnd        = 10
	ActorNameBegin       = 1
	ActorNameEnd         = 150
)

func HashPassword(password string) []byte {
	hashPassword := sha512.Sum512([]byte(password))
	passwordByteSlice := hashPassword[:]
	return passwordByteSlice
}

func RandStringRunes(seed int) string {
	symbols := make([]rune, seed)
	for i := range symbols {
		symbols[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(symbols)
}

func ValidateStringSize(validatedString string, begin int, end int, validateError string, logger *logrus.Logger) error {
	validateStringLength := utf8.RuneCountInString(validatedString)
	if validateStringLength > end || validateStringLength < begin {
		logger.Error(validateError)
		return fmt.Errorf(validateError)
	}
	return nil
}

const (
	InvalidEmailOrPasswordError     = "Invalid email or password"
	SessionRepositoryNotActiveError = "Session repository not active"
	ProfileRepositoryNotActiveError = "Profile repository not active"
	CreateProfileError              = "Create profile failed"
	ProfileNotFoundError            = "Profile not found"
	GetProfileError                 = "Get profile failed"
	GetProfileRoleError             = "Get profile role failed"
	RatingSizeError                 = "Rating must be from 0 to 10"
	TitleSizeError                  = "Title size must be from 1 to 150"
	DescriptionSizeError            = "Description size must be from 1 to 1000"
	FilmsListNotFoundError          = "Films list not found"
	ActorNameSizeError              = "Actor name size must be from 1 to 150"
	GrpcRecievError                 = "gRPC recieve error"
)

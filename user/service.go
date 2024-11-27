package user

import (
	"encoding/json"
	"go-clean-architecture/kafka"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUserInput(input RegisterUserInput) (User, error)
	RegisterUsersInput(input RegisterUsersInput) ([]User, error)
}

type service struct {
	repository   Repository
	kafkaService KafkaService
}

func NewService(repository Repository, kafkaService KafkaService) *service {
	return &service{
		repository:   repository,
		kafkaService: kafkaService,
	}
}

type KafkaService interface {
	SendToKafka(data interface{}) error
}

type kafkaService struct {
	producer *kafka.KafkaProducer
}

func NewKafkaService(producer *kafka.KafkaProducer) KafkaService {
	return &kafkaService{producer: producer}
}

func (k *kafkaService) SendToKafka(data interface{}) error {
	message, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return k.producer.SendMessage(string(message))
}

func (s *service) RegisterUsersInput(input RegisterUsersInput) ([]User, error) {
	var users []User

	for _, userInput := range input.Users {
		user := User{
			Name:       userInput.Name,
			Occupation: userInput.Occupation,
			Email:      userInput.Email,
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.MinCost)
		if err != nil {
			return users, err
		}
		user.PasswordHash = string(passwordHash)
		user.Role = "user"

		newUser, err := s.repository.Save(user)
		if err != nil {
			return users, err
		}

		users = append(users, newUser)
	}

	err := s.kafkaService.SendToKafka(users) // s.kafkaService undefined (type *service has no field or method kafkaService)
	if err != nil {
		log.Printf("Failed to send data to Kafka: %s", err)
		return users, err
	}

	return users, nil
}

func (s *service) RegisterUserInput(input RegisterUserInput) (User, error) {
	user := User{
		Name:       input.Name,
		Occupation: input.Occupation,
		Email:      input.Email,
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	// Save user to database
	newUser, err := s.repository.Save(user)
	if err != nil {
		return user, err
	}

	return newUser, nil
}

package person_repository_test

import (
	"github.com/mrminglang/go-rigger/boot"
	"github.com/mrminglang/go-rigger/repositories/models"
	"github.com/mrminglang/go-rigger/repositories/mongo/person_repository"
	genid "github.com/srlemon/gen-id"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	boot.Boot()
	m.Run()
}

func TestPerson_CreatePerson(t *testing.T) {
	person := &models.Person{
		Name:     genid.NewGeneratorData().Name,
		Age:      10,
		Email:    genid.NewGeneratorData().Email,
		CreateAt: time.Now(),
	}
	err := person_repository.NewPerson().CreatePerson(person)
	assert.Nil(t, err)

	boot.Destroy()
}

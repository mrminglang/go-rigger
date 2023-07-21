package news_main_repository_test

import (
	"github.com/mrminglang/go-rigger/boot"
	"github.com/mrminglang/go-rigger/repositories/mysql/news_main_repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M) {
	boot.Boot()
	m.Run()
}

func TestQueryRecord(t *testing.T) {
	sSql := "SELECT * FROM news_main_202305 WHERE news_id = 'Kn-172411';"
	var vMysqlData []map[string]string
	err := news_main_repository.QueryRecord(sSql, &vMysqlData)
	if err != nil {
		assert.Error(t, err)
	}
}

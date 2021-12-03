package utils_test

// import (
// 	"github.com/DATA-DOG/go-sqlmock"
// 	"gitlab.com/www.monitordomercado.com.br/mm/src/database"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func MockDB() (sqlmock.Sqlmock, error) {
// 	mockSQLDB, sqlMock, err := sqlmock.New()
// 	if err != nil {
// 		return sqlMock, err
// 	}

// 	gormDB, err := gorm.Open(mysql.New(mysql.Config{
// 		Conn:                      mockSQLDB,
// 		SkipInitializeWithVersion: true,
// 	}), &gorm.Config{})
// 	if err != nil {
// 		return sqlMock, err
// 	}

// 	database.DB = gormDB
// 	return sqlMock, nil
// }

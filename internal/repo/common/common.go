package common

import (
	"github.com/kargotech/gokargo/unitofwork/consistency"
	"gorm.io/gorm"
)

func DecideDBTxn(existingTxn *gorm.DB, cs consistency.ConsistencyItf) *gorm.DB {

	if cs != nil && cs.GetDBTxn() != nil {
		return cs.GetDBTxn()
	}
	return existingTxn
}

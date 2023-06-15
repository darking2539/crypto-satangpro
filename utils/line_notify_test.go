package utils

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestUpdateOneCollectCpeFlagReturnEquipmentAllDataInArrayTransactionDAO(t *testing.T) {

	godotenv.Load("../.env")
	LineNotify("boss")

	t.Errorf("PrintALL DATA")
}
package api

import (
	"crypto-satangpro/models"
	"crypto-satangpro/repositories"
	"encoding/json"
	"strings"
	"time"

	"github.com/darahayes/go-boom"
	"github.com/gin-gonic/gin"
)

func GetTransactionListService(c *gin.Context) {

	var request models.GetListRequest
	if payloadErr := c.ShouldBindJSON(&request); payloadErr != nil {
		boom.BadRequest(c.Writer, payloadErr.Error());
		return
	}
	
	lowerAddress := strings.ToLower(request.Address)

	transactionList, paigination, repoErr := repositories.GetTransactionListRepo(request.Page, request.PerPage, lowerAddress)
	if repoErr != nil {
		boom.BadRequest(c.Writer, repoErr.Error());
		return
	}
	
	dataResp := []models.TransactionData{}
	for _, tran := range transactionList {
		dataResp = append(dataResp, models.TransactionData{
			BlockNo: tran.BlockNo,
			TransactionIndex: tran.TransactionIndex,
			Hash: tran.Hash,
			From: tran.From,
			To: tran.To,
			Value: tran.Value,
			Gas: tran.Gas,
			GasPrice: tran.GasPrice,
		})
	}
	
	resp := models.GetListResponse{
		Data: dataResp,
		Pagination: paigination,
	}
	
	c.Writer.Header().Set("Content-Type", "application/json");
	json.NewEncoder(c.Writer).Encode(&resp);
}

func AddAddressService(c *gin.Context) {

	var request models.AddUserRequest
	if payloadErr := c.ShouldBindJSON(&request); payloadErr != nil {
		boom.BadRequest(c.Writer, payloadErr.Error());
		return
	}
	
	lowerAddress := strings.ToLower(request.Address)

	dao := models.UserModel{
		Address: lowerAddress,
		CreatedBy: request.Email,
		CreatedDate: time.Now(),
	}

	response, err := repositories.CreateUserRepo(dao)
	if err != nil {
		boom.BadRequest(c.Writer, err.Error());
		return
	}

	resp := response
	
	c.Writer.Header().Set("Content-Type", "application/json");
	json.NewEncoder(c.Writer).Encode(&resp);
}
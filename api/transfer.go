package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/arya2004/xyfin/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {

	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//
	_, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
		
	if !valid{
		return
	}

	


	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid{
		return
	}

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountID,
		ToAccountId:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)

}


func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountId)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
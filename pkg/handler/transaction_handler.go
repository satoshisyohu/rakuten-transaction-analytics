package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/handler/dto"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/infrastructure"
	"net/http"
)

type TransactionHandler struct {
	// usecase
	u infrastructure.IAnalyticsTransactionUsecase
}

func NewTransactionHandler(u infrastructure.IAnalyticsTransactionUsecase) *TransactionHandler {
	return &TransactionHandler{u: u}
}

func (h *TransactionHandler) Analysis(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "setting", "a")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// フォームデータの解析
	err := r.ParseMultipartForm(10 << 20) // 最大10MB
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	// metadataフィールドからJSON文字列を取得
	metadataStr := r.FormValue("metadata")
	var req dto.TransactionRequest
	err = json.Unmarshal([]byte(metadataStr), &req)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid JSON in metadata field", http.StatusBadRequest)
		return
	}
	var res *dto.TransactionResponse
	if res, err = h.u.Run(ctx, req, files); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// response書き込み
	helper.WriteJsonResponse(w, http.StatusOK, res)

	return
}

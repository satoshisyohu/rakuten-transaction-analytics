package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/handler/dto"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/helper"
	"github.com/satoshisyohu/rakuten-transaction-analytics/pkg/infrastructure"
)

type TransactionHandler struct {
	// usecase
	atu infrastructure.IAnalyticsTransactionUsecase
	rtu infrastructure.IRetrieveTransactionUsecase
}

func NewTransactionHandler(atu infrastructure.IAnalyticsTransactionUsecase,
	rtu infrastructure.IRetrieveTransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		atu: atu,
		rtu: rtu,
	}
}

func (h *TransactionHandler) Analysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	if res, err = h.atu.Run(ctx, req, files); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// response書き込み
	helper.WriteJsonResponse(w, http.StatusOK, res)

}

func (h *TransactionHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	var res *dto.RetrieveTransactionResponses

	if res, err = h.rtu.Run(ctx, files); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// response書き込み
	helper.WriteJsonResponse(w, http.StatusOK, res)

}

func (h *TransactionHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var (
		req dto.TransactionRequest
		res *dto.TransactionResponse
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}
	if res, err = h.atu.Execute(ctx, req); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(res)
	// response書き込み
	helper.WriteJsonResponse(w, http.StatusOK, res)

}

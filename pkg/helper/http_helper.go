package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJsonResponse(writer http.ResponseWriter, httpStatus int, response any) {

	// status設定
	writer.WriteHeader(httpStatus)
	// json
	writer.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		http.Error(writer, "Failed to create JSON response", http.StatusInternalServerError)
		return
	}
	_, err = writer.Write(res)
	if err != nil {
		return
	}

}

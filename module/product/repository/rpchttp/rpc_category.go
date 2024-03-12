package rpchttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"log"
	"my-app/module/product/query"
	"net/http"
	"time"
)

type rpcGetCategoriesByIds struct {
	url string
}

func NewRpcGetCategoriesByIds(url string) rpcGetCategoriesByIds {
	return rpcGetCategoriesByIds{url: url}
}

func (rpc rpcGetCategoriesByIds) FindWithIds(ctx context.Context, ids []uuid.UUID) ([]query.CategoryDTO, error) {
	url := rpc.url
	method := "GET"

	data := struct {
		Ids []uuid.UUID `json:"ids"`
	}{
		ids,
	}

	dataByte, _ := json.Marshal(data)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(dataByte))

	if err != nil {
		log.Println(err)
		return nil, errors.New("cannot get categories")
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return nil, errors.New("cannot get categories")
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
		return nil, errors.New("cannot get categories")
	}

	var responseData struct {
		Data []query.CategoryDTO `json:"data"`
	}

	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Println(err)
		return nil, errors.New("cannot get categories")
	}

	return responseData.Data, nil
}

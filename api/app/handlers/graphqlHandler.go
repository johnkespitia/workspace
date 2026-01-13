package handlers

import "net/http"

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func GraphQLHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: decode JSON body into GraphQLRequest struct
	// TODO: if Query is categories -> return list
	// TODO: if Mutation is createCategory -> add item + return it
	// TODO: send JSON response in {"data": {...}} format
}

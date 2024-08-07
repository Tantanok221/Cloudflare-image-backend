package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/tantanok221/cloudflare-image-backend/db"
	"github.com/tantanok221/cloudflare-image-backend/internal/helper"
	"github.com/tantanok221/cloudflare-image-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

func PostImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		db, closeDB := db.Init()
		defer cancel()
		defer closeDB()
		image, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer image.Close()
		contentType := header.Header.Get("Content-Type")
		if !utils.IsImage(contentType) {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}
		name := r.FormValue("name")
		if len(name) == 0 {
			http.Error(w, "Lack of name", http.StatusBadRequest)
			return
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", header.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(part, image)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = writer.Close()
		ir, err := http.NewRequest("POST", utils.GetEndpoint(), body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		client := &http.Client{}

		ir.Header.Set("Content-Type", writer.FormDataContentType())
		ir.Header.Set("Authorization", utils.GetAuthToken())

		res, err := client.Do(ir)
		if err != nil {
			http.Error(w, "Error on uploading to cloudflare image"+err.Error(), http.StatusInternalServerError)
		}
		defer res.Body.Close()
		var cloudflareResponse struct {
			Result struct {
				Variants []string `json:"variants"`
			} `json:"result"`
		}
		if err := json.NewDecoder(res.Body).Decode(&cloudflareResponse); err != nil {
			http.Error(w, "Error decoding Cloudflare response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if len(cloudflareResponse.Result.Variants) == 0 {
			http.Error(w, "No image variants returned from Cloudflare", http.StatusInternalServerError)
			return
		}
		imageLink := cloudflareResponse.Result.Variants[0]
		data := utils.Path{
			ImageLink: imageLink,
			Name:      name,
			Paths:     make(primitive.M),
		}
		helper.WriteJSON(w, data)
		_, err = db.InsertOne(ctx, data)
		if err != nil {
			http.Error(w, "Insert MongoDB went wrong: \n"+err.Error(), http.StatusInternalServerError)
		}
	}
}

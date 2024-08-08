package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tantanok221/cloudflare-image-backend/db"
	"github.com/tantanok221/cloudflare-image-backend/internal/helper"
	"github.com/tantanok221/cloudflare-image-backend/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type cloudflareResponse *struct {
	Result struct {
		Variants []string `json:"variants"`
	} `json:"result"`
}

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
		image.Close()
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

		body, writer, err := WriteFormData(image, header)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response, err := ReturnCloudflareImageResponse(body, writer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		imageLink := response.Result.Variants[0]
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

func WriteFormData(image multipart.File, header *multipart.FileHeader) (*bytes.Buffer, *multipart.Writer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", header.Filename)
	print("Writing Form data \n")
	if err != nil {
		return nil, writer, err
	}
	_, err = io.Copy(part, image)
	if err != nil {
		return nil, writer, err
	}
	err = writer.Close()
	return body, writer, nil
}

func PostCloudflareImageAPI(body *bytes.Buffer, writer *multipart.Writer) (*http.Response, error) {
	ir, err := http.NewRequest("POST", utils.GetEndpoint(), body)
	print("Posting API Data \n")
	if err != nil {
		return nil, err
	}
	client := &http.Client{}

	ir.Header.Set("Content-Type", writer.FormDataContentType())
	ir.Header.Set("Authorization", utils.GetAuthToken())
	fmt.Printf("%v", ir.Header)
	res, err := client.Do(ir)
	return res, err
}

func ReturnCloudflareImageResponse(body *bytes.Buffer, writer *multipart.Writer) (cloudflareResponse, error) {
	res, err := PostCloudflareImageAPI(body, writer)
	var response cloudflareResponse
	if err != nil {
		return nil, errors.New("Error on uploading to cloudflare image" + err.Error())
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.New("Error decoding Cloudflare response: " + err.Error())
	}

	if len(response.Result.Variants) == 0 {
		return nil, errors.New("No image variants returned from Cloudflare")
	}
	err = writer.Close()
	return response, nil
}

package handler

import (
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"time"

	"github.com/goamz/goamz/s3"
	"github.com/guregu/kami"
	"github.com/satori/go.uuid"
	"gopkg.in/gorp.v1"

	"encoding/json"
	"github.com/keichi/scribble/auth"
	"github.com/keichi/scribble/model"
)

const (
	imageSizeLimit = 3 * 1024 * 1024 * 1024
)

func AddImage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.AuthContext)
	bucket := ctx.Value("s3").(*s3.Bucket)

	var ownerId int64
	if auth.IsLoggedIn {
		ownerId = auth.User.Id
	}

	image := &model.Image{
		Id:          0,
		ContentType: r.Header.Get("Content-Type"),
		Uuid:        uuid.NewV4().String(),
		OwnerId:     ownerId,
		ShareState:  model.SHARE_STATE_PUBLIC,
		CreatedAt:   time.Now().UnixNano(),
		UpdatedAt:   time.Now().UnixNano(),
	}

	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		resp := ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Content-Length format is invalid: %v", err),
		}
		resp.Render(w)
		return
	}
	if length > imageSizeLimit {
		resp := ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Image size is too large: %v", err),
		}
		resp.Render(w)
		return
	}

	err = bucket.PutReader(image.Uuid, r.Body, int64(length),
		image.ContentType, s3.PublicRead, s3.Options{})
	defer r.Body.Close()

	if err != nil {
		resp := ErrorResponse{
			http.StatusInternalServerError,
			fmt.Sprintf("Could not store image: %v", err),
		}
		resp.Render(w)
		return
	}

	err = db.Insert(image)
	if err != nil {
		resp := ErrorResponse{
			http.StatusInternalServerError,
			fmt.Sprintf("Insert failed: %v", err),
		}
		resp.Render(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(image)
}

func getImage(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.AuthContext)
	imageId, err := strconv.Atoi(kami.Param(ctx, "imageId"))

	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Invalid image id format: %v", err),
		}
	}

	image := new(model.Image)
	err = db.SelectOne(image, "select * from images where id = ?", imageId)
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Query failed: %v", err),
		}
	}

	if !image.Authorize(auth.User, model.ACTION_READ) {
		return nil, &ErrorResponse{
			http.StatusUnauthorized,
			"Unauthorized action",
		}
	}

	return image, nil
}

var GetImage = WrapJsonHandler(nil, getImage)

func updateImage(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.AuthContext)
	newImage := req.(*model.Image)
	imageId, err := strconv.Atoi(kami.Param(ctx, "imageId"))

	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Invalid image id format: %v", err),
		}
	}

	image := new(model.Image)
	err = db.SelectOne(image, "select * from images where id = ?", imageId)
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Query failed: %v", err),
		}
	}

	if !image.Authorize(auth.User, model.ACTION_UPDATE) {
		return nil, &ErrorResponse{
			http.StatusUnauthorized,
			"Unauthorized action",
		}
	}

	image.ContentType = newImage.ContentType
	image.Uuid = newImage.Uuid
	image.OwnerId = newImage.OwnerId
	image.NoteId = newImage.NoteId
	image.ShareState = newImage.ShareState
	image.UpdatedAt = time.Now().UnixNano()

	if _, err := db.Update(image); err != nil {
		return nil, &ErrorResponse{
			http.StatusInternalServerError,
			fmt.Sprintf("Update failed: %v", err),
		}
	}

	return image, nil
}

var UpdateImage = WrapJsonHandler(model.Image{}, updateImage)

func deleteImage(ctx context.Context, req interface{}) (interface{}, *ErrorResponse) {
	db := ctx.Value("db").(*gorp.DbMap)
	auth := ctx.Value("auth").(*auth.AuthContext)
	bucket := ctx.Value("s3").(*s3.Bucket)
	imageId, err := strconv.Atoi(kami.Param(ctx, "imageId"))

	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Invalid note id format: %v", err),
		}
	}

	image := new(model.Image)
	err = db.SelectOne(image, "select * from images where id = ?", imageId)
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusBadRequest,
			fmt.Sprintf("Query failed: %v", err),
		}
	}

	if !image.Authorize(auth.User, model.ACTION_DELETE) {
		return nil, &ErrorResponse{
			http.StatusUnauthorized,
			"Unauthorized action",
		}
	}

	err = bucket.Del(image.Uuid)
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusInternalServerError,
			"Could not delete image from image store",
		}
	}

	_, err = db.Delete(image)
	if err != nil {
		return nil, &ErrorResponse{
			http.StatusInternalServerError,
			fmt.Sprintf("Delete failed: %v", err),
		}
	}

	return nil, nil
}

var DeleteImage = WrapJsonHandler(nil, deleteImage)

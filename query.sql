-- name: CreateImage :exec
INSERT INTO Cloudflare_Image (path,image_link) VALUES ($1,$2);
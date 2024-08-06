-- name: CreateImage :exec
INSERT INTO Cloudflare_Image (path,author_name) VALUES ($1,$2);
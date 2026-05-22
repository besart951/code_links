module github.com/besart951/code-links/apps/planer-link/backend

go 1.26

require github.com/besart951/code-links/packages/productauth v0.0.0

require github.com/golang-jwt/jwt/v5 v5.3.0 // indirect

replace github.com/besart951/code-links/packages/productauth => ../../../packages/productauth

Creates a mock auth server which can be used to replace the authentication and jwt verification process from auth0. 
This should be used for things like integration and automated testing so that we do not need to use the auth0 tenancy 
in our tests.

Includes a docker image for use in gitlab pipelines. Run locally using `go run cmd/main.go`.

Docker commands to run:

 - `docker build . -t mock-auth:v1`
 - `docker run -d -p 7000:7000 --name mock-auth -t mock-auth:v1`

Get a token from the service with: `curl --url http://localhost:7000/generate-token -X POST -H "Content-Type: application/json" -d '{ "customClaims": { "http://hearth/userServiceId": <the user id that should be in the token> } }' -o -`
Get a token in kube : `curl --url http://mock-auth.default:7000/generate-token -X POST -H "Content-Type: application/json" -d '{ "customClaims": { "http://hearth/userServiceId": <the user id that should be in the token> } }' -o -`
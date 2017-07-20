# corkboard

## Contributors
* [Jason Moore](https://github.com/jasonmoore30)
* [Ben Wornom](https://github.com/bwornom7)

## Configuration
#### Create .env
In order for Corkboard to successfully connect to the couchbase server, a .env file needs to be created and contain the following fields:
* CB_CONNECTION: Where your couchbase server is located e.g. couchbase://localhost
* CB_BUCKET: The name of the bucket you want to open
* CB_BUCKET_PASS: The password to access that bucket (if you have one)
* CB_PRIVATE_RSA: the name of the file where your private RSA key is e.g. id_rsa
* CB_PORT: the port number the server listens on

## Docs
For in depth API documentation and details on current data constraints, visit our [wiki](https://github.com/acstech/corkboard/wiki) page.

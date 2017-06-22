# corkboard

## Contributors
* Jason Moore

## Configuration
#### Create .env
In order for Corkboard to successfully connect to the couchbase server, a .env file needs to be created and contain the following fields:
* CB_CONNECTION: Where your couchbase server is located e.g. couchbase://localhost
* CB_BUCKET: The name of the bucket you want to open
* CB_BUCKET_PASS: The password to access that bucket (if you have one)
* CB_PRIVATE_RSA: the name ofthe file where your private RSA key is e.g. id_rsa

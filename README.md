# sshca
Proof of concept for web API-based SSH certificate authority.

This is my first project using Go.

## Usage
To get the public CA host key: `GET /host`  
And for the public CA user key: `GET /user`


To sign a public host key: `POST /sign`
```json
{
    "type": 2,
    "principals": [
        "example.com",
        "www.example.com"
    ],
    "pubkey": "<SSH Host Public Key Here>"
}
```

To sign a public user key: `POST /sign`
```json
{
  "type": 1,
  "principals": [
    "root",
    "ec2-user"
  ],
  "pubkey": "<SSH User Public Key Here>",
  "source-address": "127.0.0.1,192.168.0.2"
}
```
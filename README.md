# stan-project

- The `stan-project` provides a RESTful API for handling risks. This service supports creating, retrieving, and managing risks.

## Table of Contents
1. [Setup](#setup)
2. [Running the Services](#running-the-services)
3. [Stopping and Restarting the Services](#stopping-and-restarting-the-services)
4. [Cleaning Up](#cleaning-up)
5. [Makefile Commands](#makefile-commands)
6. [API Endpoints](#api-endpoints)


## Setup

### Prerequisities

- Docker
- Make

## Running the services

### 1. cloning the Repository
Clone the repository to your local machine:

```bash
    git clone https://github.com/sobystanly/stan-project
    cd stan-project
```

### 2. Running Docker Compose
- Start PostgresSQL using Docker Compose

```bash
    docker compose up -d
```

### 3. Building and Running the service

- Navigate to the `stan-project` directory

```bash
    cd stan-project
```

Build and run `risk service` using the Makefile:

```bash
    make start
```

### Accessing the service

- Once the service is running you can access the APIs:
-  risk service: http://localhost:8080


### stopping and restarting the service

- To stop or restart the service you can use the Makefile

```bash
    make stop
    make restart
```

### Cleaning up

- To stop postgresSQL and remove container
- Navigate to project folder
- 
```bash
    cd stan-project
    docker compose down
```

### Additional Information
Makefile Commands
- `make build`: Build the binaries
- `make run`: Build and run the services
- `make clean`: Clean up binaries.
- `make test`: Run tests.

## API ENDPOINTS

- The risk service exposes the following endpoints

**Create a new Risk**

```http request
   POST localhost:8080/v1/risks
```

# Payload

```json
  {
    "title": "risk 1",
    "description": "cyber risk",
    "state": "open"
  }
```
# Response
- Status Codes

- 201 Created
```json
    {
      "id": "3adf28e9-c4f8-418a-b08b-2c070cd9653b",
      "state": "open",
      "title": "risk 1",
      "description": "cyber risk"
    }
```
- 400 Bad Request, if the request payload is invalid
- 500 Internal Server Error on all other errors

** GET a Risk By ID**

```http request
   GET localhost:8080/v1/risks/<id>
```

- This API enables you to fetch a risk by ID.

# Response

```json
    {
      "id": "c6778c3b-9e4e-45c9-99b5-a122566d4648",
      "state": "open",
      "title": "Man in the middle attack",
      "description": "Eavesdropping and listening to data exchange"
    }
```

# Status Codes
- 200 OK for successful GET
- 400 Bad Request if the riskID in the path param is invalid
- 500 Internal server error for internal server errors.

** GET ALL Risks**
- This API enables you to fetch all risks and this API supports pagination.

```http request
    GET localhost:8080/v1/risks
```

- Query Parameters

1. offset: The starting point for the list of risks(default 0)
2. limit: The maximum number of risks to return(default 10)
3. sortBy: The field to sort by(default title)
4. sortOrder: The sort order(`asc` or `desc`, default: asc)

```http request
    GET localhost:8080/v1/risks?offset=0&limit=10&sortBy=state&sortOrder=desc
```

Response

- The response contains a total count of risks and a list of paginated risks

```json
    {
  "totalCount": 9,
  "risks": [
    {
      "id": "68337c00-a37c-420a-8594-bbac72a77166",
      "state": "open",
      "title": "Social Engineering",
      "description": "Manipulating individuals into divulging confidential information"
    },
    {
      "id": "3adf28e9-c4f8-418a-b08b-2c070cd9653b",
      "state": "open",
      "title": "risk 1",
      "description": "cyber risk"
    },
    {
      "id": "c7e079c0-2f25-4059-a204-1200872356d8",
      "state": "closed",
      "title": "Quid pro Quo social engineering",
      "description": "Offering service or benefit in exchange of information"
    },
    {
      "id": "ac347e8f-7b5b-4aec-bf71-914e52db4e3b",
      "state": "open",
      "title": "Phishing",
      "description": "Fradulent attempts to obtain sensitive information"
    },
    {
      "id": "c6778c3b-9e4e-45c9-99b5-a122566d4648",
      "state": "open",
      "title": "Man in the middle attack",
      "description": "Eavesdropping and listening to data exchange"
    },
    {
      "id": "c99a0bdf-ede3-4061-881e-761a3433a814",
      "state": "open",
      "title": "Malware",
      "description": "malicious software designed to damage, disrupt or gain unauthorized access to computer systems"
    },
    {
      "id": "68402263-df56-4726-beed-4496d5dca1b0",
      "state": "open",
      "title": "Denial of Service",
      "description": "overwhelming the network or system resources with a flood of illegitimate requests"
    },
    {
      "id": "4e6202f3-07bc-4837-82d4-3b362cb51033",
      "state": "investigating",
      "title": "Credential Stuffing",
      "description": "Using compromised Credentials"
    },
    {
      "id": "219b186a-b307-41b1-b01a-48341bf7cee6",
      "state": "accepted",
      "title": "Baiting social engineering ",
      "description": "Creating fabricated scenario to obtain information"
    }
  ]
}
```

Status Code
- 200 OK for successful GET.
- 500 internal server error on all internal server errors.


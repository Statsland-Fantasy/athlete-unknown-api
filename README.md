# athlete-unknown-api

Backend API to facilitate creation and playing of Athlete Unknown - a sports trivia game application

## Getting Started

### Prerequisites

- Go 1.25 or higher
- AWS DynamoDB (or DynamoDB Local for development)
- AWS credentials configured (via AWS CLI, environment variables, or IAM role)

### Configuration

The API requires the following environment variables for DynamoDB configuration:

- `DYNAMODB_ENDPOINT` (optional): Custom DynamoDB endpoint URL. Use this for DynamoDB Local or custom endpoints. Leave empty for standard AWS DynamoDB.
- `ROUNDS_TABLE_NAME` (optional): Name of the rounds DynamoDB table. Defaults to `AthleteUnknownRoundsDev`.
- `USER_STATS_TABLE_NAME` (optional): Name of the user stats DynamoDB table. Defaults to `AthleteUnknownUserStatsDev`.
- `AWS_REGION` (optional): AWS region for DynamoDB. Defaults to `us-west-2`.
- `PORT` (optional): Port for the HTTP server. Defaults to `8080`.

### DynamoDB Table Structure

The application uses two separate DynamoDB tables:

#### 1. Rounds Table (AthleteUnknownRoundsDev)

**Primary Key:**

- `playDate` (String): Partition key in format `YYYY-MM-DD` (e.g., `2025-11-24`)
- `sport` (String): Sort key (e.g., `basketball`, `baseball`, `football`)

**Attributes:**
The table stores Round objects with all their nested attributes (Player, Stats, etc.)

**Example DynamoDB Local table creation:**

```bash
aws dynamodb create-table \
    --table-name AthleteUnknownRoundsDev \
    --attribute-definitions \
        AttributeName=playDate,AttributeType=S \
        AttributeName=sport,AttributeType=S \
    --key-schema \
        AttributeName=playDate,KeyType=HASH \
        AttributeName=sport,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://localhost:8000
```

#### 2. User Stats Table (AthleteUnknownUserStatsDev)

**Primary Key:**

- `userId` (String): Partition key (user's unique identifier)

**Attributes:**
The table stores UserStats objects with all their nested attributes (Sports, aggregate statistics, etc.)

**Example DynamoDB Local table creation:**

```bash
aws dynamodb create-table \
    --table-name AthleteUnknownUserStatsDev \
    --attribute-definitions \
        AttributeName=userId,AttributeType=S \
    --key-schema \
        AttributeName=userId,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://localhost:8000
```

**Performance Note:**
The `GetRoundsBySport` endpoint uses a Scan operation to filter by sport. For better performance with large datasets, consider adding a Global Secondary Index (GSI) with `sport` as the partition key and `playDate` as the sort key.

### Running the API

1. **Using AWS DynamoDB:**

```bash
export AWS_REGION=us-west-2
export ROUNDS_TABLE_NAME=AthleteUnknownRoundsDev
export USER_STATS_TABLE_NAME=AthleteUnknownUserStatsDev
go run .
```

2. **Using DynamoDB Local:**

```bash
export DYNAMODB_ENDPOINT=http://localhost:8000
export ROUNDS_TABLE_NAME=AthleteUnknownRoundsDev
export USER_STATS_TABLE_NAME=AthleteUnknownUserStatsDev
export AWS_REGION=us-west-2
go run .
```

3. **Change the server port:**

```bash
PORT=3000 go run .
```

## API Documentation

### Base URL

All API endpoints are prefixed with `/v1`

### Supported Sports

- `basketball`
- `baseball`
- `football`

---

## Endpoints

### Health Check

```
GET /health
```

Returns the health status of the API.

**Response:**

```json
{
  "status": "healthy"
}
```

---

### Game Rounds

#### Get a Round

```
GET /v1/round?sport={sport}&playDate={date}
```

Retrieves a round containing player information for the trivia game.

**Query Parameters:**

- `sport` (required): The sport to retrieve (`basketball`, `baseball`, or `football`)
- `playDate` (optional): The play date in `YYYY-MM-DD` format. Defaults to current date.

**Example:**

```bash
curl "http://localhost:8080/v1/round?sport=basketball&playDate=2025-11-15"
```

**Response:** `200 OK`

```json
{
  "roundId": "Basketball100",
  "sport": "basketball",
  "playDate": "2025-11-15",
  "created": "2025-11-11T10:00:00Z",
  "lastUpdated": "2025-11-11T14:30:00Z",
  "theme": "GOAT",
  "player": {
    "sport": "basketball",
    "sportsReferenceURL": "https://www.basketball-reference.com/players/j/jamesle01.html",
    "name": "LeBron James",
    "bio": "DOB: December 30, 1984 in Akron, Ohio",
    "playerInformation": "6'9\", 250 lbs, Forward, Shoots Right",
    "draftInformation": "Round 1 (1st overall) from St. Vincent-St. Mary High School",
    "yearsActive": "2003-Present",
    "teamsPlayedOn": "CLE, MIA, LAL",
    "jerseyNumbers": "#23, #6",
    "careerStats": "PPG: 27.2, RPG: 7.5, APG: 7.3, WS: 273.5",
    "personalAchievements": "4x NBA Champion, 4x NBA MVP, 19x NBA All-Star, 2x Olympic Gold Medalist",
    "photo": "https://cdn.triviagame.com/players/lebron-james.jpg"
  },
  "stats": {
    "playDate": "2025-11-15",
    "name": "LeBron James",
    "sport": "basketball",
    "totalPlays": 1247,
    "percentageCorrect": 68.5,
    "highestScore": 9,
    "averageCorrectScore": 7.8,
    "mostCommonFirstTileFlipped": "tile1",
    "mostCommonLastTileFlipped": "tile9",
    "mostCommonTileFlipped": "tile5",
    "leastCommonTileFlipped": "tile3"
  }
}
```

---

#### Create a Round

```
POST /v1/round
```

Creates a new game round with player information. Admin access required.

**Request Body:**

```json
{
  "roundId": "Basketball100",
  "sport": "basketball",
  "playDate": "2025-11-15",
  "theme": "GOAT",
  "player": {
    "sport": "basketball",
    "sportsReferenceURL": "https://www.basketball-reference.com/players/j/jamesle01.html",
    "name": "LeBron James",
    "bio": "DOB: December 30, 1984 in Akron, Ohio",
    "playerInformation": "6'9\", 250 lbs, Forward, Shoots Right",
    "draftInformation": "Round 1 (1st overall) from St. Vincent-St. Mary High School",
    "yearsActive": "2003-Present",
    "teamsPlayedOn": "CLE, MIA, LAL",
    "jerseyNumbers": "#23, #6",
    "careerStats": "PPG: 27.2, RPG: 7.5, APG: 7.3, WS: 273.5",
    "personalAchievements": "4x NBA Champion, 4x NBA MVP, 19x NBA All-Star",
    "photo": "https://cdn.triviagame.com/players/lebron-james.jpg"
  },
  "stats": {
    "playDate": "2025-11-15",
    "name": "LeBron James",
    "sport": "basketball",
    "totalPlays": 0,
    "percentageCorrect": 0.0,
    "highestScore": 0,
    "averageCorrectScore": 0.0
  }
}
```

**Example:**

```bash
curl -X POST http://localhost:8080/v1/round \
  -H "Content-Type: application/json" \
  -d @round.json
```

**Response:** `201 Created`

---

#### Delete a Round

```
DELETE /v1/round?sport={sport}&playDate={date}
```

Deletes an existing round. Admin access required.

**Query Parameters:**

- `sport` (required): The sport of the round to delete
- `playDate` (required): The play date in `YYYY-MM-DD` format

**Example:**

```bash
curl -X DELETE "http://localhost:8080/v1/round?sport=basketball&playDate=2025-11-15"
```

**Response:** `204 No Content`

---

#### Get Upcoming Rounds

```
GET /v1/upcoming-rounds?sport={sport}&startDate={date}&endDate={date}
```

Retrieves upcoming rounds for a specific sport. Admin access required.

**Query Parameters:**

- `sport` (required): The sport to retrieve rounds for
- `startDate` (optional): Start date for filtering in `YYYY-MM-DD` format
- `endDate` (optional): End date for filtering in `YYYY-MM-DD` format

**Example:**

```bash
curl "http://localhost:8080/v1/upcoming-rounds?sport=basketball&startDate=2025-11-15&endDate=2025-11-30"
```

**Response:** `200 OK` - Returns an array of Round objects

---

### Game Results

#### Submit Results

```
POST /v1/results?sport={sport}&playDate={date}
```

Submits the results of a completed trivia round.

**Query Parameters:**

- `sport` (required): The sport for the results
- `playDate` (required): The date of the round in `YYYY-MM-DD` format

**Request Body:**

```json
{
  "score": 9,
  "isCorrect": true,
  "tilesFlipped": [
    "tile1",
    "tile2",
    "tile3",
    "tile4",
    "tile5",
    "tile6",
    "tile7",
    "tile8",
    "tile9"
  ]
}
```

**Example:**

```bash
curl -X POST "http://localhost:8080/v1/results?sport=basketball&playDate=2025-11-15" \
  -H "Content-Type: application/json" \
  -d '{"score": 9, "isCorrect": true, "tilesFlipped": ["tile1", "tile2", "tile3"]}'
```

**Response:** `200 OK`

---

### Statistics

#### Get Round Statistics

```
GET /v1/stats/round?sport={sport}&playDate={date}
```

Retrieves statistics for a specific round.

**Query Parameters:**

- `sport` (required): The sport
- `playDate` (required): The play date in `YYYY-MM-DD` format

**Example:**

```bash
curl "http://localhost:8080/v1/stats/round?sport=basketball&playDate=2025-11-15"
```

**Response:** `200 OK`

```json
{
  "playDate": "2025-11-15",
  "name": "LeBron James",
  "sport": "basketball",
  "totalPlays": 1247,
  "percentageCorrect": 68.5,
  "highestScore": 9,
  "averageCorrectScore": 7.8,
  "mostCommonFirstTileFlipped": "tile1",
  "mostCommonLastTileFlipped": "tile9",
  "mostCommonTileFlipped": "tile5",
  "leastCommonTileFlipped": "tile3"
}
```

---

#### Get User Statistics

```
GET /v1/stats/user?userId={userId}
```

Retrieves comprehensive statistics for a specific user.

**Query Parameters:**

- `userId` (required): The user ID

**Example:**

```bash
curl "http://localhost:8080/v1/stats/user?userId=123e4567-e89b-12d3-a456-426614174000"
```

**Response:** `200 OK`

---

#### Update Username

```
PUT /v1/user/username
```

Updates the display username for the authenticated user. Requires JWT authentication.

**Authentication:** Required (JWT Bearer token)

**Permission Required:** `update:athlete-unknown:profile`

**Request Body:**

```json
{
  "userName": "MyNewUsername"
}
```

**Username Requirements:**
- Length: 3-20 characters
- Allowed characters: Letters, numbers, and spaces
- No consecutive spaces
- No inappropriate content

**Example:**

```bash
curl -X PUT "http://localhost:8080/v1/user/username" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"userName": "CoolPlayer123"}'
```

**Response:** `200 OK`

```json
{
  "userId": "auth0|123456",
  "userName": "CoolPlayer123",
  "message": "Username updated successfully"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid username (too short/long, invalid characters, inappropriate content)
- `401 Unauthorized` - Missing or invalid JWT token
- `500 Internal Server Error` - Database error

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "Bad Request",
  "message": "Sport parameter is required",
  "code": "MISSING_REQUIRED_PARAMETER",
  "timestamp": "2025-11-11T10:30:00Z"
}
```

### Common Error Codes

- `MISSING_REQUIRED_PARAMETER` - A required parameter is missing
- `INVALID_PARAMETER` - A parameter has an invalid value
- `ROUND_NOT_FOUND` - The requested round does not exist
- `ROUND_ALREADY_EXISTS` - A round already exists for the sport/date
- `STATS_NOT_FOUND` - Statistics not found
- `USER_STATS_NOT_FOUND` - User statistics not found
- `METHOD_NOT_ALLOWED` - HTTP method not supported

---

## Features

- RESTful API design following OpenAPI 3.0 specification
- DynamoDB integration for persistent data storage
- Configurable DynamoDB endpoint (supports DynamoDB Local)
- CORS support for cross-origin requests
- Request logging middleware
- Comprehensive error handling with specific error codes
- Support for three sports: basketball, baseball, and football

## Project Structure

```
athlete-unknown-api/
├── main.go                      # HTTP server setup, routing, and middleware
├── handlers.go                  # API endpoint handlers
├── models.go                    # Data models and structures
├── database.go                  # DynamoDB operations
├── config.go                    # Configuration management
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── AthleteUnknownAPISpec.yaml  # OpenAPI specification
└── README.md                    # This file
```

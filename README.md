# athlete-unknown-api
Backend API to facilitate creation and playing of Athlete Unknown - a sports trivia game application

## Getting Started

### Prerequisites
- Go 1.25 or higher

### Running the API

1. Run the server:
```bash
go run .
```

The server will start on port 8080 by default. You can change the port by setting the `PORT` environment variable:
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
  "previouslyPlayedDates": ["2025-11-01", "2025-11-08"],
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
  "previouslyPlayedDates": ["2025-11-01", "2025-11-08"],
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
  "tilesFlipped": ["tile1", "tile2", "tile3", "tile4", "tile5", "tile6", "tile7", "tile8", "tile9"]
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

## Project Structure

```
athlete-unknown-api/
├── main.go                      # HTTP server setup, routing, and middleware
├── handlers.go                  # API endpoint handlers
├── models.go                    # Data models and structures
├── go.mod                       # Go module definition
├── AthleteUnknownAPISpec.yaml  # OpenAPI specification
└── README.md                    # This file
```

## Features

- RESTful API design following OpenAPI 3.0 specification
- In-memory data storage (easily replaceable with database)
- CORS support for cross-origin requests
- Request logging middleware
- Comprehensive error handling with specific error codes
- Thread-safe operations with mutex locks
- Support for three sports: basketball, baseball, and football
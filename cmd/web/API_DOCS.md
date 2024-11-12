## API Documentation

### Base URL

`/api/v1/`

### Authentication

- **JWT Tokens**: Used for authentication. Tokens are passed in the `Authorization` header as `Bearer <token>`.

### Endpoints

#### User Registration

- **URL**: `/register`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }

- **Response**:
  ```json
  {
    "message": "User registered successfully",
    "token": "jwt_token"
  }

#### User login

- **URL**: `/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
- **Response**:

  ```json
  {
    "message": "Login successful",
    "token": "jwt_token"
  }

#### Create a Yap


- **URL**: `/yap`
- **Method**: `POST`
- **Request Body**:
- **Headers**: ***Authorization***: Bearer <token> 

- **Response**:
  ```json
  {
    "message": "Yap created successfully",
    "yap_id": "uuid"
  }

#### Get a User's Feed

-**URL**: `/feed/?userId={user_id}&page={page_number}`
-**Method**: `GET`
-**Headers**: ***Authorization***: Bearer <token>

-**Response**:
```json
[
  {
    "yap_id": "uuid",
    "content": "string",
    "username": "string",
    "timestamp": "datetime"
  },
  {
    "yap_id": "uuid",
    "content": "string",
    "username": "string",
    "timestamp": "datetime"
  }
]


Follow a User

URL: /follow/{username}

Method: POST

Headers: Authorization: Bearer <token>

Response:

json
Copy
{
  "message": "You are now following {username}"
}


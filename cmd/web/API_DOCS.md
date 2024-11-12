## API Documentation

### Base URL

`/api/v1/`

### Authentication

- **JWT Tokens**: Used for authentication. Tokens are passed in the `Authorization` header as `Bearer <token>`.

### Endpoints

#### User Registration

- **URL**: `/signup`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "username":    "string",
    "email":       "string",
    "password":    "string",   
    "expirtes_in": "number"   optional  
  }
- **Response**:
  ```json User
  {
		"ID":        "uuid",
		"CreatedAt": "datetime",
		"UpdatedAt": "datetime",
		"Email":     "string",
		"Username":  "string",
		"Premuim":   "bool",
		"Token":     "JWT token",
  }
#### User login

- **URL**: `/login`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "email":       "string",
    "password":    "string"
    "expirtes_in": "number"   optional  
  }
- **Response**:

  ```json
  {
    "ID":        "uuid",
		"CreatedAt": "datetime",
		"UpdatedAt": "datetime",
		"Email":     "string",
		"Token":     "JWT token",
  }
#### Yap Creation


- **URL**: `/yap`
- **Method**: `POST`
- **Request Body**:
- **Headers**: ***Authorization***: Bearer <token> 

- **Request Body**:
  ```json
  {
    "body":       "string",
  }
- **Response**:
  ```json
  {
    "ID":        "uuid",
		"UpdatedAt": "datetime",
		"CreatedAt": "datetime",
		"Body":      "string",
		"UserID":    "uuid",
  }
#### Getting Yaps

- **URL**: `/yaps/user/{user_id}`
- **Method**: `GET`

- **Response**:
  ```json
  [
    {
    "ID":        "uuid",
		"UpdatedAt": "datetime",
		"CreatedAt": "datetime",
		"Body":      "string",
		"UserID":    "uuid",
    }
  ]

- **URL**: `/yaps/{yap_id}`
- **Method**: `GET`

- **Response**:
  ```json
    {
    "ID":        "uuid",
		"UpdatedAt": "datetime",
		"CreatedAt": "datetime",
		"Body":      "string",
		"UserID":    "uuid",
    }
#### Get a User's Feed

- **URL**: `/feed/?userId={user_id}&page={page_number}`
- **Method**: `GET`
- **Headers**: ***Authorization***: Bearer <token>

- **Response**:
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
#### Get a User's Followers
- **URL**: `/followers/{user_id}`
- **Method**: `GET`

- **Response**:
  ```json
  [
    {
    "ID":        "uuid",
		"CreatedAt": "datetime",
		"UpdatedAt": "datetime",
		"Email":     "string",
		"Token":     "JWT token",
    }
  ]
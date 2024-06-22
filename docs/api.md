# ENDPOINTS

## AUTH

### 1. Sign-up
- **Endpoint:** /auth/sign-up
- **Method:** POST
- **Path params:**
- **Query params:**
- **Request body:**
```
{
    "login": "adsa",
    "password": "sdas"
}
```
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": {
        "id": 1
    }
}
```

### 2. Sign-in
- **Endpoint:** /auth/sign-in
- **Method:** POST
- **Path params:**
- **Query params:**
- **Request body:**
```
{
    "login": "adsa",
    "password": "sdas"
}
```
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": {
        "accessToken":  "accessToken",
        "refreshToken": "refreshToken",
    }
}
```

### 3. Refresh Tokens
- **Endpoint:** /auth/refresh
- **Method:** POST
- **Path params:**
- **Query params:**
- **Request body:**
```
{
    "refreshToken": "adsa"
}
```
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": {
        "accessToken":  "accessToken",
        "refreshToken": "refreshToken",
    }
}
```

### 4. Logout (Authorization header required)
- **Endpoint:** /auth/logout
- **Method:** POST
- **Path params:**
- **Query params:**
- **Request body:**
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": null
}
```

## API (Authorization header required)

### ORGANIZATIONS
#### 5. Get All Users Organizations
- **Endpoint:** /api/organizations
- **Method:** GET
- **Path params:**
- **Query params:**
- **Request body:**
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": [
        {
            "id": 1,
            "name": "org1",
            "role": "admin"
        }
    ]
}
```

#### 6. Create Organization
- **Endpoint:** /api/organizations
- **Method:** POST
- **Path params:**
- **Query params:**
- **Request body:**
```
{
    "name": "org1"
}
```
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": {
        "id":  1
    }
}
```

#### 7. Delete Organization
- **Endpoint:** /api/organizations/:id
- **Method:** DELETE
- **Path params:**
```
id - organization id
```
- **Query params:**
- **Request body:**
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": {
        "id":  1
    }
}
```

#### 8. Change Organization Name
- **Endpoint:** /api/organizations/:id
- **Method:** PATCH
- **Path params:**
```
id - organization id
```
- **Query params:**
- **Request body:**
```
{
    "name": "cool name"
}
```
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": {
        "id": 1
    }
}
```

### USERS ORGANIZATIONS INTERACTION

#### 9. Invite User To Organization
- **Endpoint:** /api/organizations/invite/:id
- **Method:** POST
- **Path params:**
```
id - user id to invite
```
- **Query params:**
- **Request body:**
```
{
    "organizationId": 1,
    "role": "reader"
}
```
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": null
}
```

#### 10. Delete User From Organization
- **Endpoint:** /api/organizations/kick/:id
- **Method:** POST
- **Path params:**
```
id - user id to kick
```
- **Query params:**
- **Request body:**
```
{
    "organizationId": 1
}
```
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": {
        "id": 1
    }
}
```

#### 11. Change User Role In Organization
- **Endpoint:** /api/organizations/change-role/:id
- **Method:** POST
- **Path params:**
```
id - user id to change role
```
- **Query params:**
- **Request body:**
```
{
    "organizationId": 1,
    "role": "redactor"
}
```
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": {
        "id": 1
    }
}
```

### USERS

#### 12. Get All Users
- **Endpoint:** /api/users
- **Method:** GET
- **Path params:**
- **Query params:**
```
limit  - limits user to get (number of users in one page)
page   - page number
search - string to find user by it
```
- **Request body:**
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": {
        "count": 10,
        "rows": [
            {
                "id": 1,
                "login": "some_user"
            }
        ]
    }
}
```

#### 12. Get All Users In Organization
- **Endpoint:** /api/users/:id
- **Method:** GET
- **Path params:**
```
id - organization id
```
- **Query params:**
```
limit  - limits user to get (number of users in one page)
page   - page number
role   - to find users with certain role
```
- **Request body:**
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": {
        "count": 10,
        "rows": [
            {
                "id": 1,
                "login": "some_user",
                "role": "admin"
            }
        ]
    }
}
```

### DATASOURCES

#### 13. Create Datasource
- **Endpoint:** /api/datasources/:id
- **Method:** POST
- **Path params:**
```
id - organization id
```
- **Query params:**
- **Request body:**
```
{
    "name": "my_datasource"
}
```
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": {
        "id": 1
    }
}
```

#### 14. Get Datasources In Organization
- **Endpoint:** /api/datasources/:id
- **Method:** GET
- **Path params:**
```
id - organization id
```
- **Query params:**
- **Request body:**
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": [
        {
            "id": 1,
            "name" "datasource1"
        }
    ]
}
```

#### 15. Delete Datasource
- **Endpoint:** /api/datasources/:id
- **Method:** DELETE
- **Path params:**
```
id - datasource id
```
- **Query params:**
- **Request body:**
- **Response body:**
```
{
    "status": 200,
    "message": "ok",
    "payload": {
        "id": 1
    }
}
```

## PROXY ALLOWED ENDPOINTS (Authorization header required)

#### 16. Get Datasource Connection Info
- **Endpoint:** /connect/:id
- **Method:** GET
- **Path params:**
```
id - datasource id
```
- **Query params:**
- **Request body:**
- **Response body:**
```
{
    "host":     "postgres-asdasd-a-d2112312",
    "port":     7891,
    "user":     "admin",
    "password": "djiqwejqcjeionqe",
    "name":     "datasource1",
}
```
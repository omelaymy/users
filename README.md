# Users service
---

## Run service

Clone repository by:
```
https://github.com/omelaymy/users
cd users
```
Build and Run service:

```
docker build . -t users:latest
docker run -p 8888:8888 users
```

### REST server for storage and lightweight management of user profiles

This web service provides an API for managing user profiles and authentication.
Each profile contains information about the user, such as a unique identifier, email, username, password, and administrator status.

### API Documentation:

`http://localhost:8888/docs`

### Authentication:

The service uses Basic Access Authentication.
Access to protected endpoints requires providing correct user credentials in the authorization header.

### Access Restrictions:

All registered users can view user profiles.
Creation, modification, and deletion of profiles can only be performed by users with the administrator role (admin).

### Data Storage:
A primitive in-memory database is implemented to store user profiles in RAM. Data will be reset upon service restart.

This web service provides a simple and lightweight way to manage user profiles and ensures basic authentication to protect user's confidential data.
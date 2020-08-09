# User Service

User Service is a microservice written in Golang which handles management of user information, and the registration/login workflow.

## Package Overview
- auth: utilities relating to authentication (JWT validation, registration, login)
- config: extracts configuration information from the environment, and initializes database connections
- grpc: implements the GRPC API for User Service
- model: abstractions over database operations for users
 
## Overview of Workflow
### Login
1. The user sends their username and password to User Service.
2. The user's plaintext password is hashed via bcrypt, and the hash is compared to the stored password hash.
3. If the password matches, we return a signed JWT with their user ID as the sole claim.

### Registration
1. User sends username, password and email
2. The password is hashed with bcrypt, and stored in the users table along with their username and email.

### JWT Validation
1. Frontend sends the JWT to User Service via the ValidateJWT API
2. User Service returns whether the JWT is valid, and the user's ID if so.


## TODOs
1. Add expiration date to JWT claims, and add expiration date validation to the token validation process.
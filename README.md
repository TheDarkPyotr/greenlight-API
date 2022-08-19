## greenlight-API
WIP


### Project structure

- **bin/**: contain application   binaries  
- cmd/  
- **api/** :code for running server, reading/writing HTTP requests, managing authentication  
- **internal/**: code for persistence, data validation, email sending, etc. Imported from _cmd/api/_  
- **migrations/**: migrations for our database  
- **remote/**: script for deploying on server  

### API Endpoints
 | Method | URL Pattern        | Handler              | Action                                     |
|--------|--------|--------------------|----------------------------------------------------------|
 | GET    | `/v1/healthcheck`  | `healthcheckHandler` | Show application information                             |
 | POST   | `/v1/movies`       | `createMovieHandler` | Create a new movie                                       |
 | GET    | `/v1/movies` | `listMovieHandler` | Show details of all movies |
 | GET    | `/v1/movies/:id`   | `showMovieHandler`   | Show details of a specific movie                         |
 | PATCH  | `/v1/movies/:id` | `updateMovieHandler`   | Update the details of a specific movie                   |
 | DELETE | `/v1/movies/:id` | `deleteMovieHandler`   | Delete a specific movie                                  |
 
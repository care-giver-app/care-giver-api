# Care Giver API

This API handles requests about user and receiver data. It uses AWS SAM to build and deploy an
Lambda API fronted by the API Gateway. It interacts with DynamoDB to create and update user and receiver 
objects. 

## Usage
The following endpoints are currently available
- `/user` - POST
- `/user/{userId}` - GET
- `/user/primary-receiver` - POST
- `/user/additional-receiver` - POST
- `/receiver/{receiverId}` - GET
- `/receiver/event` - POST


## Running Locally
Prerequisite: Make sure you have local dynamodb running with the following tables created:
- `user-table-local`
- `receiver-table-local`
- `event-table-local`

To start the api:
```sh
make start-api
```

To invoke via an event in the `events/` directory:
```sh
make invoke EVENT=someEvent.json
```

## Testing

### Unit Tests
Run unit tests
```sh
make test 
```

Run unit tests with html coverage report
```sh
make test-report
```

### Component Tests
TODO
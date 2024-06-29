# Full Stack Job Scheduler - Go, React, and WebSockets

## Instructions to run locally
For the server,

- `cd backend`
- `got run main.go`

For the client,

- `cd frontend`
- `npm install`
- `npm start`

Open `http://localhost:3000/`

## Design choices
- Used priority queue + sync.Mutex to implement SJF algorithm for executing the highest priority (shortest duration) job as a blocking task
- Mocked execution of task by sleeping for job duration
- Used `github.com/gorilla/websocket` to create a websocket server in Go backend
- Utilized Tailwind CSS for styling the form and list components in React UI
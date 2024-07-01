# Full Stack Job Scheduler - Go, React, and WebSockets

## Instructions to run locally
For the server,

- `cd backend`
- `got run main.go`

For the client,

- `cd frontend`
- `npm install`
- `npm start`

Open `http://localhost:3000/`, add a few jobs, and it should look as follows:

![Screen Shot 2024-07-01 at 2 40 04 PM](https://github.com/kalpitf1/job_scheduler/assets/37945736/8ed492e5-f110-4e17-98ee-c881d0967609)

## Design choices
- Used priority queue + sync.Mutex to implement SJF algorithm for executing the highest priority (shortest duration) job as a blocking task
- Mocked execution of task by sleeping for job duration
- Used `github.com/gorilla/websocket` to create a websocket server in Go backend
- Utilized Tailwind CSS for styling the form and list components in React UI

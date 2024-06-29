import React, { useState, useEffect } from 'react';
import axios from 'axios';

const App = () => {
  const [jobs, setJobs] = useState([]);
  const [name, setName] = useState('');
  const [duration, setDuration] = useState('');

  useEffect(() => {
    fetchJobs();
    const ws = new WebSocket('ws://localhost:8080/ws');

    ws.onmessage = (event) => {
      const updatedJob = JSON.parse(event.data);
      setJobs((prevJobs) => {
        const jobIndex = prevJobs.findIndex(job => job.id === updatedJob.id);
        if (jobIndex !== -1) {
          const newJobs = [...prevJobs];
          newJobs[jobIndex] = updatedJob;
          return newJobs;
        } else {
          return [...prevJobs, updatedJob];
        }
      });
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed');
    };

    return () => ws.close();
  }, []);

  const fetchJobs = async () => {
    try {
      const response = await axios.get('http://localhost:8080/jobs');
      if (response.data) {
        setJobs(response.data);
      }
    } catch (error) {
      console.error('Error fetching jobs:', error);
    }
  };

  const addJob = async () => {
    const newJob = { name, duration: parseInt(duration) * Math.pow(10, 9) };  // server accepts time in nanoseconds
    try {
      const response = await axios.post('http://localhost:8080/jobs', newJob);
      setJobs([...jobs, response.data]);
      setName('');
      setDuration('');
    } catch (error) {
      console.error('Error adding job:', error);
    }
  };

  return (
    <div className="App">
      <h1>Job Management</h1>
      <div>
        <h2>Add Job</h2>
        <input
          type="text"
          placeholder="Job Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <input
          type="text"
          placeholder="Duration"
          value={duration}
          onChange={(e) => setDuration(e.target.value)}
        />
        <button onClick={addJob}>Add Job</button>
      </div>
      <div>
        <h2>Jobs List</h2>
        {jobs.length === 0 ? (
          <p>No jobs available</p>
        ) : (
          <ul>
            {jobs.map((job, index) => (
              <li key={index}>
                {job.name} - {job.duration / Math.pow(10, 9)} - {job.status}  {/* server returns time in nanoseconds */}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default App;

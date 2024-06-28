import React, { useState, useEffect } from 'react';
import axios from 'axios';

const App = () => {
  const [jobs, setJobs] = useState([]);
  const [name, setName] = useState('');
  const [duration, setDuration] = useState('');
  const [status, setStatus] = useState('');

  useEffect(() => {
    fetchJobs();
  }, []);

  const fetchJobs = async () => {
    try {
      const response = await axios.get('http://localhost:8080/jobs');
      setJobs(response.data);
    } catch (error) {
      console.error('Error fetching jobs:', error);
    }
  };

  const addJob = async () => {
    const newJob = { name, duration, status };
    try {
      const response = await axios.post('http://localhost:8080/jobs', newJob);
      setJobs([...jobs, response.data]);
      setName('');
      setDuration('');
      setStatus('');
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
        <input
          type="text"
          placeholder="Status"
          value={status}
          onChange={(e) => setStatus(e.target.value)}
        />
        <button onClick={addJob}>Add Job</button>
      </div>
      <div>
        <h2>Jobs List</h2>
        <ul>
          {jobs.map((job, index) => (
            <li key={index}>
              {job.name} - {job.duration} - {job.status}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default App;

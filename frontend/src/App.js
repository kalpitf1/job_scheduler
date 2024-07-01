import React, { useState, useEffect } from "react";
import axios from "axios";
import JobForm from "./components/form/JobForm";
import JobList from "./components/list/JobList";

const App = () => {
  const [jobs, setJobs] = useState([]);

  useEffect(() => {
    fetchJobs();
    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onmessage = (event) => {
      const updatedJob = JSON.parse(event.data);
      setJobs((prevJobs) => {
        const jobIndex = prevJobs.findIndex((job) => job.id === updatedJob.id);
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
      console.log("WebSocket connection closed");
    };

    return () => ws.close();
  }, []);

  const fetchJobs = async () => {
    try {
      const response = await axios.get("http://localhost:8080/jobs");
      if (response.data) {
        setJobs(response.data);
      }
    } catch (error) {
      console.error("Error fetching jobs:", error);
    }
  };

  const addJob = (newJob) => {
    setJobs([...jobs, newJob]);
  };

  return (
    <div className="bg-white">
      <div className="mx-auto max-w-7xl px-4 py-24 sm:px-6 sm:py-32 lg:px-8">
        <div className="mx-auto max-w-2xl">
          <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
            Job Scheduler
          </h2>
          <JobForm addJob={addJob} />
          <JobList jobs={jobs} />
        </div>
      </div>
    </div>
  );
};

export default App;

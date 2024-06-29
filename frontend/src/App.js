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

  const resetJob = async () => {
    setName('');
    setDuration('');
  }

  const getStatusClass = (status) => {
    switch (status) {
      case 'Pending':
        return 'text-red-500';
      case 'In Progress':
        return 'text-yellow-500';
      case 'Completed':
        return 'text-green-500';
      default:
        return 'text-gray-500';
    }
  };

  return (
    <div class="bg-white">
      <div class="mx-auto max-w-7xl px-4 py-24 sm:px-6 sm:py-32 lg:px-8">
        <div class="mx-auto max-w-2xl">
          <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
            Job Scheduler
          </h2>
          <form style={{ "margin": "50px 0" }}>
            <div className="space-y-12">
              <div className="border-b border-gray-900/10 pb-12">
                <h2 className="text-base font-semibold leading-7 text-gray-900">Add Job</h2>
                <p className="mt-1 text-sm leading-6 text-gray-600">Jobs will be processed on a Shortest Job First (SJF) basis.</p>

                <div className="mt-10 grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
                  <div className="sm:col-span-3">
                    <label htmlFor="first-name" className="block text-sm font-medium leading-6 text-gray-900">
                      Job Name
                    </label>
                    <div className="mt-2">
                      <input
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      />
                    </div>
                  </div>

                  <div className="sm:col-span-3">
                    <label htmlFor="last-name" className="block text-sm font-medium leading-6 text-gray-900">
                      Duration (in seconds)
                    </label>
                    <div className="mt-2">
                      <input
                        type="text"
                        value={duration}
                        onChange={(e) => setDuration(e.target.value)}
                        className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      />
                    </div>
                  </div>

                </div>
              </div>
            </div>

            <div className="mt-6 flex items-center justify-end gap-x-6">
              <button
                type="button"
                onClick={resetJob}
                className="text-sm font-semibold leading-6 text-gray-900">
                Cancel
              </button>
              <button
                type="submit"
                onClick={addJob}
                className="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                Save
              </button>
            </div>
          </form>

          <div class="mx-auto max-w-4xl">
            <h2 className="text-base font-semibold leading-7 text-gray-900">Job List</h2>

            <div className="mt-10">

              {jobs.length === 0 ? (
                <p className="mt-1 text-sm leading-6 text-gray-600">No jobs available</p>
              ) : (
                <ul role="list" class="divide-y divide-gray-100">
                  {jobs.map((job, index) => (
                    <li class="flex justify-between gap-x-6 py-5" key={index}>
                      <div class="min-w-0 flex-auto">
                        <p class="text-sm leading-6 text-gray-900">{job.name}</p>
                        <p class="mt-1 truncate text-xs leading-5 text-gray-500">{job.duration / Math.pow(10, 9)} seconds</p> {/* server returns time in nanoseconds */}
                      </div>
                      <div class="hidden shrink-0 sm:flex sm:flex-col sm:items-end">
                        <p class={`mt-1 text-xs leading-5 text-gray-500 ${getStatusClass(job.status)}`}>{job.status}</p>
                      </div>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>


        </div>
      </div>
    </div>
  );
};

export default App;

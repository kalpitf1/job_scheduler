import React, { useState } from "react";
import axios from "axios";

const JobForm = ({ addJob }) => {
  const [name, setName] = useState("");
  const [duration, setDuration] = useState("");

  const handleAddJob = async (e) => {
    e.preventDefault();
    const newJob = { name, duration: parseInt(duration) * Math.pow(10, 9) }; // server accepts time in nanoseconds
    try {
      const response = await axios.post("http://localhost:8080/jobs", newJob);
      addJob(response.data);
      setName("");
      setDuration("");
    } catch (error) {
      console.error("Error adding job:", error);
    }
  };

  const resetJob = () => {
    setName("");
    setDuration("");
  };

  return (
    <form style={{ margin: "50px 0" }} onSubmit={handleAddJob}>
      <div className="space-y-12">
        <div className="border-b border-gray-900/10 pb-12">
          <h2 className="text-base font-semibold leading-7 text-gray-900">
            Add Job
          </h2>
          <p className="mt-1 text-sm leading-6 text-gray-600">
            Jobs will be processed on a Shortest Job First (SJF) basis.
          </p>

          <div className="mt-10 grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
            <div className="sm:col-span-3">
              <label
                htmlFor="name"
                className="block text-sm font-medium leading-6 text-gray-900"
              >
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
              <label
                htmlFor="duration"
                className="block text-sm font-medium leading-6 text-gray-900"
              >
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
          className="text-sm font-semibold leading-6 text-gray-900"
        >
          Cancel
        </button>
        <button
          type="submit"
          className="rounded-md bg-indigo-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
        >
          Save
        </button>
      </div>
    </form>
  );
};

export default JobForm;

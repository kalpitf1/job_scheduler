import React from "react";

const JobList = ({ jobs }) => {
  const getStatusClass = (status) => {
    switch (status) {
      case "Pending":
        return "text-red-500";
      case "In Progress":
        return "text-yellow-500";
      case "Completed":
        return "text-green-500";
      default:
        return "text-gray-500";
    }
  };

  return (
    <div className="mx-auto max-w-4xl">
      <h2 className="text-base font-semibold leading-7 text-gray-900">
        Job List
      </h2>
      <div className="mt-10">
        {jobs.length === 0 ? (
          <p className="mt-1 text-sm leading-6 text-gray-600">
            No jobs available
          </p>
        ) : (
          <ul className="divide-y divide-gray-100">
            {jobs.map((job, index) => (
              <li className="flex justify-between gap-x-6 py-5" key={index}>
                <div className="min-w-0 flex-auto">
                  <p className="text-sm leading-6 text-gray-900">{job.name}</p>
                  <p className="mt-1 truncate text-xs leading-5 text-gray-500">
                    {job.duration / Math.pow(10, 9)} seconds
                  </p>{" "}
                  {/* server returns time in nanoseconds */}
                </div>
                <div className="hidden shrink-0 sm:flex sm:flex-col sm:items-end">
                  <p
                    className={`mt-1 text-xs leading-5 ${getStatusClass(
                      job.status
                    )}`}
                  >
                    {job.status}
                  </p>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default JobList;

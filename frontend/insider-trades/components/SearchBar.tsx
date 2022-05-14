import React from "react";

function SearchBar() {
  return (
    <div className="flex-1 max-w-lg">
      <form className="py-4 px-0 sm:px-4 flex justify-between items-center">
        <label className="mb-2 text-sm font-medium text-gray-900 sr-only dark:text-gray-300">Search</label>
        <div className="flex-1 relative w-auto">
          <input
            className="block p-2 pl-4 w-full text-sm text-gray-900 bg-slate-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            placeholder="Search for tickers, insider etc.."
          />
        </div>
        <button className="mx-2 px-6 py-2 bg-indigo-500 hover:bg-indigo-600 text-gray-50 rounded-xl flex items-center">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth="2"
          >
            <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </button>
      </form>
      <div className="">
      </div>
    </div>
  );
}

export default SearchBar;

import React from "react";

function Search() {
  return (
    <div className="main-search flex justify-center w-auto relative text-gray-600 m-5 md:m-0 md:mx-5">
      <input
        type="search"
        name="search"
        placeholder="Search for trades..."
        className="bg-white h-10 px-5 pr-12 rounded-full text-sm focus:outline-none"
      />
      <button type="submit" className="absolute right-0 top-0 mt-2 mr-4">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="block m-auto h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth="1.8"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
      </button>
    </div>
  );
}

export default Search;

import React, { useState } from "react";

import { useLocation, Link } from 'react-router-dom';

const Home = () => {
  return (
    <>
    <main className="mt-10 mx-auto max-w-7xl px-4 sm:mt-12 sm:px-6 md:mt-16 lg:mt-20 lg:px-8 xl:mt-28">
      <div className="sm:text-center lg:text-left">
        <h1 className="text-4xl tracking-tight font-extrabold text-gray-900 sm:text-5xl md:text-6xl">
          <span className="block xl:inline">chainlink</span>
          <span className="block text-indigo-600 xl:inline"> separating dev from prod</span>
        </h1>
        <p className="mt-3 text-base text-gray-500 sm:mt-5 sm:text-lg sm:max-w-xl sm:mx-auto md:mt-5 md:text-xl lg:mx-0">Chainlink is a service that intercepts your DNS requests and rewrites them to your local running services so you can test production changes locally.</p>
        <div className="mt-5 sm:mt-8 sm:flex sm:justify-center lg:justify-start">
          <div className="rounded-md shadow">
            <a href="https://github.com/mikemackintosh/chainlink" className="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 md:py-4 md:text-lg md:px-10"> Learn more</a>
          </div>
        </div>
      </div>
    </main>
    </>
  );
}

export default Home;

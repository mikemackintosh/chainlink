import React, { useState } from "react";

import { useLocation, Link } from 'react-router-dom';

const Home = () => {
  return (
    <>
    <svg class="hidden lg:block absolute right-0 inset-y-0 h-full w-48 text-white transform translate-x-1/2" fill="currentColor" viewBox="0 0 100 100" preserveAspectRatio="none" aria-hidden="true">
      <polygon points="50,0 100,0 50,100 0,100" />
    </svg>

    <main class="mt-10 mx-auto max-w-7xl px-4 sm:mt-12 sm:px-6 md:mt-16 lg:mt-20 lg:px-8 xl:mt-28">
      <div class="sm:text-center lg:text-left">
        <h1 class="text-4xl tracking-tight font-extrabold text-gray-900 sm:text-5xl md:text-6xl">
          <span class="block xl:inline">chainlink</span>
          <span class="block text-indigo-600 xl:inline"> separating dev from prod</span>
        </h1>
        <p class="mt-3 text-base text-gray-500 sm:mt-5 sm:text-lg sm:max-w-xl sm:mx-auto md:mt-5 md:text-xl lg:mx-0">Chainlink is a service that intercepts your DNS requests and rewrites them to your local running services so you can test production changes locally.</p>
        <div class="mt-5 sm:mt-8 sm:flex sm:justify-center lg:justify-start">
          <div class="rounded-md shadow">
            <a href="/hosts" class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 md:py-4 md:text-lg md:px-10"> Manage Sites </a>
          </div>
          <div class="mt-3 sm:mt-0 sm:ml-3">
            <a href="https://github.com/mikemackintosh/chainlink" class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-indigo-700 bg-indigo-100 hover:bg-indigo-200 md:py-4 md:text-lg md:px-10"> Learn more</a>
          </div>
        </div>
      </div>
    </main>
    </>
  );
}


export default Home;
import React from 'react';
import {
  withRouter,
  Route,
  Switch,
  BrowserRouter as Router,
} from "react-router-dom";

import CookieConsent from "react-cookie-consent";

import Navbar from './components/navbar.js';
import Home from './pages/home.js';
import Hosts from './pages/hosts.js';
import NotFound from './pages/not_found.js';
import './App.css';

const routes = [
  { path: "/", component: Home },
  { path: "/hosts", component: Hosts },
];

const App = () => {
  return (
    <>
    <div className="App">
<div class="relative bg-white overflow-hidden">
  <div class="max-w-7xl mx-auto">
    <div class="relative z-10 pb-8 bg-white sm:pb-16 md:pb-20 lg:max-w-2xl lg:w-full lg:pb-28 xl:pb-32">

      <Navbar />

      <Switch>
        {routes.map((route, idx) => (
          <Route path={route.path} exact component={route.component} key={idx} />
        ))}
        <Route component={NotFound} />
      </Switch>

      </div>
      </div>
      </div>
      </div>

      <CookieConsent
        location="bottom"
        buttonText="Continue"
        cookieName="ags-cc-store1"
        style={{ background: "rgb(32 40 64)" }}
        buttonStyle={{ background: "#335eea", color: "#ffffff", fontSize: "14px", borderRadius: '8px' }}
        expires={7}>
        <b>We use cookies to enhance your experience.{" "}</b>
        <span style={{ fontSize: "12px" }}>For the record, we don't sell any of your data.</span>
      </CookieConsent>
    </>
  );
}

export default withRouter(App);

import React, { useState, useEffect} from "react";
import classNames from "classnames";
import {
  withRouter,
  Route,
  Switch,
} from "react-router-dom";

import Home from './pages/home.js';
import AddNew from './pages/addnew.js';
import Hosts from './pages/hosts.js';
import NotFound from './pages/not_found.js';
import Sidebar from "./components/sidebar.js";

import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';

const routes = [
  { path: "/", component: Home },
  { path: "/hosts/addnew", component: AddNew },
  { path: "/hosts/:zone", component: Hosts },
];

const App = () => {
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [updatedZones, updateZones] = useState(false);
  const [zones, setZones] = useState([])

  function getZones(){
    fetch("http://chainlink.config/api/config")
      .then((res) => res.json())
      .then((json) => {
        setZones(json)
      });
  }

  useEffect(() => {
    getZones()
  }, []);

  useEffect(() => {
    if (updatedZones) {
      getZones();
      updateZones(false);
    }
  }, [updatedZones]);

  return (
    <>
    <div className="App">
      <div className="relative overflow-hidden">
        <div className="mx-auto">
          <div className="relative ">
            <div className="page">
              <Sidebar
                sidebarOpen={sidebarOpen}
                setSidebarOpen={setSidebarOpen}
                setZones={setZones}
                zones={zones}
              />

              <div className={classNames("main-content", "w-100", "h-100", {
                  "padding-none": !sidebarOpen
                })}>

                <div className={"page-content"}>

                  <Switch>
                    {routes.map((route, idx) => (
                      <Route path={route.path} exact render={(props) => <route.component
                        {...props}
                        setZones={setZones}
                        zones={zones}
                        updateZones={updateZones} />}
                        key={idx} />
                    ))}
                    <Route component={NotFound} />
                  </Switch>

                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    </>
  );
}

export default withRouter(App);

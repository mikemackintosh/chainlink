import React, { useState, useEffect } from "react";
import classNames from "classnames";
import { useHistory,useLocation, Link } from "react-router-dom";

import AddNew from '../pages/addnew.js';

function Sidebar(props) {
  const { sidebarOpen, setSidebarOpen, setZones, zones } = props;

  let history = useHistory();

  return (
    <div className={classNames("sidebar", { nonVisible: !sidebarOpen })}>
      <div className="logo d-flex justify-content-between">
        <div className="logo-name">
          Chainlink
        </div>
      </div>
      <div className="vertical-menu pt-2">
        <ul>
          <li>
            <span>All zones</span>
          </li>
        </ul>
        <ul className="zones">
          {props.zones.length &&
            props.zones.map((zone, idx) => {
              console.log("zone", zone);

              return (
              <li
                className={classNames("channel-item", {})}
                key={idx}
              >
                <Link to={"/hosts/"+zone.zone}><span data-id={idx}># {zone.zone} (<b>{Object.keys(zone.endpoints).length || 0}</b>)</span></Link>
              </li>
            )})}
            <li className={classNames("channel-item", {})}>

              <Link to={"/hosts/addnew"}><span data-id={"addnew"}><b>Add New</b></span></Link>
            </li>

        </ul>
      </div>
    </div>
  );
}

export default Sidebar;

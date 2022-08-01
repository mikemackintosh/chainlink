import React, { useState, useEffect} from "react";
import classNames from "classnames";
import { Container, Row, Col } from "reactstrap";
import Sidebar from "../components/sidebar.js";

import { useLocation, Link, useParams } from 'react-router-dom';


const Hosts = (props) => {
  const [zoneData, setZoneData] = useState({});

  let { zone } = useParams();

  useEffect(() => {
    const found = props.zones.find(obj => {
      return obj.zone === zone;
    });
    console.log("found", found)
    setZoneData(found)
  }, [props])

  return (
    <>
      <div className="chat-content h-100">
        <div className="page-title-box">
          <h4>DNS Configuration for {zone}</h4>
        </div>
        <div className="col-12 col h-100 d-flex flex-column justify-content-space">
          <div className="w-100">
            <div className="flex flex-col justify-center h-full">
                <div className="w-full ma-w-2xl mx-auto shadow-lg rounded-sm border border-gray-200">
                    <div className="p-3">
                        <div className="overflow-x-auto">
                            <table className="table-auto w-full">
                                <thead className="text-xs font-semibold uppercase text-gray-400 bg-gray-50">
                                    <tr>
                                        <th className="p-2">
                                            <div className="font-semibold text-center">Subdomain</div>
                                        </th>
                                        <th className="p-2 whitespace-nowrap">
                                            <div className="font-semibold text-left">IP Address</div>
                                        </th>
                                        <th className="p-2 whitespace-nowrap">
                                            <div className="font-semibold text-left">Path</div>
                                        </th>
                                        <th className="p-2 whitespace-nowrap">
                                            <div className="font-semibold text-center">Upstream</div>
                                        </th>
                                    </tr>
                                </thead>
                                <tbody className="text-sm divide-y divide-gray-100">
                                {zoneData && zoneData.endpoints && Object.keys(zoneData.endpoints).map((item) => {
                                  return(
                                    <tr key={item}>
                                        <td className="p-2">
                                            <div className="flex-wrap items-center">
                                                <div className="font-medium text-gray-800">{item}</div>
                                                <div className="text-gray-400">{zoneData.endpoints[item].resolve.fqdn}</div>
                                            </div>
                                        </td>
                                        <td className="p-2 whitespace-nowrap">
                                            <div className="text-left">{zoneData.endpoints[item].resolve.value}</div>
                                        </td>
                                        <td className="p-2 whitespace-nowrap">
                                            <div className="text-left font-medium">{zoneData.endpoints[item].route.path}</div>
                                        </td>
                                        <td className="p-2 whitespace-nowrap">
                                            <div className="text-lg text-center">{zoneData.endpoints[item].route.upstreeam}</div>
                                        </td>
                                    </tr>)})}
                                </tbody>
                            </table>
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

export default Hosts;

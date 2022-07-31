import React, { useState, useEffect} from "react";

import { useLocation, Link } from 'react-router-dom';

const Zone = (props) => {
  useEffect(() => {
    console.log(props.zone.Endpoints);
  }, [])
  return (
    <>
  <div className="flex flex-col justify-center h-full">
      <div className="w-full ma-w-2xl mx-auto bg-white shadow-lg rounded-sm border border-gray-200">
          <header className="px-5 py-4 border-b border-gray-100">
              <h2 className="font-semibold text-gray-800">{props.zone.Zone}</h2>
          </header>
          <div className="p-3">
              <div className="overflow-x-auto">
                  <table className="table-auto w-full">
                      <thead className="text-xs font-semibold uppercase text-gray-400 bg-gray-50">
                          <tr>
                              <th className="p-2 whitespace-nowrap">
                                  <div className="font-semibold text-left">Subdomain</div>
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
                      {Object.keys(props.zone.Endpoints).map((key) => {
                          <tr>
                              <td className="p-2 whitespace-nowrap">
                                  <div className="flex items-center">
                                      <div className="font-medium text-gray-800">{key}sdfsd</div>
                                  </div>
                              </td>
                              <td className="p-2 whitespace-nowrap">
                                  <div className="text-left">alexshatov@gmail.com</div>
                              </td>
                              <td className="p-2 whitespace-nowrap">
                                  <div className="text-left font-medium text-green-500">$2,890.66</div>
                              </td>
                              <td className="p-2 whitespace-nowrap">
                                  <div className="text-lg text-center">??</div>
                              </td>
                          </tr>})}
                      </tbody>
                  </table>
              </div>
          </div>
      </div>
  </div>
    </>
  );
}

const Hosts = () => {
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState([])

  useEffect(() => {
    fetch("http://chainlink.config/api/config")
      .then((res) => res.json())
      .then((json) => {
    console.log(json)
        setData(json)
        setLoading(false)
      });

  }, []);

  return (
    <>
    {loading && <div>Loading</div>}
    {data.map(e => (<Zone key={e.Zone} zone={e} />))}

    </>
  );
}


export default Hosts;

import React, { useState, useEffect} from "react";
import classNames from "classnames";
import { Form, FormGroup, Label, FormFeedback, FormText, Input, Row, Col, Button } from "reactstrap";
import Sidebar from "../components/sidebar.js";

import { useLocation, Link, useParams } from 'react-router-dom';


const AddNew = (props) => {
  const [zoneData, setZoneData] = useState({});
  const [fqdnInput, setFqdnInput] = useState("");
  const [validate, setValidate] = useState({fqdn: null, upstream: null});
  const [upstreamInput, setUpstreamInput] = useState("");
  const { setZones, zones } = props;

  let { zone } = useParams();

  function onClickSave(e) {
    e.preventDefault();

    const fqdnPattern = /^(((?!\-))(xn\-\-)?[a-z0-9\-_]{0,61}[a-z0-9]{1,1}\.)*(xn\-\-)?([a-z0-9\-]{1,61}|[a-z0-9\-]{1,30})\.[a-z]{2,}$/i;
    if (fqdnPattern.test(fqdnInput)) {
      validate.fqdn = true
      setValidate({...validate, validate})
    } else {
      validate.fqdn = false
      setValidate({...validate,validate})
      return;
    }
    const upstreamPattern = /https?:\/\/(?:w{1,3}\.)?[^\s.]+(?:\.[a-z]+)*(?::\d+)?(?![^<]*(?:<\/\w+>|\/?>))/i;
    if (upstreamPattern.test(upstreamInput)) {
      validate.upstream = true
      setValidate({...validate, validate})
    } else {
      validate.upstream = false
      setValidate({...validate,validate})
      return;
    }

    const newMessage = {
      fqdn: fqdnInput,
      upstream: upstreamInput,
     };

     (async () => {
       const rawResponse = await fetch('http://chainlink.config/api/config', {
         method: 'POST',
         headers: {
           'Accept': 'application/json',
           'Content-Type': 'application/json'
         },
         body: JSON.stringify(newMessage)
       });

       const content = await rawResponse.json();
       console.log(content);

       const found = props.zones.find(obj => {
         return obj.zone === content.zone;
       });

       var newZoneList = zones
       console.log("zones", zones)
       newZoneList.push(content)
       console.log("newZoneList", newZoneList)
       setZones(() => [...zones, newZoneList]);
     })();


     setFqdnInput(() => "");
   }

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
          <h4>Add a new DNS Configuration</h4>
        </div>
        <div className="col-12 col h-100 d-flex flex-column justify-content-space">
          <div className="w-100">
            <div className="flex flex-col justify-center h-full">
                <div className="w-full ma-w-2xl mx-auto shadow-lg rounded-sm border border-gray-200">
                    <div className="p-3">
                        <div className="overflow-x-auto">

<Form className="w-100 position-relative">
  <Row>
    <Col md={6}>
      <FormGroup className="position-relative">
        <Label>Fully Qualified Domain Name</Label>
        <Input
          type="text"
          id="fqdn"
          invalid={validate.fqdn === false}
          valid={validate.fqdn === true}
          placeholder="example.service.tld"
          value={fqdnInput}
          onChange={(e) => setFqdnInput(e.target.value)}
        />
        <FormFeedback tooltip>
           This should be a FQDN, including subdomain, domain and TLD.
         </FormFeedback>
      </FormGroup>
    </Col>
    <Col md={6}>
      <FormGroup className="position-relative">
        <Label>Local Service URL</Label>
        <Input
          type="text"
          id="upstream"
          invalid={validate.upstream === false}
          valid={validate.upstream === true}
          placeholder="http://127.0.0.1:3000"
          value={upstreamInput}
          onChange={(e) => setUpstreamInput(e.target.value)}
        />
        <FormFeedback tooltip>
           Make sure to add a scheme and port number explicitly.
         </FormFeedback>
      </FormGroup>
    </Col>
  </Row>
 </Form>

 <Button onClick={onClickSave}      className="btn-rounded chat-send btn btn-primary">
  Add New Entry
</Button>

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

export default AddNew;

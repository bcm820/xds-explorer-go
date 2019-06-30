import React from "react";
import useForm from "./useForm";
import axios from "axios";

import { Column, Row } from "./styles";

function Request({ getResources }) {
  const { values, handleChange, handleSubmit } = useForm(
    formData => {
      console.log(formData);
      const reqBody = { ...formData };
      reqBody.resourceNames = reqBody.resourceNames
        .split(",")
        .map(s => s.trim());
      axios.post("/request", reqBody).then(setTimeout(getResources, 1000));
    },
    {
      resourceType: "ClusterLoadAssignment",
      node: "",
      zone: "",
      cluster: "",
      resourceNames: ""
    }
  );

  return (
    <form onSubmit={handleSubmit}>
      <Column>
        <Row>
          <label htmlFor={"resourceType"}>ResourceType</label>
          <select onChange={handleChange} name={"resoureType"} required>
            <option value="ClusterLoadAssignment">ClusterLoadAssignment</option>
            <option value="Cluster">Cluster</option>
            <option value="RouteConfiguration">RouteConfiguration</option>
            <option value="Listener">Listener</option>
            <option value="auth.Secret">auth.Secret</option>
          </select>
        </Row>
        <input
          type={"text"}
          onChange={handleChange}
          name={"node"}
          value={values.node}
          placeholder={"Node"}
        />
        <input
          type={"text"}
          onChange={handleChange}
          name={"zone"}
          value={values.zone}
          placeholder={"Zone"}
          required
        />
        <input
          type={"text"}
          onChange={handleChange}
          name={"cluster"}
          value={values.cluster}
          placeholder={"Cluster"}
          required
        />
        <input
          type={"text"}
          onChange={handleChange}
          name={"resourceNames"}
          value={values.resourceNames}
          placeholder={"ResourceNames"}
        />
        <input type={"submit"} value={"Request"} />
      </Column>
    </form>
  );
}

export default Request;

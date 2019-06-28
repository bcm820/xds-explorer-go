import React from "react";
import useForm from "./useForm";
import axios from "axios";

function Request({ incrementCount }) {
  const { values, handleChange, handleSubmit } = useForm(
    formData => {
      const request = { ...formData };
      request.resourceNames = request.resourceNames
        .split(",")
        .map(s => s.trim());
      axios.post("/request", request).then(() => incrementCount());
    },
    {
      resourceType: "",
      node: "",
      zone: "",
      cluster: "",
      resourceNames: ""
    }
  );

  return (
    <form onSubmit={handleSubmit}>
      <input
        type={"text"}
        onChange={handleChange}
        name={"resourceType"}
        value={values.resourceType}
        placeholder={"ResourceType"}
      />
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
      />
      <input
        type={"text"}
        onChange={handleChange}
        name={"cluster"}
        value={values.cluster}
        placeholder={"Cluster"}
      />
      <input
        type={"text"}
        onChange={handleChange}
        name={"resourceNames"}
        value={values.resourceNames}
        placeholder={"ResourceNames"}
      />
      <input type={"submit"} value={"Request"} />
    </form>
  );
}

export default Request;

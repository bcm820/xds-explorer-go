import React, { useState } from "react";

import Request from "./Request";
import JSONTree from "react-json-tree";
import axios from "axios";

import { Container, theme } from "./styles";

function App() {
  const [resources, setResources] = useState([]);
  const getResources = () =>
    axios.get("/listen").then(res => setResources(res.data));
  return (
    <Container>
      <Request getResources={getResources} />
      <JSONTree
        data={resources}
        theme={theme}
        invertTheme={true}
        shouldExpandNode={() => true}
      />
    </Container>
  );
}

export default App;

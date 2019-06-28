import React, { useState } from "react";
import Request from "./Request";
import Response from "./Response";
import "./App.css";

function App() {
  const [requestCount, setRequestCount] = useState(0);
  return (
    <>
      <Request incrementCount={() => setRequestCount(requestCount + 1)} />
      <Response requestCount={requestCount} />
    </>
  );
}

export default App;

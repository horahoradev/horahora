import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter as Router } from "react-router-dom";
import "antd/dist/antd.css";
import "./index.css";
import App from "./App";
import nightwind from "nightwind/helper";
nightwind.toggle();

ReactDOM.render(
  <React.StrictMode>
    <Router>
      <App />
    </Router>
  </React.StrictMode>,
  document.getElementById("root")
);

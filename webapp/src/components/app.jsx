import { BrowserRouter } from "react-router-dom";

import { Pages } from "../pages/_index";

import { Footer } from "./footer";

export function App() {
  return (
    <BrowserRouter>
      <Pages />
      <Footer />
    </BrowserRouter>
  );
}

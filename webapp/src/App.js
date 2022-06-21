import { BrowserRouter } from "react-router-dom";

import { Pages } from "./pages/_index";
import Footer from "./Footer";

function App() {
  return (
    <BrowserRouter >
      <Pages />
      <Footer></Footer>
    </BrowserRouter >
  );
}

export default App;

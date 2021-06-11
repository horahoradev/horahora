import { Switch, Route } from "react-router";
import HomePage from './HomePage';
import LoginPage from "./LoginPage";

function App() {
  return <div className="bg-yellow-50 min-h-screen font-serif">
    <Switch>
      <Route exact path="/">
        <HomePage />
      </Route>
      <Route exact path="/login">
        <LoginPage />
      </Route>
      <Route>
        TODO(ivan): 404
      </Route>
    </Switch>
  </div>;
}

export default App;

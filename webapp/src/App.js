import { Switch, Route } from "react-router";
import HomePage from './HomePage';

function App() {
  return <div className="bg-yellow-50 min-h-screen">
    <Switch>
      <Route exact path="/">
        <HomePage />
      </Route>
      <Route>
        TODO(ivan): 404
      </Route>
    </Switch>
  </div>;
}

export default App;

import { Route, Switch } from "react-router";

import HomePage from "./HomePage";
import LoginPage from "./LoginPage";
import LogoutPage from "./LogoutPage";
import VideoPage from "./VideoPage";
import UserPage from "./UserPage";
import ArchivalPage from "./ArchivalPage";
import RegisterPage from "./RegisterPage";
import PasswordResetPage from "./Passwordreset";
import AuditPage from "./AuditPage";
import Footer from "./Footer";
import { ThemeSwitcherProvider } from 'react-css-theme-switcher';


function App() {
  const themes = {
    light: '/antd.min.css',
    dark: '/antd.dark.min.css',
  };

  return (
    <ThemeSwitcherProvider defaultTheme="dark" themeMap={themes}>
    <div className=" bg-yellow-50 dark:bg-gray-900 min-h-screen font-sans-serif">
      <Switch>
        <Route exact path="/">
          <HomePage />
        </Route>
        <Route exact path="/login">
          <LoginPage />
        </Route>
        <Route exact path="/logout">
          <LogoutPage />
        </Route>
        <Route exact path="/videos/:id">
          <VideoPage />
        </Route>
        <Route exact path="/users/:id">
          <UserPage />
        </Route>
        <Route exact path="/register">
          <RegisterPage></RegisterPage>
        </Route>
        <Route exact path ="/archive-requests">
          <ArchivalPage/>
        </Route>
        <Route exact path ="/password-reset">
            <PasswordResetPage/>
        </Route>
        <Route exact path ="/audits">
          <AuditPage/>
        </Route>
        <Route>TODO(ivan): 404</Route>
      </Switch>
    </div>
    <Footer></Footer>
    </ThemeSwitcherProvider>
  );
}

export default App;

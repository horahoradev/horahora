import { Route, Switch } from "react-router";

import { HomePage } from "./home";
import { LoginPage } from "./login";
import { LogoutPage } from "./logout";
import { VideoPage } from "./video";
import UserPage from "./UserPage";
import { ArchivalPage } from "./archival";
import { RegisterPage } from "./register";
import PasswordResetPage from "./Passwordreset";
import { AuditPage } from "./AuditPage";

export function Pages() {
  return (
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
      <Route exact path="/archive-requests">
        <ArchivalPage />
      </Route>
      <Route exact path="/password-reset">
        <PasswordResetPage />
      </Route>
      <Route exact path="/audits">
        <AuditPage />
      </Route>
      <Route>TODO(ivan): 404</Route>
    </Switch>
  )
}
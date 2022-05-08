import { useFormik } from "formik";
import { Link, Button, Input } from "antd";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faKey, faUser } from "@fortawesome/free-solid-svg-icons";
import { useEffect, useRef } from "react";
import { useHistory } from "react-router-dom";
import Footer from "./Footer";

import Header from "./Header";

import * as API from "./api";

function LoginForm() {
  let history = useHistory();
  // TODO(ivan): validation, form errors
  // TODO(ivan): submitting state
  let formik = useFormik({
    initialValues: {
      username: "",
      password: "",
    },
    onSubmit: async (values) => {
      await API.postLogin(values);
      history.push("/");
    },
  });

  // automatically focus input on first input on render
  let usernameInputRef = useRef();
  useEffect(() => {
    usernameInputRef.current && usernameInputRef.current.focus();
  }, [usernameInputRef]);

  return (
    <div className="max-w-xs w-full border rounded shadow bg-white p-4">
      <h2 className="text-xl mb-4 inline-block">Welcome back!</h2> <a className="float-right -top-5" href="/register">register</a>
      <form onSubmit={formik.handleSubmit}>
        <Input.Group>
          <Input
            name="username"
            values={formik.values.username}
            onChange={formik.handleChange}
            size="large"
            placeholder="Username"
            ref={usernameInputRef}
            prefix={
              <FontAwesomeIcon className="mr-1 text-gray-400" icon={faUser} />
            }
          />
        </Input.Group>
        <br />
        <Input.Group>
          <Input.Password
            name="password"
            values={formik.values.username}
            onChange={formik.handleChange}
            size="large"
            placeholder="Password"
            prefix={
              <FontAwesomeIcon className="mr-1 text-gray-400" icon={faKey} />
            }
          />
        </Input.Group>
        <br />
        <Input.Group>
          <Button block type="primary" htmlType="submit" size="large">
            Submit
          </Button>
        </Input.Group>
      </form>
    </div>
  );
}

function LoginPage() {
  return (
    <>
      <Header dataless />
      <div className="flex justify-center mx-4">
        <div className="max-w-screen-lg w-screen my-6 flex justify-center items-center pt-32">
          <LoginForm />
        </div>
      </div>
    </>
  );
}

export default LoginPage;

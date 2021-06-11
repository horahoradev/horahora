import Header from "./Header";
import { useFormik } from "formik";
import { Button, Input } from "antd";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faKey, faUser } from "@fortawesome/free-solid-svg-icons";
import { useRef, useEffect } from "react";
import { useHistory } from "react-router-dom";
import { postLogin } from "./api";

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
      await postLogin(values);
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
      <h2 className="text-xl mb-4">Welcome back!</h2>
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

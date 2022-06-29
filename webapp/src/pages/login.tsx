import { useFormik } from "formik";
import { Button, Input } from "antd";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faKey, faUser } from "@fortawesome/free-solid-svg-icons";
import { useEffect, useRef } from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import type { InputRef } from "antd";

import { Header } from "#components/header";
import { postLogin } from "#api/index";

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

function LoginForm() {
  const router = useRouter();
  // TODO(ivan): validation, form errors
  // TODO(ivan): submitting state
  let formik = useFormik({
    initialValues: {
      username: "",
      password: "",
    },
    onSubmit: async (values) => {
      await postLogin(values);
      router.push("/");
    },
  });

  // automatically focus input on first input on render
  let usernameInputRef = useRef<InputRef>(null);
  useEffect(() => {
    usernameInputRef.current && usernameInputRef.current.focus();
  }, [usernameInputRef]);

  return (
    <div className="max-w-xs w-full border rounded shadow bg-white dark:bg-gray-800 p-4">
      <h2 className="text-xl text-black dark:text-white mb-4 inline-block">
        Welcome back!
      </h2>{" "}
      <Link
        className="float-right -top-5 text-black dark:text-white"
        href="/register"
      >
        register
      </Link>
      <form onSubmit={formik.handleSubmit}>
        <Input.Group>
          <Input
            name="username"
            // @ts-expect-error types
            values={formik.values.username}
            onChange={formik.handleChange}
            size="large"
            placeholder="Username"
            ref={usernameInputRef}
            prefix={
              <FontAwesomeIcon
                className="max-h-4 mr-1 text-gray-400"
                icon={faUser}
              />
            }
          />
        </Input.Group>
        <br />
        <Input.Group>
          <Input.Password
            name="password"
            // @ts-expect-error types
            values={formik.values.username}
            onChange={formik.handleChange}
            size="large"
            placeholder="Password"
            prefix={
              <FontAwesomeIcon
                className="max-h-4 mr-1 text-gray-400"
                icon={faKey}
              />
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

export default LoginPage;

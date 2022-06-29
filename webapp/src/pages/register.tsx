import { useFormik } from "formik";
import { Button, Input } from "antd";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faKey, faMailBulk, faUser } from "@fortawesome/free-solid-svg-icons";
import { useEffect, useRef } from "react";
import { useRouter } from "next/router";
import type { InputRef } from "antd";

import { PrivacyPolicy } from "#components/privacy-policy";
import { TermsOfService } from "#components/terms-of-service";
import { Header } from "#components/header";
import { postRegister } from "#api/index";

function RegisterPage() {
  return (
    <>
      <Header dataless />
      <div className="flex justify-center mx-4">
        <div className="max-w-screen-lg w-screen my-6 flex justify-center items-center pt-32">
          <RegistrationForm />
        </div>
      </div>
    </>
  );
}

function RegistrationForm() {
  const router = useRouter();
  // TODO(ivan): validation, form errors
  // TODO(ivan): submitting state
  let formik = useFormik({
    initialValues: {
      email: "",
      username: "",
      password: "",
    },
    onSubmit: async (values) => {
      await postRegister(values);
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
      <h2 className="text-xl text-black dark:text-white mb-4">Register</h2>
      <form onSubmit={formik.handleSubmit}>
        <Input.Group>
          <Input
            name="email"
            // @ts-expect-error form types
            values={formik.values.email}
            onChange={formik.handleChange}
            size="large"
            placeholder="email"
            ref={usernameInputRef}
            prefix={
              <FontAwesomeIcon
                className="max-h-4 mr-1 text-gray-400"
                icon={faUser}
              />
            }
          />
        </Input.Group>
        <Input.Group>
          <Input
            name="username"
            // @ts-expect-error form types
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
        <Input.Group>
          <Input.Password
            name="password"
            // @ts-expect-error form types
            values={formik.values.password}
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
        <div className="text-black dark:text-white">
          By submitting, you agree to the <PrivacyPolicy /> and{" "}
          <TermsOfService />
        </div>
        <Input.Group>
          <Button block type="primary" htmlType="submit" size="large">
            Submit
          </Button>
        </Input.Group>
      </form>
    </div>
  );
}

export default RegisterPage;

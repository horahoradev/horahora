import { useFormik } from "formik";
import { Button, Input } from "antd";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faKey, faMailBulk, faUser } from "@fortawesome/free-solid-svg-icons";
import { useEffect, useRef } from "react";
import { useRouter } from "next/router";
import type { InputRef } from "antd";

import { Header } from "#components/header";
import { postPasswordReset } from "#api/index";

function PasswordResetPage() {
  return (
    <>
      <Header dataless />
      <div className="flex justify-center mx-4">
        <div className="max-w-screen-lg w-screen my-6 flex justify-center items-center pt-32">
          <PasswordResetForm />
        </div>
      </div>
    </>
  );
}

function PasswordResetForm() {
  const router = useRouter();
  // TODO(ivan): validation, form errors
  // TODO(ivan): submitting state
  let formik = useFormik({
    initialValues: {
      old_password: "",
      new_password: "",
    },
    onSubmit: async (values) => {
      await postPasswordReset(values);
      router.push("/");
    },
  });

  // automatically focus input on first input on render
  let usernameInputRef = useRef<InputRef>(null);

  useEffect(() => {
    usernameInputRef.current && usernameInputRef.current.focus();
  }, [usernameInputRef]);

  return (
    <div className="max-w-xs w-full border rounded shadow bg-white dark:bg-black p-4">
      <h2 className="text-xl text-black dark:text-white mb-4">
        Reset Password
      </h2>
      <form onSubmit={formik.handleSubmit}>
        <Input.Group>
          <Input
            name="old_password"
            // @ts-expect-error form types
            values={formik.values.old_password}
            onChange={formik.handleChange}
            size="large"
            placeholder="current password"
            ref={usernameInputRef}
            prefix={
              <FontAwesomeIcon className="max-h-4 mr-1 text-gray-400" icon={faUser} />
            }
          />
        </Input.Group>
        <Input.Group>
          <Input
            name="new_password"
            // @ts-expect-error form types
            values={formik.values.new_password}
            onChange={formik.handleChange}
            size="large"
            placeholder="new password"
            ref={usernameInputRef}
            prefix={
              <FontAwesomeIcon className="max-h-4 mr-1 text-gray-400" icon={faUser} />
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

export default PasswordResetPage;

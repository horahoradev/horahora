import { useFormik } from "formik";
import { Button, Input } from "antd";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faKey, faMailBulk, faUser } from "@fortawesome/free-solid-svg-icons";
import { useEffect, useRef } from "react";
import { useHistory } from "react-router-dom";

import Header from "./Header";
import * as API from "./api";

function PasswordResetForm() {
  let history = useHistory();
  // TODO(ivan): validation, form errors
  // TODO(ivan): submitting state
  let formik = useFormik({
    initialValues: {
      old_password: "",
      new_password: "",
    },
    onSubmit: async (values) => {
      await API.postPasswordReset(values);
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
      <h2 className="text-xl mb-4">Reset Password</h2>
      <form onSubmit={formik.handleSubmit}>
      <Input.Group>
          <Input
            name="old_password"
            values={formik.values.old_password}
            onChange={formik.handleChange}
            size="large"
            placeholder="current password"
            ref={usernameInputRef}
            prefix={
              <FontAwesomeIcon className="mr-1 text-gray-400" icon={faUser} />
            }
          />
          </Input.Group>
        <Input.Group>
          <Input
            name="new_password"
            values={formik.values.new_password}
            onChange={formik.handleChange}
            size="large"
            placeholder="new password"
            ref={usernameInputRef}
            prefix={
              <FontAwesomeIcon className="mr-1 text-gray-400" icon={faUser} />
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

export default PasswordResetPage;

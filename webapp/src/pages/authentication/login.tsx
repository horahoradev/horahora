import { useRouter } from "next/router";

import { loginAccount } from "#api/authentication";
import { LinkInternal } from "#components/links";
import {
  FormClient,
  type IFormElements,
  type ISubmitEvent,
} from "#components/forms";
import { Page } from "#components/page";
import { Password, Text } from "#components/inputs";

const FIELD_NAMES = {
  NAME: "username",
  PASSWORD: "password",
} as const;
type IFieldName = typeof FIELD_NAMES[keyof typeof FIELD_NAMES];

function LoginPage() {
  const router = useRouter();

  async function handleSubmit(event: ISubmitEvent) {
    const fields = event.currentTarget.elements as IFormElements<IFieldName>;
    const formData = Object.values(FIELD_NAMES).reduce(
      (formData, fieldName) => {
        switch (fieldName) {
          case FIELD_NAMES.NAME:
          case FIELD_NAMES.PASSWORD: {
            const fieldElement = fields[fieldName];
            formData.set(fieldName, fieldElement.value);
            break;
          }

          default:
            throw new Error(
              `The field "${fieldName}" is missing from the form.`
            );
        }

        return formData;
      },
      new FormData()
    );
    await loginAccount(formData);
    router.push("/");
  }

  return (
    <Page>
      <FormClient id="auth-login" onSubmit={handleSubmit}>
        <p>
          Welcome back!{" "}
          <LinkInternal
            className="float-right -top-5 text-black dark:text-white"
            href="/register"
          >
            register
          </LinkInternal>
        </p>
        <Text id="auth-login-name" name={FIELD_NAMES.NAME}>
          Name
        </Text>
        <Password id="auth-login-password" name={FIELD_NAMES.PASSWORD}>
          Password
        </Password>
      </FormClient>
    </Page>
  );
}

export default LoginPage;

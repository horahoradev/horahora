import { useRouter } from "next/router";

import { LinkInternal } from "#components/links";
import {
  FormClient,
  type IFormElements,
  type ISubmitEvent,
} from "#components/forms";
import { Page } from "#components/page";
import { Password, Text } from "#components/inputs";
import { useAccount } from "#hooks";
import { type IAccountLogin } from "#lib/account";

const FIELD_NAMES = {
  NAME: "username",
  PASSWORD: "password",
} as const;
type IFieldName = typeof FIELD_NAMES[keyof typeof FIELD_NAMES];

function LoginPage() {
  const router = useRouter();
  const { login } = useAccount();

  async function handleSubmit(event: ISubmitEvent) {
    const fields = event.currentTarget.elements as IFormElements<IFieldName>;
    const accLogin = Object.values(FIELD_NAMES).reduce(
      (accLogin, fieldName) => {
        switch (fieldName) {
          case FIELD_NAMES.NAME:
          case FIELD_NAMES.PASSWORD: {
            const fieldElement = fields[fieldName];
            accLogin[fieldName] = fieldElement.value;
            break;
          }

          default:
            throw new Error(
              `The field "${fieldName}" is missing from the form.`
            );
        }

        return accLogin;
      },
      {} as IAccountLogin
    );

    await login(accLogin);
    router.push("/");
  }

  return (
    <Page title="Log in">
      <FormClient id="auth-login" onSubmit={handleSubmit}>
        <p>
          Not registered?{" "}
          <LinkInternal href="/authentication/register">Register</LinkInternal>
        </p>
        <Text
          id="auth-login-name"
          name={FIELD_NAMES.NAME}
          autoComplete="username"
        >
          Name
        </Text>
        <Password
          id="auth-login-password"
          name={FIELD_NAMES.PASSWORD}
          autoComplete="current-password"
        >
          Password
        </Password>
      </FormClient>
    </Page>
  );
}

export default LoginPage;

import { useRouter } from "next/router";

import { registerAccount } from "#api/authentication";
import { LinkInternal } from "#components/links";
import { Page } from "#components/page";
import {
  FormClient,
  type IFormElements,
  type ISubmitEvent,
} from "#components/forms";
import { Text, Password, Email } from "#components/inputs";

const FIELD_NAMES = {
  NAME: "username",
  EMAIL: "email",
  PASSWORD: "password",
} as const;
type IFieldName = typeof FIELD_NAMES[keyof typeof FIELD_NAMES];

function RegisterPage() {
  const router = useRouter();

  async function handleSubmit(event: ISubmitEvent) {
    const fields = event.currentTarget.elements as IFormElements<IFieldName>;
    const formParams = Object.values(FIELD_NAMES).reduce(
      (formParams, fieldName) => {
        switch (fieldName) {
          case FIELD_NAMES.NAME:
          case FIELD_NAMES.EMAIL:
          case FIELD_NAMES.PASSWORD: {
            const fieldElement = fields[fieldName];
            formParams.set(fieldName, fieldElement.value);
            break;
          }

          default:
            throw new Error(
              `The field "${fieldName}" is missing from the form.`
            );
        }

        return formParams;
      },
      new URLSearchParams()
    );
    await registerAccount(formParams);
    router.push("/");
  }

  return (
    <Page title="Register">
      <FormClient id="auth-register" onSubmit={handleSubmit}>
        <Text id="auth-register-username" name={FIELD_NAMES.NAME}>
          Name
        </Text>
        <Email id="auth-register-email" name={FIELD_NAMES.EMAIL}>
          Email
        </Email>
        <Password
          id="auth-register-password"
          name={FIELD_NAMES.PASSWORD}
          autoComplete="new-password"
        >
          Password
        </Password>
        {/* These should be checkboxes. */}
        <p>
          By submitting, you agree to the{" "}
          <LinkInternal href="/privacy-policy">Privacy Policy</LinkInternal> and{" "}
          <LinkInternal href="/terms-of-service">TermsOfService</LinkInternal>.
        </p>
      </FormClient>
    </Page>
  );
}

export default RegisterPage;

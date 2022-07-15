import { useRouter } from "next/router";

import { resetAccountPassword } from "#api/authentication";
import { Page } from "#components/page";
import {
  FormClient,
  type IFormElements,
  type ISubmitEvent,
} from "#components/forms";
import { Password } from "#components/inputs";

const FIELD_NAMES = {
  OLD: "old_password",
  NEW: "new_password",
} as const;
type IFieldName = typeof FIELD_NAMES[keyof typeof FIELD_NAMES];

function PasswordResetPage() {
  const router = useRouter();

  async function handleSubmit(event: ISubmitEvent) {
    const fields = event.currentTarget.elements as IFormElements<IFieldName>;

    const formData = Object.values(FIELD_NAMES).reduce(
      (formData, fieldName) => {
        switch (fieldName) {
          case FIELD_NAMES.OLD:
          case FIELD_NAMES.NEW: {
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

    await resetAccountPassword(formData);
    router.push("/");
  }

  return (
    <Page>
      <FormClient id="auth-reset" onSubmit={handleSubmit}>
        <p>Reset Password</p>
        <Password
          id="auth-reset-old"
          name={FIELD_NAMES.OLD}
          autoComplete="current-password"
        >
          Current Password
        </Password>
        <Password
          id="auth-reset-new"
          name={FIELD_NAMES.NEW}
          autoComplete="new-password"
        >
          New Password
        </Password>
      </FormClient>
    </Page>
  );
}

export default PasswordResetPage;

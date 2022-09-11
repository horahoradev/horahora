import { Input } from "./input";
import { type IPasswordAutoComplete, PASSWORD_AUTOCOMPLETE } from "./types";

import { FormSection, IFormSectionProps, Label } from "#components/forms";
import { blockComponent } from "#components/meta";

export interface IPasswordProps extends IFormSectionProps {
  id: string;
  name: string;
  autoComplete?: IPasswordAutoComplete;
  minLength?: number;
  maxLength?: number;
  required?: boolean;
}

export const Password = blockComponent(undefined, Component);

function Component({
  id,
  name,
  autoComplete = PASSWORD_AUTOCOMPLETE.ON,
  required,
  minLength,
  maxLength,
  children,
  ...blockProps
}: IPasswordProps) {
  return (
    <FormSection {...blockProps}>
      <Label htmlFor={id}>{children}</Label>
      <Input
        id={id}
        name={name}
        type="password"
        autoComplete={autoComplete}
        required={required}
        minLength={minLength}
        maxLength={maxLength}
      />
    </FormSection>
  );
}

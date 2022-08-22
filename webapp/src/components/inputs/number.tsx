import { Input } from "./input";
import { IInputSectionProps } from "./types";

import { FormSection } from "#components/forms";
import { blockComponent } from "#components/meta";

export interface INumberProps extends IInputSectionProps {}

export const Number = blockComponent(undefined, Component);

function Component({ id, name, disabled, form, required, ...blockProps }: INumberProps) {
  return (
    <FormSection {...blockProps}>
      <Input
        id={id}
        name={name}
        type="number"
        disabled={disabled}
        form={form}
        required={required}
      />
    </FormSection>
  );
}

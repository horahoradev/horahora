import { Input } from "./input";

import { FormSection, Label, type IFormSectionProps } from "#components/forms";
import { blockComponent } from "#components/meta";

export interface IRadioProps extends IFormSectionProps {
  id: string;
  name: string;
  value: string;
  checked?: boolean;
}

export const Radio = blockComponent(undefined, Component);

function Component({
  id,
  name,
  value,
  checked,
  children,
  ...blockProps
}: IRadioProps) {
  return (
    <FormSection {...blockProps}>
      <Input
        id={id}
        name={name}
        type="radio"
        value={value}
        defaultChecked={checked}
      />
      <Label htmlFor={id}>{children}</Label>
    </FormSection>
  );
}

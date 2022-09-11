import { Input } from "./input";
import styles from "./radio.module.scss";

import { FormSection, Label, type IFormSectionProps } from "#components/forms";
import { blockComponent } from "#components/meta";

export interface IRadioProps extends IFormSectionProps {
  id: string;
  name: string;
  value: string;
  checked?: boolean;
}

export const Radio = blockComponent(styles.block, Component);

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
        className={styles.input}
        name={name}
        type="radio"
        value={value}
        defaultChecked={checked}
      />
      <Label className={styles.label} htmlFor={id}>
        {children}
      </Label>
    </FormSection>
  );
}

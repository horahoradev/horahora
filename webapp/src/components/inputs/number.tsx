import { Ref } from "react";

import { Input } from "./input";
import { IInputSectionProps } from "./types";

import { FormSection } from "#components/forms";
import { blockComponent } from "#components/meta";

// eslint-disable-next-line
import styles from "./number.module.scss";

export interface INumberInputProps extends IInputSectionProps {
  inputRef?: Ref<HTMLInputElement>;
  min?: number;
  max?: number;
  step?: number;
}

/**
 * @TODO antd blocks the style reset somehow
 */
export const NumberInput = blockComponent(styles.block, Component);

function Component({
  id,
  name,
  min,
  max,
  step,
  disabled,
  form,
  required,
  defaultValue,
  inputRef,
  ...blockProps
}: INumberInputProps) {
  return (
    <FormSection {...blockProps}>
      <Input
        id={id}
        className={styles.input}
        name={name}
        type="number"
        min={min}
        max={max}
        step={step}
        disabled={disabled}
        form={form}
        required={required}
        ref={inputRef}
        defaultValue={defaultValue}
      />
    </FormSection>
  );
}

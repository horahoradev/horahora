import { Ref } from "react";

import { Input } from "./input";
import { IInputSectionProps } from "./types";

import { FormSection } from "#components/forms";
import { blockComponent } from "#components/meta";

export interface INumberProps extends IInputSectionProps {
  inputRef?: Ref<HTMLInputElement>;
  min?: number;
  max?: number;
  step?: number;
}

export const Number = blockComponent(undefined, Component);

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
}: INumberProps) {
  return (
    <FormSection {...blockProps}>
      <Input
        id={id}
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

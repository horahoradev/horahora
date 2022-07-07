import { type ReactNode } from "react";

import { Radio } from "./radio";

import { Fieldset, Legend, type IFieldsetProps } from "#components/forms";
import { blockComponent } from "#components/meta";

export interface IRadioGroupProps extends IFieldsetProps {
  name: string;
  options: IRadioOption[];
}

export interface IRadioOption {
  id: string;
  title: ReactNode;
  value: string;
}

export const RadioGroup = blockComponent(undefined, Component);

function Component({
  name,
  options,
  children,
  ...blockProps
}: IRadioGroupProps) {
  return (
    <Fieldset {...blockProps}>
      <Legend>{children}</Legend>
      {options.map(({ id, title, value }) => (
        <Radio key={id} id={id} name={name} value={value}>
          {title}
        </Radio>
      ))}
    </Fieldset>
  );
}

import { type ReactNode } from "react";

import { Radio } from "./radio";
import styles from "./radio-group.module.scss";

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

export const RadioGroup = blockComponent(styles.block, Component);

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
        <Radio key={id} id={id} name={name} className={styles.radio} value={value}>
          {title}
        </Radio>
      ))}
    </Fieldset>
  );
}

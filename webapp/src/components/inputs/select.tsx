import styles from "./select.module.scss";

import { FormSection, Label, type IFormSectionProps } from "#components/forms";
import { blockComponent, IBlockProps } from "#components/meta";

export interface ISelectProps extends IFormSectionProps {
  id: string;
  name: string;
  options: IOptionProps[];
}

export interface IOptionProps extends IBlockProps<"option"> {
  title: string;
}

export const Select = blockComponent(styles.block, Component);

function Component({
  id,
  name,
  options,
  children,
  ...blockProps
}: ISelectProps) {
  return (
    <FormSection {...blockProps}>
      <Label htmlFor={id}>{children}</Label>
      <select id={id} className={styles.select} name={name}>
        {options.map(({ title, ...optionProps }, index) => (
          <option key={index} className={styles.option} {...optionProps}>
            {title}
          </option>
        ))}
      </select>
    </FormSection>
  );
}

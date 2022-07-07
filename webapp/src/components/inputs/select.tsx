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

export const Select = blockComponent(undefined, Component);

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
      <select id={id} name={name}>
        {options.map(({ title, ...optionProps }, index) => (
          <option key={index} {...optionProps}>
            {title}
          </option>
        ))}
      </select>
    </FormSection>
  );
}

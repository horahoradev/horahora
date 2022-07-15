import { Form, type IFormProps } from "./form";
import { type ISubmitEvent } from "./types";

import { blockComponent, ClientComponent } from "#components/meta";

export interface IFormClientProps extends Omit<IFormProps, "method"> {
  onSubmit: (event: ISubmitEvent) => Promise<void>;
}

/**
 * Form which only renders upon hydration.
 *
 * Submit is default prevented and `onSubmit()` is always async and is required.
 */
export const FormClient = blockComponent(undefined, Component);

export function Component({
  onSubmit,
  children,
  ...blockProps
}: IFormClientProps) {
  async function handleSubmit(
    ...args: Parameters<IFormClientProps["onSubmit"]>
  ) {
    const [event] = args;
    event.preventDefault();
    await onSubmit(event);
  }

  return (
    <ClientComponent>
      <Form onSubmit={handleSubmit} {...blockProps}>
        {children}
      </Form>
    </ClientComponent>
  );
}

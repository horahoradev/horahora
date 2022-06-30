import type { FocusEvent } from "react";

/**
 * `onBlur` callback but only triggers the passed `callback`
 * when element's descendants aren't focused.
 * @param callback A function to invoke on successful blur.
 * @param callbackArgs Arguments to pass to the said function, if any.
 */
export function onParentBlur<DOMInterface extends Element>(
  callback: (...args: any[]) => void,
  ...callbackArgs: any[]
) {
  return (event: FocusEvent<DOMInterface>) => {
    if (event.currentTarget.contains(event.relatedTarget)) {
      event.preventDefault();
      return;
    }

    if (callbackArgs) {
      callback(...callbackArgs);
      return;
    }

    callback();
  };
}
